package service_test

import (
	"github.com/bernardosecades/feeder/pkg/service"
	"github.com/bernardosecades/feeder/pkg/value"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestServiceReportWithoutConcurrency(t *testing.T) {
	mock := &MockSkuRepository{}
	svc := service.NewService(mock, MockLoggerSvc{})

	svc.AddSku("KASL-3423") // valid
	svc.AddSku("KASL-7770") // valid
	svc.AddSku("KASL-1234") // valid
	svc.AddSku("KASL-1234") // duplicated
	svc.AddSku("765-1234")  // invalid

	totalUnique, totalDuplicated, totalInvalid := svc.Report()

	assert.EqualValues(t, 3, totalUnique)
	assert.EqualValues(t, 1, totalDuplicated)
	assert.EqualValues(t, 1, totalInvalid)
}

func TestServiceReportRunSafelyConcurrently(t *testing.T) {
	svc := service.NewService(MockSkuRepository{}, MockLoggerSvc{})

	numberRoutines := 500
	var wg sync.WaitGroup
	wg.Add(numberRoutines)

	for i := 0; i < numberRoutines; i++ {
		go func() {
			svc.AddSku("765-1234")  // invalid
			svc.AddSku("KASL-1234") // valid
			wg.Done()
		}()
	}
	wg.Wait()

	totalUnique, totalDuplicated, totalInvalid := svc.Report()

	assert.EqualValues(t, 1, totalUnique)
	assert.EqualValues(t, numberRoutines, totalInvalid)
	// we put 1 to give context -> 1 = "KASL-1234" is the unique valid so will
	// duplicate the valid minus the first time we add (is not duplicated)
	assert.EqualValues(t, (numberRoutines * 1) - 1, totalDuplicated)
}

func TestServicePersistWhenStorageDontContainAnySkuAddedInThisRunning(t *testing.T) {
	mock := &MockSkuRepository{}
	svc := service.NewService(mock, MockLoggerSvc{})

	svc.AddSku("KASL-3423")
	svc.AddSku("KASL-7770")

	mock.fnPersist = func(block map[string]value.Sku) (int64, error) {
		return 1, nil
	}

	totalInserted, totalSkipped, err :=svc.Persist()

	assert.Nil(t, err)
	assert.EqualValues(t, 1, totalInserted)
	assert.EqualValues(t, 1, totalSkipped)
}

func TestServicePersistWhenStorageAlreadyContainOneSkuAddedInThisRunning(t *testing.T) {
	mock := &MockSkuRepository{}
	svc := service.NewService(mock, MockLoggerSvc{})

	svc.AddSku("KASL-3423")
	svc.AddSku("KASL-7770")

	mock.fnPersist = func(block map[string]value.Sku) (int64, error) {
		return 2, nil
	}

	totalInserted, totalSkipped, err :=svc.Persist()

	assert.Nil(t, err)
	assert.EqualValues(t, 2, totalInserted)
	assert.EqualValues(t, 0, totalSkipped)
}

type MockSkuRepository struct {
	fnPersist func(block map[string]value.Sku) (int64, error)
	fnDelete func(block map[string]value.Sku) (int64, error)
}

func (m MockSkuRepository) Persist(block map[string]value.Sku) (int64, error) {
	if m.fnPersist != nil {
		return m.fnPersist(block)
	}
	return 0, nil
}
func (m MockSkuRepository) Delete(block map[string]value.Sku) (int64, error) {
	if m.fnDelete != nil {
		return m.fnDelete(block)
	}
	return 0, nil
}

type MockLoggerSvc struct {
}

func (m MockLoggerSvc) Log(v ...interface{}) {
}