package http

import (
    "net/http"
)

const NPM_REGISTRY_URL = "https://registry.npmjs.org"

type HttpRequest struct {
    Url string
    Method string
    Header map[string]string
}

// Initializes the request data.
func CreateRequest(name string, version string) *HttpRequest {
    r := &HttpRequest{}

    r.Url = NPM_REGISTRY_URL + "/" + name + "/" + version
    r.Method = "GET"

    r.Header = map[string]string{
        "Accept": "application/vnd.npm.install-v1+json; q=1.0, application/json; q=0.8, */*",
    }

    return r
}

// Send a request to the npm registry
func (r *HttpRequest) Send() (*http.Response, error) {
    
    req, err := http.NewRequest(r.Method, r.Url, nil)
    if err != nil {
        return nil, err
    }

    for k, v := range r.Header {
        req.Header.Set(k, v)
    }

    client := &http.Client{}

    response, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    
    return response, nil
}




