package repository_test

import (
	"github.com/bernardosecades/feeder/pkg/repository"
	"github.com/bernardosecades/feeder/pkg/tools/env"
	"github.com/bernardosecades/feeder/pkg/value"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPersistAndDelete(t *testing.T) {
	r := repository.NewSkuPostgreSQL(
		env.GetEnvOrFallback("DB_HOST", "localhost"),
		env.GetEnvOrFallback("DB_PORT", "5416"),
		env.GetEnvOrFallback("DB_USER", "feeder"),
		env.GetEnvOrFallback("DB_PASS", "feeder"),
		env.GetEnvOrFallback("DB_NAME", "feeder"),
	)

	// Persist and Delete with items in data
	data := make(map[string]value.Sku)

	sku1, _ := value.NewSku("KASL-3423")
	sku2, _ := value.NewSku("KASL-7777")
	sku3, _ := value.NewSku("KASL-7777")

	data[sku1.String()] = sku1
	data[sku2.String()] = sku2
	data[sku2.String()] = sku3

	rowsInserted, err := r.Persist(data)
	assert.Nil(t, err)
	assert.EqualValues(t, 2, rowsInserted) // duplicate sku is ignored

	rowsDeleted, err := r.Delete(data)
	assert.Nil(t, err)
	assert.EqualValues(t, 2, rowsDeleted)

	// Persist and Delete with empty data
	data = make(map[string]value.Sku)

	rowsInserted, err = r.Persist(data)
	assert.Nil(t, err)
	assert.EqualValues(t, 0, rowsInserted)

	rowsDeleted, err = r.Delete(data)
	assert.Nil(t, err)
	assert.EqualValues(t, 0, rowsDeleted)
}
