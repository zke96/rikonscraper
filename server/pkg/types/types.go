package types

import "github.com/google/uuid"

type Product struct {
	ID          uuid.UUID `json:"id,omitempty"`
	Display     string    `json:"label,omitempty"`
	ProductCode string    `json:"productCode,omitempty" sql:"product_code"`
	Price       string    `json:"price,omitempty" db:"-"`
	InStock     bool      `json:"inStock,omitempty" db:"-"`
	Description string    `json:"description,omitempty" db:"-"`
	URL         string    `json:"url,omitempty"`
}

type AlertRecord struct {
	ID      uuid.UUID `json:"id,omitempty"`
	PartID  uuid.UUID `json:"partID,omitempty"`
	Email   string    `json:"email,omitempty"`
	Display string    `json:"display,omitempty"`
	URL     string    `json:"url,omitempty"`
}
