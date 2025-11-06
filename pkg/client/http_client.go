package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/mangudaigb/dhauli-base/logger"
)

type AuthCredentials struct {
	AccessToken  string
	RefreshToken string
	RefreshUrl   string
}

type TokenRefresher func(refreshToken string) (string, error)

type HttpClient[T any] struct {
	log     *logger.Logger
	Auth    AuthCredentials
	client  *http.Client
	baseUrl string
	IsSsl   bool
}

type HttpClientOption[T any] func(*HttpClient[T])

func WithAuth[T any](auth AuthCredentials) HttpClientOption[T] {
	return func(c *HttpClient[T]) {
		c.Auth = auth
	}
}

func WithSSL[T any](ssl bool) HttpClientOption[T] {
	return func(c *HttpClient[T]) {
		c.IsSsl = ssl
	}
}

func (sc *HttpClient[T]) Get(path string, params map[string]string) (*T, error) {
	return sc.doSimpleRequest(http.MethodGet, path, params)
}

func (sc *HttpClient[T]) Post(path string, body *T, optionalParams map[string]string) (*T, error) {
	return sc.doBodyRequest(http.MethodPost, path, body, optionalParams)
}

func (sc *HttpClient[T]) Put(path string, body *T, optionalParams map[string]string) (*T, error) {
	return sc.doBodyRequest(http.MethodPut, path, body, optionalParams)
}

func (sc *HttpClient[T]) Delete(path string, optionalParams map[string]string) (*T, error) {
	return sc.doSimpleRequest(http.MethodDelete, path, optionalParams)
}

func (sc *HttpClient[T]) doSimpleRequest(httpMethod, path string, optionalParams map[string]string) (*T, error) {
	fullUrl := sc.buildUrl(path, optionalParams)
	sc.log.Info("%s %s", httpMethod, fullUrl)

	response, err := sc.client.Get(fullUrl)
	if err != nil {
		sc.log.Errorf("Error while getting %s by error: %v", fullUrl, err)
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			sc.log.Errorf("Error while closing response body: %v", err)
		}
	}(response.Body)

	body, err := io.ReadAll(response.Body)
	if err != nil {
		sc.log.Errorf("Error while reading response body: %v", err)
		return nil, err
	}
	sc.log.Info("Response body: %s", string(body))

	var result T
	if err = json.Unmarshal(body, &result); err != nil {
		sc.log.Errorf("Error while unmarshalling response body: %v", err)
		return nil, err
	}
	return &result, nil
}

func (sc *HttpClient[T]) doBodyRequest(httpMethod, path string, body *T, requestParams map[string]string) (*T, error) {
	fullUrl := sc.buildUrl(path, requestParams)
	sc.log.Info("%s %s", httpMethod, fullUrl)
	bodyJson, err := json.Marshal(body)
	if err != nil {
		sc.log.Errorf("Error while marshalling body: %v", err)
		return nil, err
	}
	sc.log.Info("Body: %s", string(bodyJson))

	req, err := http.NewRequest(httpMethod, fullUrl, io.NopCloser(bytes.NewReader(bodyJson)))
	if err != nil {
		sc.log.Errorf("Error while creating request: %v", err)
		return nil, fmt.Errorf("error while creating %s request: %v", httpMethod, err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := sc.client.Do(req)
	if err != nil {
		sc.log.Errorf("Error while sending request: %v", err)
		return nil, fmt.Errorf("error while executing %s request: %v", httpMethod, err)
	}
	defer func(Body io.ReadCloser) {
		err := resp.Body.Close()
		if err != nil {
			sc.log.Errorf("Error while closing response body: %v", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		respBodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code: %d, response body: %s", resp.StatusCode, string(respBodyBytes))
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		sc.log.Errorf("Error while reading response body: %v", err)
		return nil, err
	}
	var result T
	if err = json.Unmarshal(responseBody, &result); err != nil {
		sc.log.Errorf("Error while unmarshalling response body: %v", err)
		return nil, fmt.Errorf("error unmarshalling response body into type %T: %w", result, err)
	}
	return &result, nil
}

func (sc *HttpClient[T]) buildUrl(path string, params map[string]string) string {
	baseUrl, _ := url.Parse(sc.baseUrl)
	fullUrl, _ := baseUrl.Parse(path)
	query := fullUrl.Query()
	for k, v := range params {
		query.Add(k, v)
	}
	fullUrl.RawQuery = query.Encode()
	return fullUrl.String()
}

func NewHttpClient[T any](log *logger.Logger, baseUrl string, opts ...HttpClientOption[T]) *HttpClient[T] {
	stateClient := &HttpClient[T]{
		log: log,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		baseUrl: baseUrl,
	}

	for _, opt := range opts {
		opt(stateClient)
	}
	return stateClient
}
