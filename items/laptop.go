package items

// Laptop provides laptop specifications
type Laptop struct {
	CPU   `json:"cpu"`
	RAM   `json:"ram"`
	GPU   `json:"gpu"`
	Drive `json:"drive"`
	Price `json:"price"`
}
