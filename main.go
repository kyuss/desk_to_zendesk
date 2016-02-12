package main

import (
	"flag"
	log "github.com/cihub/seelog"
	"./desk"
	"./zendesk"
	"os"
)

var (
	deskUser     = ""
	deskPassword = ""

	zenDeskUser     = ""
	zenDeskPassword = ""
)

func main() {

	var page = flag.String("page", "", "page")
	var perPage = flag.String("per_page", "", "Per Page")
	var path = flag.String("path", "", "path")

	flag.Parse()

	migrator := NewMigrator(
		desk.NewClient(deskUser, deskPassword), zendesk.NewClient(zenDeskUser, zenDeskPassword), *path)

	code := migrator.Run(*page, *perPage)
	log.Flush()
	os.Exit(code)
}
