package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

const (
	ApiVersion = "2.251"
	AuthURL    = "/ipa/session/login_password"
	JsonURL    = "/ipa/session/json"
)

type JSON = map[string]interface{}

type APIResponse[RT any, VT any] struct {
	Result struct {
		Result  RT     `json:"result"`
		Summary string `json:"summary"`
		Value   VT     `json:"value"`
	} `json:"result"`
	Error     *APIError `json:"error"`
	ID        int       `json:"id"`
	Principal string    `json:"principal"`
	Version   string    `json:"version"`
}

type APIClient struct {
	httpClient *http.Client

	server string
}

func NewClient(server string, user string, password string) (*APIClient, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	client := &APIClient{
		httpClient: &http.Client{
			Jar: jar,
		},
		server: server,
	}

	var AuthData = url.Values{}
	AuthData.Set("user", user)
	AuthData.Set("password", password)

	req, err := http.NewRequest("POST", client.server+AuthURL, strings.NewReader(AuthData.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "text/plain")
	req.Header.Set("Referer", client.server+"/ipa/")

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	client.httpClient.Jar.SetCookies(req.URL, resp.Cookies())

	return client, nil
}

func apiRequest[RT any, VT any](c *APIClient, method string, options JSON, params ...string) (*RT, error) {
	// Add version to options (required)
	if options == nil {
		options = JSON{}
	}
	options["version"] = ApiVersion

	if params == nil {
		params = make([]string, 0)
	}

	var JSONReqData = map[string]interface{}{
		"method": method,
		"params": []interface{}{
			params,
			options,
		},
		"id": 0,
	}
	reqData, err := json.Marshal(JSONReqData)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", c.server+JsonURL, bytes.NewBuffer(reqData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Referer", c.server+"/ipa/")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("request failed with status " + resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var APIresp APIResponse[RT, VT]
	err = json.Unmarshal(body, &APIresp)
	if err != nil {
		return nil, err
	}

	if APIresp.Error != nil {
		return nil, APIresp.Error
	}

	return &APIresp.Result.Result, nil
}
