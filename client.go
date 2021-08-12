package xsolla

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
	"time"
)

type Client interface {
	NewTokenSettings() TokenSettings
	GetSubscriptionUserId(id string) (string, error)
	GetSubscription(id int) (*Subscription, error)
	GetUser(userId string) (*User, error)
	GetTransaction(id string) (*Transaction, error)
	CreateToken(token *Token) (string, error)
}

type StdClient struct {
	MerchantId     int
	MerchantSecret string
	ProjectId      int
	ProjectSecret  string
	Sandbox        bool
	Timeout        time.Duration
}

func (c *StdClient) doReq(req *http.Request, out interface{}) error {
	req.SetBasicAuth(strconv.Itoa(c.MerchantId), c.MerchantSecret)
	req.Header.Set("Accept", "application/json; charset=UTF-8")
	client := http.Client{
		Timeout: c.Timeout,
	}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	isJSON := strings.Contains(res.Header.Get("Content-Type"), "json")
	body, err := io.ReadAll(res.Body)
	// All of Xsollas non-error responses should be within the 200 range.
	// https://developers.xsolla.com/api/v2/getting-started/#api_errors_handling
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		err := RequestError{Raw: string(body)}
		if isJSON {
			json.Unmarshal(body, &err)
		}
		return &err
	} else if isJSON {
		return json.Unmarshal(body, out)
	}
	return nil
}

func (c *StdClient) newMerchantEndpoint(pathname string) *url.URL {
	u, err := url.Parse(EndpointMerchant)
	if err != nil {
		panic(err)
	}
	u.Path = path.Join(u.Path, fmt.Sprintf("%d", c.MerchantId), pathname)
	return u
}

func (c *StdClient) newProjectEndpoint(pathname string) *url.URL {
	u, err := url.Parse(EndpointProject)
	if err != nil {
		panic(err)
	}
	u.Path = path.Join(u.Path, fmt.Sprintf("%d", c.ProjectId), pathname)
	return u
}

// I am not sure about this style, I don't think project should be part of the main client and doing things like this.
// It would be better to find a better way to do newRequest by just passing the url.URL in directly instead. For now
// I am doing it like this to reduce current code verbosity.
//
// Secondly we are using panics because all request creations shouldn't fail during testing. If it panics there is
// something wrong with the codebase.
func (c *StdClient) newRequest(endpoint, method, pathname string, body io.Reader) *http.Request {
	var u *url.URL
	switch endpoint {
	case EndpointMerchant:
		u = c.newMerchantEndpoint(pathname)
	case EndpointProject:
		u = c.newProjectEndpoint(pathname)
	default:
		panic("no valid endpoint!")
	}
	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		panic(err)
	}
	return req
}

func (c *StdClient) newJSONRequest(endpoint, method, pathname string, body interface{}) (*http.Request, error) {
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req := c.newRequest(endpoint, method, pathname, bytes.NewReader(buf))
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

func (c *StdClient) NewTokenSettings() TokenSettings {
	var mode string
	if c.Sandbox {
		mode = modeSandbox
	}
	return TokenSettings{
		ProjectId: c.ProjectId,
		Mode:      mode,
		UI:        NewUISettings(),
	}
}

func (c *StdClient) GetSubscriptionUserId(id string) (string, error) {
	req := c.newRequest(EndpointMerchant, http.MethodGet, fmt.Sprintf("subscriptions/%s", id), nil)
	var resPayload struct {
		User struct {
			Id string `json:"id"`
		} `json:"user"`
	}
	err := c.doReq(req, &resPayload)
	return resPayload.User.Id, err
}

func (c *StdClient) GetSubscription(id int) (*Subscription, error) {
	req := c.newRequest(EndpointProject, http.MethodGet, fmt.Sprintf("subscriptions/%d", id), nil)
	var resPayload Subscription
	if err := c.doReq(req, &resPayload); err != nil {
		return nil, err
	}
	return &resPayload, nil
}

func (c *StdClient) GetUser(userId string) (*User, error) {
	req := c.newRequest(EndpointProject, http.MethodGet, fmt.Sprintf("users/%s", userId), nil)
	var resPayload User
	if err := c.doReq(req, &resPayload); err != nil {
		return nil, err
	}
	return &resPayload, nil
}

func (c *StdClient) GetTransaction(id string) (*Transaction, error) {
	req := c.newRequest(EndpointMerchant, http.MethodGet, fmt.Sprintf("reports/transactions/%s/details", id), nil)
	var resPayload Transaction
	if err := c.doReq(req, &resPayload); err != nil {
		return nil, err
	}
	return &resPayload, nil
}

func (c *StdClient) CreateToken(token *Token) (string, error) {
	req, err := c.newJSONRequest(EndpointMerchant, http.MethodPost, "token", token)
	if err != nil {
		return "", err
	}
	var resPayload struct {
		Token string `json:"token"`
	}
	err = c.doReq(req, &resPayload)
	return resPayload.Token, err
}
