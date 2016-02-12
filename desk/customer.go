package desk

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"
)

type Customer struct {
	Id          int64     `json:"id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Company     string    `json:"company"`
	CompanyName string    `json:"company_name"`
	Title       string    `json:"title"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Language    string    `json:"language"`
	Avatar      string    `json:"avatar"`
	Emails      []Contact `json:"emails"`
	Phones      []Contact `json:"phone_numbers"`
	Addresses   []string  `json:"addresses"`
}

type Contact struct {
	Type  string
	Value string
}

func (cus *Customer) FullName() string {
	return fmt.Sprintf("%s %s", cus.FirstName, cus.LastName)
}

func (cus *Customer) Email() string {
	if len(cus.Emails) == 0 {
		return ""
	}
	return cus.Emails[0].Value
}

func (client *Client) GetCustomer(id int64) (*Customer, error) {
	url := fmt.Sprintf("%s/customers/%d", DeskURL, id)
	resp, err := client.do("GET", url, "application/json", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	customer := Customer{}
	err = json.Unmarshal(body, &customer)
	if err != nil {
		return nil, err
	}

	return &customer, nil
}
