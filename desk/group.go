package desk

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Group struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func (client *Client) GetGroup(id int64) (*Group, error) {
	url := fmt.Sprintf("%s/groups/%d", DeskURL, id)
	resp, err := client.do("GET", url, "application/json", nil)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	group := Group{}
	err = json.Unmarshal(body, &group)
	if err != nil {
		return nil, err
	}

	return &group, nil
}
