// Package items provides models for extracted data
package items

// ModificationType represents category of a product
type ModificationType interface{}

// Product provides information for a product
type Product struct {
	Name          string             `json:"name"`
	Modifications []ModificationType `json:"modifications"`
}

// Laptop provides laptop specifications
type Laptop struct {
	CPU   `json:"cpu"`
	RAM   `json:"ram"`
	GPU   `json:"gpu"`
	Drive `json:"drive"`
	Price `json:"price"`
}

// CPU provides CPU specifications
type CPU struct {
	Series string `json:"series"`
	Model  string `json:"model"`
}

// RAM provides RAM specifications
type RAM struct {
	Capacity int `json:"capacity"`
}

// GPU provides GPU specifications
type GPU struct {
	Model string `json:"model"`
}

// Drive provides drive specifications
type Drive struct {
	Type     string `json:"type"`
	Capacity int    `json:"capacity"`
}

// Price provides min/max prices of the laptop
type Price struct {
	Min int `json:"min"`
	Max int `json:"max"`
}
