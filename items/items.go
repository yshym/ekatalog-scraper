// Package items provides models for extracted data
package items

// Product provides information for a product
type Product struct {
	Name          string        `json:"name"`
	Modifications []interface{} `json:"modifications"`
}

// Laptop provides laptop specifications
type Laptop struct {
	Processor `json:"processor"`
	RAM       `json:"ram"`
	GPU       `json:"gpu"`
	Drive     `json:"drive"`
	Price     `json:"price"`
}

// Processor provides processor specifications
type Processor struct {
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
