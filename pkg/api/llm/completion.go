package llm

type Completion struct {
	Model        string
	SystemPrompt string
	UserPrompt   string
}

func (c Completion) ToCompletionRequest() CompletionRequest {
	return CompletionRequest{
		Model: c.Model,
		Messages: []Message{
			{
				Role:    "system",
				Content: c.SystemPrompt,
			},
			{
				Role:    "user",
				Content: c.UserPrompt,
			},
		},
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
