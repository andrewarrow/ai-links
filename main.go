package main

import (
	"embed"
	"fmt"
	"io/ioutil"
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

	if arg == "template" {
		name := os.Args[2]
		asBytes, _ := ioutil.ReadFile("t")
		top := `{{ define "_%s" }}`
		end := `{{ end }}`
		s := fmt.Sprintf(top, name) + "\n" + string(asBytes) + "\n" + end
		ioutil.WriteFile("views/_"+name+".html", []byte(s), 0644)
		result := fmt.Sprintf(`{{ template "_%s" . }}`, name)
		fmt.Printf("\n\n%s\n\n", result)
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
		r.Paths["files"] = app.HandleFiles
		r.Paths["dash"] = app.HandleDash
		r.Prefix = "sd"
		r.BucketPath = "/Users/aa/bucket"
		r.ListenAndServe(":3003")
	} else if arg == "help" {
	}

}
