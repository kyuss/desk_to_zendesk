package desk

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"time"
)

type Case struct {
	Id           int64             `json:"id"`
	ExternalId   string            `json:"external_id"`
	Blurb        string            `json:"blurb"`
	Subject      string            `json:"subject"`
	Priority     int               `json:"priority"`
	Description  string            `json:"description"`
	Status       string            `json:"status"`
	Type         string            `json:"type"`
	Labels       []string          `json:"labels"`
	LabelIds     []int             `json:"label_ids"`
	Language     string            `json:"language"`
	CustomFields map[string]string `json:"custom_fields"`
	CreatedAt    time.Time         `json:"created_at,omitempty"`
	UpdatedAt    time.Time         `json:"updated_at,omitempty"`
	SupressRules bool              `json:"supress_rules"`
	Links        map[string]Link   `json:"_links"`
	Message      *Message
	Customer     *Customer
	Replies      []Reply
	Notes        []Note
	Attachments  []Attachment
	Assignee     *User
	Group        *Group
}

type CaseList struct {
	TotalEntries int               `json:"total_entries"`
	Page         int               `json:"page"`
	Values       map[string][]Case `json:"_embedded"`
}

func (c *Case) customerId() int64 {
	re := regexp.MustCompile(`/api/v2/customers/([0-9]+)$`)
	match := re.FindStringSubmatch(c.Links["customer"].Href)
	if len(match) != 0 {
		id, err := strconv.ParseInt(match[1], 10, 64)
		if err != nil {
			return 0
		}
		return id
	}
	return 0
}

func (c *Case) assigneeId() int64 {
	re := regexp.MustCompile(`/api/v2/users/([0-9]+)$`)
	match := re.FindStringSubmatch(c.Links["assigned_user"].Href)
	if len(match) != 0 {
		id, err := strconv.ParseInt(match[1], 10, 64)
		if err != nil {
			return 0
		}
		return id
	}
	return 0
}

func (c *Case) groupId() int64 {
	re := regexp.MustCompile(`/api/v2/groups/([0-9]+)$`)
	match := re.FindStringSubmatch(c.Links["assigned_group"].Href)
	if len(match) != 0 {
		id, err := strconv.ParseInt(match[1], 10, 64)
		if err != nil {
			return 0
		}
		return id
	}
	return 0
}

func (client *Client) ListCases(page, perPage string) ([]Case, error) {
	url := fmt.Sprintf("%s/cases/search?labels=%s&page=%s&per_page=%s&sort_field=created_at&sort_direction=asc", DeskURL, DeskLabel, page, perPage)
	resp, err := client.do("GET", url, "application/json", nil)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	list := CaseList{}

	err = json.Unmarshal(body, &list)
	if err != nil {
		return nil, err
	}

	return list.Values["entries"], nil
}

func (client *Client) EmbedCase(c *Case) error {
	return client.embedAll(c)
}

func (client *Client) embedAll(c *Case) error {
	replies, err := client.ListReplies(c)
	if err != nil {
		return err
	}
	c.Replies = replies

	c.Message, err = client.GetMessage(c)
	if err != nil {
		return err
	}
	c.Customer, err = client.GetCustomer(c.customerId())
	if err != nil {
		return err
	}
	c.Notes, err = client.ListNotes(c)
	if err != nil {
		return err
	}
	c.Attachments, err = client.ListAttachments(c)
	if err != nil {
		return err
	}
	c.Assignee, err = client.GetUser(c.assigneeId())
	if err != nil {
		return err
	}
	c.Group, err = client.GetGroup(c.groupId())
	if err != nil {
		return err
	}

	return nil
}
