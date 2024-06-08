package llm

import (
	"fmt"
	"os"
)

type OpenAIConfig struct {
	BaseUrl   string
	AuthToken string
}

func (o *OpenAIConfig) Init() (*OpenAIConfig, error) {

	*o = OpenAIConfig{}

	return o.fromEnv()
}

func (o *OpenAIConfig) fromEnv() (*OpenAIConfig, error) {

	baseUrl, baseUrlOk := os.LookupEnv("OPENAI_BASE_URL")

	if !baseUrlOk {
		return nil, fmt.Errorf("env variable OPENAI_BASE_URL required")
	}

	authToken, authTokenOk := os.LookupEnv("OPENAI_AUTH_TOKEN")

	if !authTokenOk {
		return nil, fmt.Errorf("env variable OPENAI_AUTH_TOKEN required")
	}

	*o = OpenAIConfig{
		BaseUrl:   baseUrl,
		AuthToken: authToken,
	}

	return o, nil
}
