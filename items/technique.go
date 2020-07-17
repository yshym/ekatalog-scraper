package items

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
