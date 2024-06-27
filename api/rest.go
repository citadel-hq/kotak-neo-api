package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type ApiException struct {
	Status int
	Reason string
	Body   string
}

func (e *ApiException) Error() string {
	return fmt.Sprintf("ApiException: Status=%d, Reason=%s, Body=%s", e.Status, e.Reason, e.Body)
}

type RESTClientObject struct {
	Configuration Configuration
}

func NewRESTClientObject(configuration Configuration) *RESTClientObject {
	return &RESTClientObject{
		Configuration: configuration,
	}
}

func (c *RESTClientObject) Request(method, urlStr string, queryParams map[string]string, headers map[string]string, body interface{}) (*http.Response, error) {
	method = strings.ToUpper(method)
	if method != "GET" && method != "HEAD" && method != "DELETE" && method != "POST" && method != "PUT" && method != "PATCH" && method != "OPTIONS" {
		return nil, &ApiException{Status: 0, Reason: "Invalid HTTP method"}
	}

	if headers == nil {
		headers = map[string]string{}
	}

	if _, exists := headers["Content-Type"]; !exists {
		headers["Content-Type"] = "application/json"
	}

	var reqBody []byte
	var err error

	if body != nil {
		if strings.Contains(strings.ToLower(headers["Content-Type"]), "json") {
			reqBody, err = json.Marshal(body)
			if err != nil {
				return nil, err
			}
		} else if strings.Contains(strings.ToLower(headers["Content-Type"]), "x-www-form-urlencoded") {
			form := url.Values{}
			form.Add("jData", fmt.Sprintf("%v", body))
			reqBody = []byte(form.Encode())
		} else {
			return nil, &ApiException{Status: 0, Reason: "Invalid Content-Type in the Header Parameters"}
		}
	}

	if queryParams != nil {
		params := url.Values{}
		for k, v := range queryParams {
			params.Add(k, v)
		}
		urlStr += "?" + params.Encode()
	}

	req, err := http.NewRequest(method, urlStr, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, &ApiException{Status: 0, Reason: fmt.Sprintf("%T\n%s", err, err.Error())}
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		return nil, &ApiException{Status: resp.StatusCode, Reason: resp.Status, Body: string(bodyBytes)}
	}

	return resp, nil
}
