package app

import (
	"github.com/andrewarrow/feedback/router"
	"github.com/andrewarrow/feedback/util"
)

func HandleClients(c *router.Context, second, third string) {
	if NotLoggedIn(c) {
		return
	}
	if second == "" && third == "" && c.Method == "GET" {
		handleClientIndex(c)
		return
	}
	if second == "" && third == "" && c.Method == "POST" {
		handleClientCreate(c)
		return
	}
	if second != "" && third == "" && c.Method == "GET" {
		handleClientShow(c, second)
		return
	}
	if second != "" && third == "" && c.Method == "POST" {
		handleClientShowPost(c, second)
		return
	}
	c.NotFound = true
}

func handleClientIndex(c *router.Context) {
	list := c.All("client", "order by created_at desc", "")

	colAttributes := map[int]string{}
	//colAttributes[0] = "w-1/2"

	m := map[string]any{}
	headers := []string{"name", "street", "city/state/zip", "created"}

	params := map[string]any{}
	m["headers"] = headers
	m["cells"] = c.MakeCells(util.ToAnyArray(list), headers, params, "_client")
	m["col_attributes"] = colAttributes

	topVars := map[string]any{}
	send := map[string]any{}
	send["bottom"] = c.Template("table_show.html", m)
	send["top"] = c.Template("clients_list_top.html", topVars)
	c.SendContentInLayout("generic_top_bottom.html", send, 200)
}
