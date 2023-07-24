package main

import (
	"embed"
	"math/rand"
	"os"
	"settle-down/app"
	"time"

	"github.com/andrewarrow/feedback/router"
)

//go:embed app/feedback.json
var embeddedFile []byte

//go:embed views/*.html
var embeddedTemplates embed.FS

//go:embed assets/**/*.*
var embeddedAssets embed.FS

var buildTag string

func main() {
	rand.Seed(time.Now().UnixNano())
	if len(os.Args) == 1 {
		//PrintHelp()
		return
	}
	arg := os.Args[1]

	if arg == "reset" {
		//r := router.NewRouter("DATABASE_URL")
		//r.ResetDatabase()
	} else if arg == "run" {
		router.BuildTag = buildTag
		router.EmbeddedTemplates = embeddedTemplates
		router.EmbeddedAssets = embeddedAssets
		r := router.NewRouter("DATABASE_URL", embeddedFile)
		r.Paths["/"] = app.HandleWelcome
		r.Paths["clients"] = app.HandleClients
		r.Paths["invoices"] = app.HandleInvoices
		r.Paths["proposals"] = app.HandleProposals
		r.Paths["templates"] = app.HandleTemplates
		r.Paths["sessions"] = app.HandleSessions
		r.Paths["users"] = app.HandleUsers
		r.Prefix = "sd"
		r.ListenAndServe(":3000")
	} else if arg == "help" {
	}

}
