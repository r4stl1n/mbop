package tools

type Tool interface {
	Name() string
	Example() string
	Description() string
	Run(values ...string) (string, error)
}

func ConvertToolArrayToPrompt(tools map[string]Tool) string {

	prompt := ""

	for _, tool := range tools {
		prompt += "\n"
		prompt += tool.Name() + "\n"
		prompt += tool.Example() + "\n"
		prompt += tool.Description() + "\n"
	}

	return prompt
}
