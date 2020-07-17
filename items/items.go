// Package items provides models for extracted data
package items

// ModificationType represents category of a product
type ModificationType interface{}

// Product provides information for a product
type Product struct {
	Name          string             `json:"name"`
	Modifications []ModificationType `json:"modifications"`
}
