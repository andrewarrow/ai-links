package app

import (
	"net/http"

	"github.com/andrewarrow/feedback/router"
)

func HandleWelcome(c *router.Context, second, third string) {
	if second == "" && third == "" && c.Method == "GET" {
		handleWelcomeIndex(c)
		return
	}
	c.NotFound = true
}

func handleWelcomeIndex(c *router.Context) {
	if len(c.User) == 0 {
		c.SendContentInLayout("welcome.html", nil, 200)
		return
	}
	returnPath := "/sd/clients"
	http.Redirect(c.Writer, c.Request, returnPath, 302)
}
