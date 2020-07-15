// Package items provides models for extracted data
package items

// Laptop provides laptop specifications
type Laptop struct {
	Processor Processor `json:"processor"`
	RAM       RAM       `json:"ram"`
	GPU       GPU       `json:"gpu"`
	Drive     Drive     `json:"drive"`
	Price     Price     `json:"price"`
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
