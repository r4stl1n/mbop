package wiki

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/http"
)

type Wikipedia struct {
}

func (w Wikipedia) Name() string {
	return "wikipedia"
}

func (w Wikipedia) Example() string {
	return "wikipedia: Django"
}

func (w Wikipedia) Description() string {
	return "Returns a summary from searching Wikipedia"
}

func (w Wikipedia) Run(values ...string) (string, error) {

	if len(values) != 1 {
		return "", fmt.Errorf("expected one argument")
	}

	url := "https://en.wikipedia.org/w/api.php"

	// Create a new request using http
	req, _ := http.NewRequest("GET", url, nil)

	q := req.URL.Query()
	q.Add("action", "query")
	q.Add("prop", "extracts")
	q.Add("titles", values[0])
	q.Add("format", "json")
	q.Add("exintro", "")
	q.Add("explaintext", "")
	q.Add("redirects", "1")

	req.URL.RawQuery = q.Encode()

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

	response := WikipediaResponse{}
	unmarshallError := json.Unmarshal(body, &response)

	if unmarshallError != nil {
		return "", unmarshallError
	}

	if len(response.Query.Pages) < 1 {
		return "", fmt.Errorf("no results found")
	}

	// Return first result
	for _, v := range response.Query.Pages {
		return v.Extract, nil
	}

	return "", fmt.Errorf("no results found")
}
