package desk

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"
)

type User struct {
	Id            int64     `json:"id"`
	Name          string    `json:"name"`
	PublicName    string    `json:"public_name"`
	Email         string    `json:"email"`
	EmailVerified bool      `json:email_verified"`
	Avatar        string    `json:"avatar"`
	Level         string    `json:"level"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (client *Client) GetUser(id int64) (*User, error) {
	url := fmt.Sprintf("%s/users/%d", DeskURL, id)
	resp, err := client.do("GET", url, "application/json", nil)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	user := User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
