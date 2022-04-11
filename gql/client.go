package gql

import (
	"agora/assignments/agora"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"strings"
)

const (
	appIDHeader       = "app-id"
	ContentTypeJSON   = "application/json"
	ContentTypeHeader = "Content-Type"

	// pagination
	PageQueryParam  = "page"
	limitQueryParam = "Limit"

	usersResource = "user"
)

// Client is a fully-featured client through which GraphQL Users API calls can be made.
type Client struct {
	url        string
	httpClient *http.Client
	appID      string
}

// GetUsers retrives all list of users.
func (c *Client) GetUsers(ctx context.Context, page, perPage int64) ([]agora.User, error) {
	var users struct {
		Data []agora.User
	}

	var queryParams = map[string]string{
		PageQueryParam:  fmt.Sprintf("%d", page),
		limitQueryParam: fmt.Sprintf("%d", perPage),
	}

	var _, err = c.makeRequest(
		ctx,
		http.MethodGet,
		usersResource,
		queryParams,
		nil,
		&users,
	)
	if err != nil {
		return nil, err
	}

	return users.Data, nil
}

// GetUserByID retrives user by given unique ID.
func (c *Client) GetUserByID(ctx context.Context, id string) (agora.User, error) {
	var user agora.User

	var _, err = c.makeRequest(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/%s", usersResource, id),
		nil,
		nil,
		&user,
	)
	if err != nil {
		return agora.User{}, err
	}

	return user, nil
}

// makeRequest is a common HTTP request execution method.
func (c *Client) makeRequest(ctx context.Context,
	method, uri string,
	params map[string]string,
	in, out interface{}) (*http.Response, error) {
	var body io.Reader

	// Marshal the request body, if one was provided.
	if in != nil {
		var bodyData, err = json.Marshal(in)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}

		body = bytes.NewReader(bodyData)
	}

	// Build the request.
	var req, err = http.NewRequest(method, fmt.Sprintf("%s/%s", c.url, uri), body)
	if err != nil {
		return nil, fmt.Errorf("failed to create new HTTP request: %w", err)
	}

	if params != nil {
		var queryParam = req.URL.Query()

		for key, value := range params {
			queryParam.Add(key, value)
		}
		req.URL.RawQuery = queryParam.Encode()
	}

	if body != nil {
		req.Header.Set(ContentTypeHeader, ContentTypeJSON)
	}
	req.Header.Set(appIDHeader, c.appID)

	// Execute the request.
	var resp *http.Response
	resp, err = c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute HTTP request: %w", err)
	}
	defer ConsumeAndCloseResponseBody(resp)

	// Return an appropriate error if the HTTP status code indicates failure.
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, errors.New("upstream server error")
	}

	// Unmarshal the response body, if a response body is expected.
	err = ParseResponseBody(resp, out)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTTP response body: %w", err)
	}

	return resp, nil
}

// ConsumeAndCloseResponseBody fully consumes and then closes a HTTP response
// body.
func ConsumeAndCloseResponseBody(r *http.Response) {
	defer r.Body.Close()
	io.Copy(ioutil.Discard, r.Body)
}

// ParseResponseBody parses an HTTP response body according to its content
// type and parses in into out. If out is nil, no action is performed.
// Currently, only an application/json content type is supported.
func ParseResponseBody(r *http.Response, out interface{}) error {
	if out == nil {
		return nil
	}

	var mediaType, _, err = mime.ParseMediaType(r.Header.Get(ContentTypeHeader))
	if err != nil {
		return fmt.Errorf("failed to parse HTTP response media type: %w", err)
	}

	if !strings.HasPrefix(mediaType, ContentTypeJSON) {
		return fmt.Errorf("unexpected content type: %q", mediaType)
	}

	var data []byte
	data, err = ioutil.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("failed to read HTTP response body: %w", err)
	}

	err = json.Unmarshal(data, out)
	if err != nil {
		return fmt.Errorf("failed to unmarshal HTTP response body: %w", err)
	}

	return nil
}

// newClient creates a new users(graphQL API) client.
func newClient(url string, appID string) (*Client, error) {
	return &Client{
		url:   url,
		appID: appID,
		httpClient: &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
			},
		},
	}, nil
}
