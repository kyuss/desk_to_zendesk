package zendesk

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

const UserAgentID = "Filestack-go 0.1"

var ErrMissingTicket = errors.New("no ticket found")
var ErrMissingEmail = errors.New("no email found")

const ZendeskURL = "https://filestack.zendesk.com/api/v2"

type Client struct {
	username string
	password string
	http     *http.Client
}

func NewClient(username, password string) *Client {

	timeout := time.Duration(120 * time.Second)

	return &Client{
		username: username,
		password: password,
		http:     &http.Client{Timeout: timeout},
	}
}

func (client *Client) do(method, urlStr, bodyType string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, urlStr, body)
	if err != nil {
		return nil, err
	}
	if bodyType != "" {
		req.Header.Set("Content-Type", bodyType)
	}
	req.Header.Set("User-Agent", UserAgentID)
	req.Header.Add("Accept-Encoding", "identity")
	req.SetBasicAuth(client.username, client.password)

	return client.http.Do(req)
}

func readError(resp *http.Response) error {
	if resp == nil {
		return errors.New("Invalid response")
	}
	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
		return nil
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	zerror := ZendeskError{
		Code: resp.StatusCode,
	}
	err = json.Unmarshal(bytes, &zerror)
	if err != nil {
		return err
	}

	return zerror
}

// ZendeskError represents an error that can be returned from filepicker.io service.
type ZendeskError struct {
	Code        int
	Message     string                   `json:"error"`
	Description string                   `json:"description"`
	Details     map[string][]ErrorDetail `json:"details"`
}

type ErrorDetail struct {
	Description string `json:"description"`
	Error       string `json:"error"`
}

func (err *ZendeskError) Includes(key string, value string) bool {
	if details, ok := err.Details[key]; ok {
		for _, detail := range details {
			if detail.Error == value {
				return true
			}
		}
	}

	return false
}

// Error satisfies builtin.error interface. It prints an error string with
// the reason of failure.
func (e ZendeskError) Error() string {
	return fmt.Sprintf("Zendesk: %d - %s, %s, %s", e.Code, e.Message, e.Description, e.Details)
}
