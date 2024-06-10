package llm

import "fmt"

type CompletionHistory struct {
	Model   string
	Context []Message
}

func (c *CompletionHistory) Init(model string) *CompletionHistory {

	*c = CompletionHistory{
		Model:   model,
		Context: []Message{},
	}

	c.Add(Message{
		Role:    "system",
		Content: "",
	})

	return c
}

func (c *CompletionHistory) PrintHistory() string {
	output := ""

	for _, x := range c.Context {
		output = output + fmt.Sprintf("Role: %s\nContent: %s\n\n", x.Role, x.Content)
	}

	return output
}

func (c *CompletionHistory) PrintLatestHistory() string {
	output := ""

	x := c.Context[len(c.Context)-1]

	output = output + fmt.Sprintf("Role: %s\nContent: %s\n\n", x.Role, x.Content)

	return output
}

func (c *CompletionHistory) Add(message Message) {
	c.Context = append(c.Context, message)
}

func (c *CompletionHistory) RemoveLatest() {
	c.Context = c.Context[:len(c.Context)-1]
}

func (c *CompletionHistory) ToCompletionRequest() CompletionRequest {
	return CompletionRequest{
		Model:    c.Model,
		Messages: c.Context,
		Length:   8192,
	}
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Choice struct {
	Index        int         `json:"index"`
	Message      Message     `json:"message"`
	Logprobs     interface{} `json:"logprobs"`
	FinishReason string      `json:"finish_reason"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type CompletionRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Length   int       `json:"max_tokens"`
}

type CompletionResponse struct {
	Id                string   `json:"id"`
	Object            string   `json:"object"`
	Created           int      `json:"created"`
	Model             string   `json:"model"`
	SystemFingerprint string   `json:"system_fingerprint"`
	Choices           []Choice `json:"choices"`
	Usage             Usage    `json:"usage"`
}
