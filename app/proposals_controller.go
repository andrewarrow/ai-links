package app

import (
	"github.com/andrewarrow/feedback/router"
	"github.com/andrewarrow/feedback/util"
)

func HandleProposals(c *router.Context, second, third string) {
	if NotLoggedIn(c) {
		return
	}
	if second == "" && third == "" && c.Method == "GET" {
		handleProposalIndex(c)
		return
	}
	c.NotFound = true
}

func handleProposalIndex(c *router.Context) {
	list := c.All("proposal", "where user_id=$1 order by created_at desc", "", c.User["id"])

	colAttributes := map[int]string{}
	//colAttributes[0] = "w-1/2"

	m := map[string]any{}
	headers := []string{"client", "items", "total", "created"}

	params := map[string]any{}
	m["headers"] = headers
	m["cells"] = c.MakeCells(util.ToAnyArray(list), headers, params, "_proposal")
	m["col_attributes"] = colAttributes

	topVars := map[string]any{}
	send := map[string]any{}
	send["bottom"] = c.Template("table_show.html", m)
	send["top"] = c.Template("proposals_list_top.html", topVars)
	c.SendContentInLayout("generic_top_bottom.html", send, 200)
}
