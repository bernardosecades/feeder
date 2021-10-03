package value

import (
	"github.com/bernardosecades/feeder/pkg/tools/numberfn"
	"github.com/bernardosecades/feeder/pkg/tools/stringfn"

	"errors"
	"strings"
)

const separator = "-"

// All errors reported by the package
var (
	ErrSeparatorSku        = errors.New("wrong format, should be have two parts separate by '-' character")
	ErrLenFirstPartSku     = errors.New("first part of sku should have length 4")
	ErrLenSecondPartSku    = errors.New("second part of sku should have length 4")
	ErrLettersFirstPartSku = errors.New("first part only can content letters")
	ErrNumberSecondPartSku = errors.New("second part only can content unsigned integer")
)

type Sku struct {
	value string
}

// NewSku create new instance of Sku
func NewSku(v string) (Sku, error) {
	sku := Sku{value: v}
	err := sku.validateAndNormalize()
	if err != nil {
		return Sku{}, err
	}

	return sku, nil
}

// String return value from Sku
func (s Sku) String() string {
	return s.value
}

// StringWithoutZeros return value without zeros from Sku
func (s Sku) StringWithoutZeros() string {
	parts := strings.Split(s.value, "-")

	secondPart := strings.TrimLeft(parts[1], "0")
	if len(secondPart) == 0 {
		secondPart = "0"
	}

	return parts[0] + separator + secondPart
}

// validateAndNormalize it will clean (special characters), validate and normalize sku value
func (s *Sku) validateAndNormalize() error {
	// clean input
	s.value = strings.ReplaceAll(s.value, "\n", "")
	s.value = strings.ReplaceAll(s.value, "\r", "")
	s.value = strings.Trim(s.value, " ")

	parts := strings.Split(s.value, separator)
	if len(parts) != 2 {
		return ErrSeparatorSku
	}

	if len(parts[0]) != 4 {
		return ErrLenFirstPartSku
	}

	if len(parts[1]) != 4 {
		return ErrLenSecondPartSku
	}

	if !stringfn.IsOnlyLetters(parts[0]) {
		return ErrLettersFirstPartSku
	}

	if !numberfn.IsUnsignedInteger(parts[1]) {
		return ErrNumberSecondPartSku
	}

	// normalize
	s.value = strings.ToUpper(parts[0]) + separator + parts[1]

	return nil
}
