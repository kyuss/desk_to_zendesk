package main

import (
	"fmt"
	log "github.com/cihub/seelog"
	"./desk"
	"./zendesk"
	"strings"
)

type Migrator struct {
	deskClient    *desk.Client
	zenDeskClient *zendesk.Client
	path          string
}

func NewMigrator(deskClient *desk.Client, zenDeskClient *zendesk.Client, path string) *Migrator {
	return &Migrator{
		deskClient:    deskClient,
		zenDeskClient: zenDeskClient,
		path:          path,
	}
}

func (migrator *Migrator) Rune(page, perPage string) int {
	tickets, err := migrator.zenDeskClient.GetTickets("tags:desk")
	if err != nil {
		log.Errorf("%v", err)
		return 1
	}
	for _, ticket := range tickets {
		log.Infof("Deleting ticket: %d", ticket.Id)
		err = migrator.zenDeskClient.DeleteTicket(&ticket)
		if err != nil {
			log.Errorf("Counldn't delete ticket: %v", err)
		}
	}

	return 0
}

func (migrator *Migrator) Run(page, perPage string) int {
	cases, err := migrator.collectCases(page, perPage)
	if err != nil {
		log.Errorf("%v", err)
		return 1
	}
	log.Infof("Collected cases: %d", len(cases))

	for _, c := range cases {
		err := migrator.migrateCase(&c)
		if err != nil {
			log.Errorf("%v", err)
			continue
		}

	}
	return 0
}

func (migrator *Migrator) collectCases(page, perPage string) ([]desk.Case, error) {
	log.Infof("Collecting cases to migrate. (%s, %s)", page, perPage)
	return migrator.deskClient.ListCases(page, perPage)
}

func (migrator *Migrator) migrateCase(c *desk.Case) error {
	log.Infof("Migrating case %d.", c.Id)

	err := migrator.checkCase(c)
	if err != nil {
		return err
	}

	err = migrator.checkCustomer(c)
	if err != nil {
		return err
	}

	log.Infof("Migrating customer.")
	requester, err := migrator.migrateCustomer(c.Customer)
	if err != nil || requester.Email == "" {
		return fmt.Errorf("Can't migrate customer: %v. Aborting.", err)
	}
	log.Infof("Customer %s migrated.", requester.Email)

	assigneeId := migrator.migrateAssignee(c)
	log.Infof("Setting assignee %d.", assigneeId)

	err = migrator.migrateTicket(c, requester, assigneeId)
	if err != nil {
		return err
	}

	log.Infof("Finished case %d.", c.Id)
	return nil
}

func (migrator *Migrator) checkCase(c *desk.Case) error {
	if strings.Contains(strings.ToLower(c.Subject), "missed payment on your account") {
		return fmt.Errorf("Missed payment case %d. Aborting.", c.Id)
	}

	exists, err := migrator.checkIfCaseExists(c)
	if err != nil {
		return err
	}

	if exists {
		return fmt.Errorf("Case already exists in Zendesk.")
	}

	err = migrator.deskClient.EmbedCase(c)
	if err != nil {
		return fmt.Errorf("Can't embed all case dependencies: %v. Aborting.", err)
	}

	return nil
}

func (migrator *Migrator) checkIfCaseExists(c *desk.Case) (bool, error) {
	_, err := migrator.zenDeskClient.GetTicketByExternalId(c.Id)
	if err == nil {
		return true, nil
	} else {
		if err != zendesk.ErrMissingTicket {
			return true, err
		}
	}
	return false, nil
}

func (migrator *Migrator) migrateTicket(c *desk.Case, requester *zendesk.Customer, assigneeId int64) error {

	ticket := &zendesk.Ticket{
		ExternalId:  fmt.Sprintf("%d", c.Id),
		RequesterId: requester.Id,
		AssigneeId:  assigneeId,
		Subject:     c.Subject,
		Description: c.Message.Body,
		Tags:        []string{"desk"},
		Type:        "question",
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
		GroupId:     defaultGroup,
		Status:      statusMap[c.Status],
		Priority:    priorityMap[c.Priority],
	}

	comments := make([]zendesk.Comment, 0)

	log.Infof("Migrating replies: %d.", len(c.Replies))
	replies, err := migrator.migrateReplies(c, requester)
	if err != nil {
		return err
	}
	comments = append(comments, replies...)

	log.Infof("Migrating notes: %d.", len(c.Notes))
	notes, err := migrator.migrateNotes(c)
	if err != nil {
		return err
	}

	comments = append(comments, notes...)

	log.Infof("Migrating attachments: %d.", len(c.Attachments))
	attachments, err := migrator.migrateAttachments(c, requester)
	if err != nil {
		return err
	}

	comments = append(comments, attachments...)

	ticket.Comments = comments

	_, err = migrator.zenDeskClient.CreateTicket(ticket)
	if err != nil {
		return fmt.Errorf("Can't migrate ticket: %v. Aborting.", err)
	}

	return nil
}

func (migrator *Migrator) migrateReplies(c *desk.Case, requester *zendesk.Customer) ([]zendesk.Comment, error) {
	comments := make([]zendesk.Comment, 0)

	for _, reply := range c.Replies {
		authorId := int64(DEFAULT_ASSIGNEE)

		if reply.Direction == "out" {
			if reply.User != nil && reply.User.Email != "" {
				authorId = getAssigneeId(reply.User.Email)
			} else {
				authorId = DEFAULT_ASSIGNEE
			}
		} else if reply.Direction == "in" {
			if reply.Customer != nil && reply.Customer.Email() != "" {
				author, err := migrator.migrateCustomer(reply.Customer)
				if err != nil {
					log.Errorf("%v", err)
					authorId = requester.Id
				} else {
					authorId = author.Id
				}
			}
		}

		comment := zendesk.Comment{
			AuthorId:  authorId,
			Value:     reply.Body,
			CreatedAt: reply.CreatedAt,
			Public:    true,
		}

		if len(strings.TrimSpace(comment.Value)) > 0 {
			comments = append(comments, comment)
		}
	}

	return comments, nil
}

func (migrator *Migrator) migrateNotes(c *desk.Case) ([]zendesk.Comment, error) {
	comments := make([]zendesk.Comment, 0)

	for _, note := range c.Notes {
		authorId := int64(DEFAULT_ASSIGNEE)
		if note.User != nil && note.User.Email != "" {
			authorId = getAssigneeId(note.User.Email)
		}

		comment := zendesk.Comment{
			AuthorId:  authorId,
			Value:     note.Body,
			CreatedAt: note.CreatedAt,
			Public:    false,
		}

		if len(strings.TrimSpace(comment.Value)) > 0 {
			comments = append(comments, comment)
		}
	}

	return comments, nil
}
