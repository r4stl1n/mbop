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
	req, _ := http.NewRequest("GET", url, nil)

	// add authorization header to the req
	req.Header.Add("Authorization", bearer)

	// Send req using http Client
	client := &http.Client{}
	resp, respError := client.Do(req)
	if respError != nil {
		zap.L().Error("failed to get response", zap.String("url", url), zap.Error(respError))
		return "", respError
	}

	defer func(Body io.ReadCloser) {
		bodyCloseError := Body.Close()
		if bodyCloseError != nil {
			zap.L().Error("failed to close body", zap.Error(bodyCloseError))
		}
	}(resp.Body)

	body, bodyError := io.ReadAll(resp.Body)

	if bodyError != nil {
		zap.L().Error("failed to read the response bytes:", zap.Error(bodyError))
		return "", bodyError
	}

	return string(body), nil
}

func (o *OpenAIAPI) postRequest(url string, data interface{}) (string, error) {

	// Create a Bearer string by appending string access token
	var bearer = "Bearer " + o.config.AuthToken

	marshall, marshallError := json.Marshal(data)

	if marshallError != nil {
		return "", marshallError
	}

	// Create a new request using http
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(marshall))

	// add authorization header to the req
	req.Header.Add("Authorization", bearer)

	// Send req using http Client
	client := &http.Client{}
	resp, respError := client.Do(req)
	if respError != nil {
		zap.L().Error("failed to get response", zap.String("url", url), zap.Error(respError))
		return "", respError
	}

	defer func(Body io.ReadCloser) {
		bodyCloseError := Body.Close()
		if bodyCloseError != nil {
			zap.L().Error("failed to close body", zap.Error(bodyCloseError))
		}
	}(resp.Body)

	body, bodyError := io.ReadAll(resp.Body)
	if bodyError != nil {
		zap.L().Error("failed to read the response bytes:", zap.Error(bodyError))
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

	response, responseError := o.getRequest(url)

	if responseError != nil {
		return Models{}, responseError
	}

	models := Models{}
	unmarshallError := json.Unmarshal([]byte(response), &models)

	if unmarshallError != nil {
		return Models{}, unmarshallError
	}

	return models, nil
}

func (o *OpenAIAPI) GetCompletion(completion CompletionHistory) (string, CompletionResponse, error) {
	url := o.config.BaseUrl + "/chat/completions"

	response, responseError := o.postRequest(url, completion.ToCompletionRequest())

	if responseError != nil {
		return "", CompletionResponse{}, responseError
	}

	completionResponse := CompletionResponse{}
	unmarshallError := json.Unmarshal([]byte(response), &completionResponse)

	if unmarshallError != nil {
		return "", CompletionResponse{}, unmarshallError
	}

	return completionResponse.Choices[0].Message.Content, completionResponse, nil

}
