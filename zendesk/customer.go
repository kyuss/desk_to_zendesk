package zendesk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"
)

type Customer struct {
	Id         int64     `json:"id"`
	ExternalId string    `json:"external_id"`
	Name       string    `json:"name"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Email      string    `json:"email"`
	Phone      string    `json:"phone"`
	Role       string    `json:"role"`
	Verified   bool      `json:"verified"`
}

type CreateCustomerRequest struct {
	User *Customer `json:"user"`
}

func (client *Client) CreateCustomer(customer *Customer) (*Customer, error) {
	request := CreateCustomerRequest{User: customer}

	body, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%s/users.json",ZendeskURL)
	resp, err := client.do("POST", url, "application/json", bytes.NewBuffer(body))

	if err := readError(resp); err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	payload, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	cus := CreateCustomerRequest{}

	err = json.Unmarshal(payload, &cus)

	return cus.User, nil
}

type SearchResult struct {
	Customers []Customer `json:"users"`
	Count     int        `json:"count"`
}

func (client *Client) GetCustomerByEmail(email string) (*Customer, error) {
	url := fmt.Sprintf("https://filestack.zendesk.com/api/v2/users/search.json?query=%s", email)
	resp, err := client.do("GET", url, "application/json", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	list := SearchResult{}
	err = json.Unmarshal(body, &list)
	if err != nil {
		return nil, err
	}

	if list.Count == 0 {
		return nil, ErrMissingEmail
	}
	cus := list.Customers[0]
	if cus.Id == 0 {
		fmt.Println(string(body))
	}

	return &cus, nil
}
