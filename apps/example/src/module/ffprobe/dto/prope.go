package dto

type ProbeQuery struct {
	Input string `query:"input"`
}

type FastConvertBody struct {
	Input   string         `json:"input"`
	Options map[string]any `json:"options"`
	Output  string         `json:"output"`
}
