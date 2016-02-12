package desk

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"
)

type Message struct {
	Subject    string    `json:"subject"`
	Body       string    `json:"body"`
	Direction  string    `json:"direction"`
	Status     string    `json:"status"`
	To         string    `json:"to"`
	From       string    `json:"from"`
	CC         string    `json:"cc"`
	BCC        string    `json:"bcc"`
	ClientType string    `json:"client_type"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (client *Client) GetMessage(c *Case) (*Message, error) {
	url := fmt.Sprintf("%s/cases/%d/message", DeskURL, c.Id)
	resp, err := client.do("GET", url, "application/json", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := readError(resp); err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	m := Message{}
	err = json.Unmarshal(body, &m)
	if err != nil {
		return nil, err
	}

	return &m, nil

}
