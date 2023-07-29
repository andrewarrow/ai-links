package app

import "github.com/andrewarrow/feedback/router"

func HandleFiles(c *router.Context, second, third string) {
	if second == "" && third == "" && c.Method == "GET" {
		handleFileGet(c)
		return
	}
	if second == "" && third == "" && c.Method == "POST" {
		handleFilePost(c)
		return
	}
	c.NotFound = true
}

func handleFileGet(c *router.Context) {
	send := map[string]any{}
	c.SendContentInLayout("files.html", send, 200)
}
