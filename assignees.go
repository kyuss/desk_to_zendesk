package main

import (
	"./desk"

)

const DEFAULT_ASSIGNEE = 3414708898

var assigneeMap = map[string]int64{
	"assignee@email.com":    3341625747,	
}

const defaultGroup = 27855858

var statusMap = map[string]string{
	"new":      "new",
	"open":     "open",
	"pending":  "pending",
	"resolved": "solved",
	"closed":   "closed",
}

var priorityMap = map[int]string{
	1:  "low",
	2:  "low",
	3:  "low",
	4:  "normal",
	5:  "normal",
	6:  "normal",
	7:  "high",
	8:  "high",
	9:  "urgent",
	10: "urgent",
}

func (migrator *Migrator) migrateAssignee(c *desk.Case) int64 {
	assigneeId := int64(0)
	if c.Assignee != nil && c.Assignee.Id != 0 {
		assigneeId = getAssigneeId(c.Assignee.Email)
	}

	return assigneeId
}

func getAssigneeId(email string) int64 {
	if id, ok := assigneeMap[email]; ok {
		return id
	} else {
		return DEFAULT_ASSIGNEE
	}
}
