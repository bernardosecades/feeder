package service

import (
	"github.com/bernardosecades/feeder/pkg/logger"
	"github.com/bernardosecades/feeder/pkg/repository"
	"github.com/bernardosecades/feeder/pkg/value"

	"sync"
)

type SkusInserted int
type SkusInsertSkipped int

type TotalUniqueSkus int
type TotalDuplicatedSkus int
type TotalInvalidSkus int

type Feeder interface {
	Persist() (SkusInserted, SkusInsertSkipped, error)
	Report() (TotalUniqueSkus, TotalDuplicatedSkus, TotalInvalidSkus)
	Log()
	AddSku(sku string)
}

type feeder struct {
	skuRepository repository.Sku
	logger        logger.Logger
	skus          map[string]value.Sku
	invalid       int
	duplicated    int
	mx            *sync.Mutex
}

// NewService create new instance from service.Feeder
func NewService(skuRepository repository.Sku, logger logger.Logger) Feeder {
	return &feeder{
		skuRepository: skuRepository,
		logger:        logger,
		skus:          map[string]value.Sku{},
		invalid:       0,
		duplicated:    0,
		mx:            new(sync.Mutex),
	}
}

// AddSku it will add new sku only if is valid and is not duplicated in the current running application.
// It will increment counter for invalid and duplicate sku for current running application. It is ready
// to be safe with concurrency using lock system.
func (s *feeder) AddSku(sku string) {
	sk, err := value.NewSku(sku)

	// we block all goroutines until the mutex is unlocked to avoid race conditions
	s.mx.Lock()
	defer s.mx.Unlock()

	if err != nil {
		s.invalid++
		return
	}

	if _, found := s.skus[sk.String()]; found {
		s.duplicated++
		return
	}

	s.skus[sk.String()] = sk
	// NOTE: we could log here (because here are uniques skus) but we use method Log called in server to improve
	// the performance because if not, each message will access to file log to write so that is a bad performance
	// so we log at the end.
}

// Log log unique sku from running application
func (s *feeder) Log() {
	for _, v := range s.skus {
		s.logger.Log("Added sku:", v.StringWithoutZeros())
	}
}

// Persist it will persist sku previously stored in memory and will return skus inserted
// and skipped: number of skipped is because can happen a valid sku in a running application was already persisted
// in other running application.
func (s *feeder) Persist() (SkusInserted, SkusInsertSkipped, error) {
	skuInserted, err := s.skuRepository.Persist(s.skus)
	if err != nil {
		return 0, 0, err
	}
	skipped := int64(len(s.skus)) - skuInserted

	return SkusInserted(skuInserted), SkusInsertSkipped(skipped), nil

}

// Report it will return summary of skus: unique, duplicated and invalid in current running application.
func (s *feeder) Report() (TotalUniqueSkus, TotalDuplicatedSkus, TotalInvalidSkus) {
	return TotalUniqueSkus(len(s.skus)), TotalDuplicatedSkus(s.duplicated), TotalInvalidSkus(s.invalid)
}
