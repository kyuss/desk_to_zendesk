package desk

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const UserAgentID = "Filestack-go 0.1"

const DeskURL = "https://xenon.desk.com/api/v2"
const DeskLabel = "Filepicker"

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
	if resp.StatusCode == http.StatusOK {
		return nil
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return DeskError{
		Code:    resp.StatusCode,
		Message: strings.TrimSpace(string(bytes)),
	}
}

// DeskError represents an error that can be returned from filepicker.io service.
type DeskError struct {
	Code    int
	Message string
}

// Error satisfies builtin.error interface. It prints an error string with
// the reason of failure.
func (e DeskError) Error() string {
	return fmt.Sprintf("Desk: %d - %s", e.Code, e.Message)
}
