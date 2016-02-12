package desk

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"time"
)

type Note struct {
	Body      string          `json:"body"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	Links     map[string]Link `json:"_links"`
	User      *User
}

type NoteList struct {
	TotalEntries int               `json:"total_entries"`
	Page         int               `json:"page"`
	Values       map[string][]Note `json:"_embedded"`
}

func (n *Note) userId() int64 {
	re := regexp.MustCompile(`/api/v2/users/([0-9]+)$`)
	match := re.FindStringSubmatch(n.Links["user"].Href)
	if len(match) > 0 {
		id, err := strconv.ParseInt(match[1], 10, 64)
		if err != nil {
			return 0
		}
		return id
	}
	return 0
}

func (client *Client) ListNotes(c *Case) ([]Note, error) {
	url := fmt.Sprintf("%s/cases/%d/notes?per_page=%d", DeskURL, c.Id, 100)
	resp, err := client.do("GET", url, "application/json", nil)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	list := NoteList{}
	err = json.Unmarshal(body, &list)
	if err != nil {
		return nil, err
	}

	notes := make([]Note, 0)

	for _, note := range list.Values["entries"] {
		n := note
		if note.userId() > 0 {
			user, err := client.GetUser(n.userId())
			if err != nil {
				return nil, err
			}
			n.User = user
		}
		notes = append(notes, n)
	}

	return notes, nil
}
