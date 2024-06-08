package llm

type Model struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	OwnedBy string `json:"owned_by"`
}

type Models struct {
	Object string  `json:"object"`
	Data   []Model `json:"data"`
}
