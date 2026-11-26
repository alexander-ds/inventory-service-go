package inventory

import (
	"errors"
	"testing"
)

func TestProductInputValidate(t *testing.T) {
	valid := ProductInput{SKU: "LAP-001", Name: "Laptop", Unit: "UND"}
	if err := valid.Validate(); err != nil {
		t.Fatalf("payload válido rechazado: %v", err)
	}

	cases := map[string]ProductInput{
		"sin sku":           {Name: "Laptop", Unit: "UND"},
		"sin name":          {SKU: "LAP-001", Unit: "UND"},
		"sin unit":          {SKU: "LAP-001", Name: "Laptop"},
		"precio negativo":   {SKU: "LAP-001", Name: "Laptop", Unit: "UND", PurchasePrice: -1},
		"precio venta negativo": {SKU: "LAP-001", Name: "Laptop", Unit: "UND", SalePrice: -1},
	}

	for name, input := range cases {
		t.Run(name, func(t *testing.T) {
			if err := input.Validate(); !errors.Is(err, ErrInvalidInput) {
				t.Errorf("esperaba ErrInvalidInput, obtuve %v", err)
			}
		})
	}
}
