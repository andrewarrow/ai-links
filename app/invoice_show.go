package app

import (
	"net/http"

	"github.com/andrewarrow/feedback/router"
)

func handleInvoiceShowPost(c *router.Context, guid string) {
	c.ReadFormValuesIntoParams("submit")
	invoice := c.One("invoice", "where guid=$1", guid)
	list := invoice["items"].([]any)
	total := invoice["total"].(int64)
	item := map[string]any{"text": "test", "amount": 1000}
	list = append(list, item)
	c.Params["items"] = list
	c.Params["total"] = total + 1000
	c.ValidateUpdate("invoice")
	message := c.ValidateUpdate("invoice")
	returnPath := "/sd/invoices/" + guid
	if message != "" {
		router.SetFlash(c, message)
		http.Redirect(c.Writer, c.Request, returnPath, 302)
		return
	}
	message = c.Update("invoice", "where guid=", guid)
	if message != "" {
		router.SetFlash(c, message)
		http.Redirect(c.Writer, c.Request, returnPath, 302)
		return
	}
	http.Redirect(c.Writer, c.Request, returnPath, 302)
}

func handleInvoiceShow(c *router.Context, guid string) {
	invoice := c.One("invoice", "where guid=$1", guid)
	regexMap := map[string]string{}
	cols, editable := GetCols(c, "invoice")
	cols = append(cols, "save")
	editable["save"] = "save"

	colAttributes := map[int]string{}
	colAttributes[1] = "w-3/4"

	m := map[string]any{}
	headers := []string{"text", "amount"}

	list := invoice["items"].([]any)

	params := map[string]any{}
	params["item"] = invoice
	params["editable"] = editable
	params["regex_map"] = regexMap
	m["headers"] = headers
	m["cells"] = c.MakeCells(list, headers, params, "_invoice_show")
	m["col_attributes"] = colAttributes
	topVars := map[string]any{}
	topVars["name"] = invoice["name"]
	topVars["guid"] = guid
	send := map[string]any{}
	send["bottom"] = c.Template("table_show.html", m)
	send["top"] = c.Template("invoices_top.html", topVars)

	c.SendContentInLayout("generic_top_bottom.html", send, 200)
}
