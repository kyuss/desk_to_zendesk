package main

import (
	"./desk"
	"./zendesk"
	"fmt"
	log "github.com/cihub/seelog"
)

func (migrator *Migrator) migrateAttachments(c *desk.Case, requester *zendesk.Customer) ([]zendesk.Comment, error) {
	comments := make([]zendesk.Comment, 0)
	for _, attachment := range c.Attachments {
		err := migrator.deskClient.DownloadFile(&attachment, fmt.Sprintf("%s/attachments/", migrator.path))
		if err != nil {
			log.Errorf("Can't download attachment: %v. Ignoring.", err)
			continue
		}
		token, err := migrator.migrateAttachment(&attachment, fmt.Sprintf("%s/attachments/", migrator.path))
		if err != nil {
			log.Errorf("Can't migrate attachment: %v. Ignoring", err)
			continue
		}

		if token != "" {
			comments = append(comments, zendesk.Comment{
				AuthorId: requester.Id,
				Value:    attachment.FileName,
				Uploads:  []string{token},
			})
		}
	}

	return comments, nil
}

func (migrator *Migrator) migrateAttachment(attachment *desk.Attachment, directory string) (string, error) {

	att := zendesk.Attachment{
		FileName:    attachment.FileName,
		ContentType: attachment.ContentType,
		Size:        attachment.Size,
	}

	token, err := migrator.zenDeskClient.CreateAttachment(&att, directory)

	if err != nil {
		return "", err
	}

	return token, nil
}
