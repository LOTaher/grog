package httpclient

import (
	"fmt"
	"net/http"
)

const NPM_REGISTRY_URL = "https://registry.npmjs.org"

func SendRequest(name, version string) (*http.Response, error) {
	url := fmt.Sprintf("%s/%s/%s", NPM_REGISTRY_URL, name, version)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %v", err)
	}

	req.Header.Set("Accept", "application/vnd.npm.install-v1+json; q=1.0, application/json; q=0.8, */*")

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send HTTP request: %v", err)
	}

	return response, nil
}
