package items

// Smartphone provides smartphone specifications
type Smartphone struct {
	CPU   `json:"cpu"`
	RAM   `json:"ram"`
	GPU   `json:"gpu"`
	Drive `json:"drive"`
	Price `json:"price"`
}
