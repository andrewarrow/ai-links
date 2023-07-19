package app

import (
	"github.com/andrewarrow/feedback/router"
)

func HandleClients(c *router.Context, second, third string) {
	if second == "" && third == "" && c.Method == "GET" {
		handleClientIndex(c)
		return
	}
	c.NotFound = true
}

func handleClientIndex(c *router.Context) {
	list := getData()

	colAttributes := map[int]string{}
	colAttributes[0] = "w-1/2"

	m := map[string]any{}
	headers := []string{"id", "name", "address", "created_at"}

	params := map[string]any{}
	m["headers"] = headers
	m["cells"] = c.MakeCells(list, headers, params, "_welcome")
	m["col_attributes"] = colAttributes

	send := map[string]any{}
	send["bottom"] = c.Template("table_show.html", m)
	c.SendContentInLayout("generic_top_bottom.html", send, 200)
}
