package desk

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

type Attachment struct {
	Id          int64     `json:"id"`
	FileName    string    `json:"file_name"`
	ContentType string    `json:"content_type"`
	Size        int       `json:"size"`
	URL         string    `json:"url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type AttachmentList struct {
	TotalEntries int                     `json:"total_entries"`
	Page         int                     `json:"page"`
	Values       map[string][]Attachment `json:"_embedded"`
}

func (client *Client) ListAttachments(c *Case) ([]Attachment, error) {
	url := fmt.Sprintf("%s/cases/%d/attachments?per_page=%d", DeskURL, c.Id, 100)
	resp, err := client.do("GET", url, "application/json", nil)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	list := AttachmentList{}
	err = json.Unmarshal(body, &list)
	if err != nil {
		return nil, err
	}

	return list.Values["entries"], nil
}

func (client *Client) DownloadFile(a *Attachment, directory string) error {
	resp, err := client.do("GET", a.URL, "", nil)
	defer resp.Body.Close()

	file, err := os.Create(filepath.Clean(filepath.Join(directory, a.FileName)))
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(file, resp.Body)
	return err
}

func (client *Client) GetAttachmentURL(a *Attachment) (io.Reader, error) {
	url := fmt.Sprintf("%s", a.URL)
	resp, err := client.do("GET", url, "", nil)

	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}
