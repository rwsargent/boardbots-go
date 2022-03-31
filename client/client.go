// An implementaiton of the public API for boardbots.dev.
// BoardBotClient wraps an HTTPClient
package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"time"

	"golang.org/x/net/publicsuffix"
)

type (
	Client interface {
		Authenticate() error
		MakeMove() error
		StartGame() error
		JoinLobby() error
		CreateLobby() error
	}

	Credentials struct {
		Username string
	}
	BoardBotClient[State any] struct {
		Credentials Credentials
		httpClient  *http.Client
		domain      string
		user        User
	}

	User struct {
		Id   json.Number `json:"id"`
		Name string      `json:"name"`
	}
)

// NewBoardBotClient builds a BoardBotClient with the users credentials, and server address.
// Parameterized with the type of Game (via its state). You'll need to instantiate
// multiple clients if you want to manipulate different types of games.
func NewBoardBotClient[State any](creds Credentials, addr string) (*BoardBotClient[State], error) {
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		return nil, err
	}

	return &BoardBotClient[State]{
		Credentials: creds,
		httpClient: &http.Client{
			Jar:     jar,
			Timeout: time.Duration(2) * time.Second,
		},
		domain: addr,
	}, nil
}

// Authenticates the client. Stores auth cookie automatically, makes
// Another call to get user information.
func (c *BoardBotClient[S]) Authenticate() error {
	resp, err := c.httpClient.Get(c.domain + "/auth/login?name=" + c.Credentials.Username)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("status code: %d", resp.StatusCode)
	}
	resp, err = c.httpClient.Get(c.domain + "/user")
	if err != nil {
		return err
	}
	err = json.NewDecoder(resp.Body).Decode(&c.user)
	defer resp.Body.Close()

	if err != nil {
		return err
	}

	return nil
}

// Post makes an HTTP post call with the supplied request, and Unmarshals
// the response from JSON to the supplied request struct.
func Post[Req any, Res any, S any](bbClient *BoardBotClient[S], path string, body Req) (Res, error) {
	b, err := json.Marshal(body)
	var resp Res
	if err != nil {
		return resp, err
	}

	req, err := http.NewRequest("POST", bbClient.domain+path, bytes.NewBuffer(b))
	if err != nil {
		return resp, err
	}

	httpResponse, err := bbClient.httpClient.Do(req)
	if err != nil {
		return resp, err
	}

	responseBody, _ := ioutil.ReadAll(httpResponse.Body)
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode != 200 {
		return resp, fmt.Errorf("response status: %d. Response body: %s", httpResponse.StatusCode, string(responseBody))
	}

	if len(responseBody) == 0 {
		return resp, nil
	}

	if _, ok := any(resp).(string); ok {
		return resp, nil
	}

	err = json.Unmarshal(responseBody, &resp)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

// Get makes an HTTP Get call, and Unmarshals the response to the provided struct.
func Get[Res any, S any](client *BoardBotClient[S], path string) (Res, error) {
	body, err := GetString(client, path)

	var resp Res
	if err != nil {
		return resp, nil
	}

	err = json.Unmarshal([]byte(body), &resp)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

// GetString makes an HTTP GET call, and returns the body of the response
// as a string
func GetString[S any](client *BoardBotClient[S], path string) (string, error) {
	httpResp, err := client.httpClient.Get(client.domain + path)
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(httpResp.Body)
	defer httpResp.Body.Close()

	if httpResp.StatusCode != 200 {
		return "", fmt.Errorf("status code: %d, %s", httpResp.StatusCode, body)
	}

	if err != nil {
		return "", err
	}
	return string(body), err
}
