package desk

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"time"
)

type ReplyList struct {
	TotalEntries int                `json:"total_entries"`
	Page         int                `json:"page"`
	Values       map[string][]Reply `json:"_embedded"`
}

type Reply struct {
	Id         int64           `json:"id"`
	Subject    string          `json:"subject"`
	Body       string          `json:"body"`
	Direction  string          `json:"direction"`
	Status     string          `json:"status"`
	To         string          `json:"to"`
	From       string          `json:"from"`
	CC         string          `json:"cc"`
	BCC        string          `json:"bcc"`
	ClientType string          `json:"client_type"`
	CreatedAt  time.Time       `json:"created_at,omitempty"`
	UpdatedAt  time.Time       `json:"updated_at,omitempty"`
	Links      map[string]Link `json:"_links"`
	User       *User
	Customer   *Customer
}

func (r *Reply) userId() int64 {
	re := regexp.MustCompile(`/api/v2/users/([0-9]+)$`)
	match := re.FindStringSubmatch(r.Links["sent_by"].Href)
	if len(match) != 0 {
		id, err := strconv.ParseInt(match[1], 10, 64)
		if err != nil {
			return 0
		}
		return id
	}
	return 0
}

func (r *Reply) customerId() int64 {
	re := regexp.MustCompile(`/api/v2/customers/([0-9]+)$`)
	match := re.FindStringSubmatch(r.Links["customer"].Href)
	if len(match) != 0 {
		id, err := strconv.ParseInt(match[1], 10, 64)
		if err != nil {
			return 0
		}
		return id
	}
	return 0
}

func (client *Client) ListReplies(c *Case) ([]Reply, error) {
	url := fmt.Sprintf("%s/cases/%d/replies?per_page=%d", DeskURL, c.Id, 100)
	resp, err := client.do("GET", url, "application/json", nil)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	list := ReplyList{}
	err = json.Unmarshal(body, &list)
	if err != nil {
		return nil, err
	}

	replies := make([]Reply, 0)

	for _, reply := range list.Values["entries"] {
		r := reply
		if reply.userId() > 0 {
			user, err := client.GetUser(r.userId())
			if err != nil {
				return nil, err
			}
			r.User = user
		}

		if reply.customerId() > 0 {
			customer, err := client.GetCustomer(r.customerId())
			if err != nil {
				return nil, err
			}
			r.Customer = customer
		}
		replies = append(replies, r)
	}

	return replies, nil
}
