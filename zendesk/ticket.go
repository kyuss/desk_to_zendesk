package zendesk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"
)

type Ticket struct {
	Id          int64     `json:"id,omitifempty"`
	ExternalId  string    `json:"external_id"`
	RequesterId int64     `json:"requester_id"`
	AssigneeId  int64     `json:"assignee_id,omitifempty"`
	Subject     string    `json:"subject"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	Priority    string    `json:"priority"`
	GroupId     int64     `json:"group_id"`
	Tags        []string  `json:"tags"`
	Type        string    `json:"type"`
	Comments    []Comment `json:"comments"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Comment struct {
	AuthorId  int64     `json:"author_id,omitifempty"`
	Value     string    `json:"value"`
	CreatedAt time.Time `json:"created_at"`
	Public    bool      `json:"public"`
	Uploads   []string  `json:"uploads,omitifempty"`
}

type CreateTicketPayload struct {
	Ticket *Ticket `json:"ticket"`
}

func (client *Client) CreateTicket(ticket *Ticket) (*Ticket, error) {
	request := CreateTicketPayload{Ticket: ticket}

	body, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%s/imports/tickets.json", ZendeskURL)
	resp, err := client.do("POST", url, "application/json", bytes.NewBuffer(body))
	if err := readError(resp); err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	payload, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	response := CreateTicketPayload{}
	err = json.Unmarshal(payload, &response)
	if err != nil {
		return nil, err
	}

	return response.Ticket, nil
}

type TicketList struct {
	Tickets []Ticket `json:"tickets"`
	Count   int64    `json:"count"`
}

func (client *Client) GetTicketByExternalId(id int64) (*Ticket, error) {
	url := fmt.Sprintf("%s/tickets.json?external_id=%d",ZendeskURL, id)
	resp, err := client.do("GET", url, "application/json", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	list := TicketList{}
	err = json.Unmarshal(body, &list)
	if err != nil {
		return nil, err
	}

	if list.Count == 0 || len(list.Tickets) == 0 {
		return nil, ErrMissingTicket
	}
	ticket := list.Tickets[0]
	return &ticket, nil
}

type TicketSearchResult struct {
	Results []Ticket `json:"results"`
	Count   int      `json:"count"`
}

func (client *Client) GetTickets(query string) ([]Ticket, error) {
	url := fmt.Sprintf("%s/search.json?query=%s&per_page=100&page=1", ZendeskURL, query)
	resp, err := client.do("GET", url, "application/json", nil)
	if err != nil {
		return []Ticket{}, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []Ticket{}, err
	}

	list := TicketSearchResult{}
	err = json.Unmarshal(body, &list)
	if err != nil {
		return []Ticket{}, err
	}
	fmt.Println(list.Count)
	return list.Results, nil
}

func (client *Client) DeleteTicket(ticket *Ticket) error {
	_, err := client.do("DELETE", fmt.Sprintf("%s/v2/tickets/%d.json", ZendeskURL, ticket.Id), "application/json", nil)
	if err != nil {
		return err
	}
	return nil
}
