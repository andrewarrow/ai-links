package app

import (
	"github.com/andrewarrow/feedback/router"
	"github.com/andrewarrow/feedback/util"
)

func HandleTemplates(c *router.Context, second, third string) {
	if NotLoggedIn(c) {
		return
	}
	if second == "" && third == "" && c.Method == "GET" {
		handleTemplateIndex(c)
		return
	}
	c.NotFound = true
}

func handleTemplateIndex(c *router.Context) {
	list := c.All("template", "where user_id=$1 order by created_at desc", "", c.User["id"])

	colAttributes := map[int]string{}
	//colAttributes[0] = "w-1/2"

	m := map[string]any{}
	headers := []string{"client", "items", "total", "created"}

	params := map[string]any{}
	m["headers"] = headers
	m["cells"] = c.MakeCells(util.ToAnyArray(list), headers, params, "_template")
	m["col_attributes"] = colAttributes

	topVars := map[string]any{}
	send := map[string]any{}
	send["bottom"] = c.Template("table_show.html", m)
	send["top"] = c.Template("templates_list_top.html", topVars)
	c.SendContentInLayout("generic_top_bottom.html", send, 200)
}
