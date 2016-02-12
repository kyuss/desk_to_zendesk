package zendesk

import (
	"encoding/json"
	"fmt"
	_ "io"
	"io/ioutil"
	"net/http"
	"os"
)

type Attachment struct {
	Id          int64  `json:"id,omitifempty"`
	FileName    string `json:"file_name"`
	ContnetURL  string `json:"content_url,omitifempty"`
	ContentType string `json:"content_type"`
	Size        int    `json:"size"`
}

type Upload struct {
	Token      string     `json:"token"`
	Attachment Attachment `json:"attachment"`
}

type UploadResponse struct {
	Upload Upload `json:"upload"`
}

func (client *Client) CreateAttachment(attachment *Attachment, directory string) (string, error) {
	url := fmt.Sprintf("%s/uploads.json?filename=%s&size=%d", ZendeskURL, attachment.FileName, attachment.Size)

	filepath := fmt.Sprintf("%s%s", directory, attachment.FileName)

	file, err := os.Open(filepath)
	defer file.Close()
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, file)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", attachment.ContentType)
	req.Header.Set("Content-Length", fmt.Sprintf("%d", attachment.Size))
	req.Header.Set("User-Agent", UserAgentID)
	req.Header.Add("Accept-Encoding", "identity")
	req.SetBasicAuth(client.username, client.password)

	resp, err := client.http.Do(req)

	if err := readError(resp); err != nil {
		return "", err
	}
	defer resp.Body.Close()

	payload, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	uploadResponse := UploadResponse{}
	err = json.Unmarshal(payload, &uploadResponse)
	if err != nil {
		return "", err
	}

	return uploadResponse.Upload.Token, nil
}
