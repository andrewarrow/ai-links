package app

import (
	"net/http"

	"github.com/andrewarrow/feedback/router"
	"github.com/andrewarrow/feedback/util"
)

func handleClientShowPost(c *router.Context, guid string) {
	cols, editable := GetCols(c, "client")
	list := []string{}
	for _, item := range cols {
		if router.IsEditable(item, editable) == false {
			continue
		}
		list = append(list, item)
	}
	list = append(list, "submit")
	c.ReadFormValuesIntoParams(list...)
	submit := c.Params["submit"].(string)
	if submit != "save" {
		handleInvoiceCreate(c, guid)
		return
	}

	c.ValidateUpdate("client")
	message := c.ValidateUpdate("client")
	returnPath := "/sd/clients"
	if message != "" {
		router.SetFlash(c, message)
		http.Redirect(c.Writer, c.Request, returnPath+"/"+guid, 302)
		return
	}
	message = c.Update("client", "where guid=", guid)
	if message != "" {
		router.SetFlash(c, message)
		http.Redirect(c.Writer, c.Request, returnPath+"/"+guid, 302)
		return
	}
	http.Redirect(c.Writer, c.Request, returnPath, 302)
}

func handleClientShow(c *router.Context, guid string) {
	client := c.One("client", "where guid=$1", guid)
	regexMap := map[string]string{}
	cols, editable := GetCols(c, "client")
	cols = append(cols, "save")
	editable["save"] = "save"

	colAttributes := map[int]string{}
	colAttributes[1] = "w-3/4"

	m := map[string]any{}
	headers := []string{"field", "value"}

	params := map[string]any{}
	params["item"] = client
	params["editable"] = editable
	params["regex_map"] = regexMap
	m["headers"] = headers
	m["cells"] = c.MakeCells(util.ToAny(cols), headers, params, "_client_show")
	m["col_attributes"] = colAttributes
	topVars := map[string]any{}
	topVars["name"] = client["name"]
	topVars["guid"] = guid
	send := map[string]any{}
	send["bottom"] = c.Template("table_show.html", m)
	send["top"] = c.Template("clients_top.html", topVars)

	c.SendContentInLayout("generic_top_bottom.html", send, 200)
}
