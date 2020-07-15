package items

type Laptop struct {
	Processor Processor `json:"processor"`
	RAM       RAM       `json:"ram"`
	GPU       GPU       `json:"gpu"`
	Drive     Drive     `json:"drive"`
	Price     Price     `json:"price"`
}

type Processor struct {
	Series string `json:"series"`
	Model  string `json:"model"`
}

type RAM struct {
	Capacity int `json:"capacity"`
}

type GPU struct {
	Model string `json:"model"`
}

type Drive struct {
	Type     string `json:"type"`
	Capacity int    `json:"capacity"`
}

type Price struct {
	Min int `json:"min"`
	Max int `json:"max"`
}
