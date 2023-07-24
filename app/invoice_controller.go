package app

import (
	"github.com/andrewarrow/feedback/router"
	"github.com/andrewarrow/feedback/util"
)

func HandleInvoices(c *router.Context, second, third string) {
	if NotLoggedIn(c) {
		return
	}
	if second == "" && third == "" && c.Method == "GET" {
		handleInvoiceIndex(c)
		return
	}
	if second != "" && third == "" && c.Method == "GET" {
		handleInvoiceShow(c, second)
		return
	}
	if second != "" && third == "" && c.Method == "POST" {
		handleInvoiceShowPost(c, second)
		return
	}
	c.NotFound = true
}

func handleInvoiceIndex(c *router.Context) {
	list := c.All("invoice", "where user_id=$1 order by updated_at desc", "", c.User["id"])
	clientIds := []any{}
	for _, item := range list {
		clientIds = append(clientIds, item["client_id"])
	}
	clientMap := c.WhereIn("client", clientIds)

	colAttributes := map[int]string{}
	//colAttributes[0] = "w-1/2"

	m := map[string]any{}
	headers := []string{"client", "number", "items", "total", "created"}

	params := map[string]any{}
	params["client_map"] = clientMap
	m["headers"] = headers
	m["cells"] = c.MakeCells(util.ToAnyArray(list), headers, params, "_invoice")
	m["col_attributes"] = colAttributes

	topVars := map[string]any{}
	send := map[string]any{}
	send["bottom"] = c.Template("table_show.html", m)
	send["top"] = c.Template("invoices_list_top.html", topVars)
	c.SendContentInLayout("generic_top_bottom.html", send, 200)
}
