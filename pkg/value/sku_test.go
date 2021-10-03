package value_test

import (
	"github.com/bernardosecades/feeder/pkg/value"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidSku(t *testing.T) {
	cases := []struct {
		skuInput    string
		skuExpected string
	}{
		{"KASL-3423",  "KASL-3423"},
		{"KASL-3423 ",  "KASL-3423"},
		{" KASL-3423 ",  "KASL-3423"},
		{" KASL-3423 \n",  "KASL-3423"},
		{"\n KASL-3423 \n\r",  "KASL-3423"},
		{"kasl-3423",  "KASL-3423"},
	}

	for _, c := range cases {
		sku, err := value.NewSku(c.skuInput)
		assert.Nil(t, err)
		assert.Equal(t, c.skuExpected, sku.String())
	}
}

func TestInvalidSku(t *testing.T) {
	_, err := value.NewSku("AAAA*1234")
	assert.Equal(t, err, value.ErrSeparatorSku)

	_, err = value.NewSku("A2AA-1234")
	assert.Equal(t, err, value.ErrLettersFirstPartSku)

	_, err = value.NewSku("AAAA-1.34")
	assert.Equal(t, err, value.ErrNumberSecondPartSku)

	_, err = value.NewSku("AAA-12345")
	assert.Equal(t, err, value.ErrLenFirstPartSku)

	_, err = value.NewSku("AAAA-123")
	assert.Equal(t, err, value.ErrLenSecondPartSku)
}

func TestSkuWithoutZeros(t *testing.T) {
	cases := []struct {
		skuInput    string
		skuExpected string
	}{
		{"KASL-0003",  "KASL-3"},
		{"KASL-0300",  "KASL-300"},
		{"KASL-0000",  "KASL-0"},
	}

	for _, c := range cases {
		sku, err := value.NewSku(c.skuInput)
		assert.Nil(t, err)
		assert.Equal(t, c.skuExpected, sku.StringWithoutZeros())
	}
}
