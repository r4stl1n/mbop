package llm

import (
	"bytes"
	"encoding/json"
	"go.uber.org/zap"
	"io"
	"net/http"
)

type OpenAIAPI struct {
	config *OpenAIConfig
}

func (o *OpenAIAPI) Init() (*OpenAIAPI, error) {

	config, configError := new(OpenAIConfig).Init()

	if configError != nil {
		return nil, configError
	}

	*o = OpenAIAPI{
		config: config,
	}

	return o, nil
}

func (o *OpenAIAPI) getRequest(url string) (string, error) {
	// Create a Bearer string by appending string access token
	var bearer = "Bearer " + o.config.AuthToken

	// Create a new request using http
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		zap.L().Error("openai http request creation failed", 
			zap.String("url", url),
			zap.String("method", "GET"),
			zap.Error(err))
		return "", err
	}

	// add authorization header to the req
	req.Header.Add("Authorization", bearer)
	req.Header.Add("Content-Type", "application/json")

	// Send req using http Client
	client := &http.Client{}
	resp, respError := client.Do(req)
	if respError != nil {
		zap.L().Error("openai http request failed", 
			zap.String("url", url),
			zap.String("method", "GET"),
			zap.Error(respError))
		return "", respError
	}

	defer func(Body io.ReadCloser) {
		bodyCloseError := Body.Close()
		if bodyCloseError != nil {
			zap.L().Error("openai response body close failed", zap.Error(bodyCloseError))
		}
	}(resp.Body)

	body, bodyError := io.ReadAll(resp.Body)
	if bodyError != nil {
		zap.L().Error("openai response body read failed", 
			zap.String("url", url),
			zap.Error(bodyError))
		return "", bodyError
	}

	return string(body), nil
}

func (o *OpenAIAPI) postRequest(url string, data interface{}) (string, error) {
	// Create a Bearer string by appending string access token
	var bearer = "Bearer " + o.config.AuthToken

	marshall, marshallError := json.Marshal(data)
	if marshallError != nil {
		zap.L().Error("openai request JSON marshalling failed", 
			zap.String("url", url),
			zap.Error(marshallError))
		return "", marshallError
	}

	// Create a new request using http
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(marshall))
	if err != nil {
		zap.L().Error("openai http request creation failed", 
			zap.String("url", url),
			zap.String("method", "POST"),
			zap.Error(err))
		return "", err
	}

	// add authorization header to the req
	req.Header.Add("Authorization", bearer)
	req.Header.Add("Content-Type", "application/json")

	// Send req using http Client
	client := &http.Client{}
	resp, respError := client.Do(req)
	if respError != nil {
		zap.L().Error("openai http request failed", 
			zap.String("url", url),
			zap.String("method", "POST"),
			zap.Error(respError))
		return "", respError
	}

	defer func(Body io.ReadCloser) {
		bodyCloseError := Body.Close()
		if bodyCloseError != nil {
			zap.L().Error("openai response body close failed", zap.Error(bodyCloseError))
		}
	}(resp.Body)

	body, bodyError := io.ReadAll(resp.Body)
	if bodyError != nil {
		zap.L().Error("openai response body read failed", 
			zap.String("url", url),
			zap.Error(bodyError))
		return "", bodyError
	}

	return string(body), nil
}

func (o *OpenAIAPI) TestConnection() error {
	_, err := o.GetModels()

	return err
}

func (o *OpenAIAPI) GetModels() (Models, error) {
	url := o.config.BaseUrl + "/models"

	zap.L().Info("retrieving openai models", zap.String("url", url))

	response, responseError := o.getRequest(url)
	if responseError != nil {
		zap.L().Error("openai models retrieval failed", zap.Error(responseError))
		return Models{}, responseError
	}

	models := Models{}
	unmarshallError := json.Unmarshal([]byte(response), &models)
	if unmarshallError != nil {
		zap.L().Error("openai models JSON parsing failed", zap.Error(unmarshallError))
		return Models{}, unmarshallError
	}

	zap.L().Info("openai models retrieved", zap.Int("modelCount", len(models.Data)))
	return models, nil
}

func (o *OpenAIAPI) GetCompletion(completion CompletionHistory) (string, CompletionResponse, error) {
	url := o.config.BaseUrl + "/chat/completions"

	zap.L().Info("requesting openai completion", 
		zap.String("model", completion.Model),
		zap.Int("messageCount", len(completion.Context)))

	response, responseError := o.postRequest(url, completion.ToCompletionRequest())
	if responseError != nil {
		zap.L().Error("openai completion request failed", 
			zap.String("model", completion.Model),
			zap.Error(responseError))
		return "", CompletionResponse{}, responseError
	}

	completionResponse := CompletionResponse{}
	unmarshallError := json.Unmarshal([]byte(response), &completionResponse)
	if unmarshallError != nil {
		zap.L().Error("openai completion JSON parsing failed", zap.Error(unmarshallError))
		return "", CompletionResponse{}, unmarshallError
	}

	if len(completionResponse.Choices) == 0 {
		zap.L().Error("openai completion returned no choices")
		return "", completionResponse, nil
	}

	content := completionResponse.Choices[0].Message.Content
	zap.L().Info("openai completion received", 
		zap.String("model", completion.Model),
		zap.Int("contentLength", len(content)),
		zap.Int("tokenUsage", completionResponse.Usage.TotalTokens))

	return content, completionResponse, nil
}
