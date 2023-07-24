package app

import (
	"net/http"

	"github.com/andrewarrow/feedback/router"
)

func handleClientShowPost(c *router.Context, id string) {
	list := []string{}
	list = append(list, "submit")
	c.ReadFormValuesIntoParams(list...)
	submit := c.Params["submit"].(string)
	if submit != "save" {
		//handleClientAddNew(c, id)
		return
	}

	c.ValidateUpdate("client")
	message := c.ValidateUpdate("client")
	returnPath := "/sd/clients"
	if message != "" {
		router.SetFlash(c, message)
		http.Redirect(c.Writer, c.Request, returnPath+"/"+id, 302)
		return
	}
	message = c.Update("client", "where id=", id)
	if message != "" {
		router.SetFlash(c, message)
		http.Redirect(c.Writer, c.Request, returnPath+"/"+id, 302)
		return
	}
	http.Redirect(c.Writer, c.Request, returnPath, 302)
}

func handleClientShow(c *router.Context, guid string) {
	client := c.One("client", "where guid=$1", guid)
	model := c.FindModel("client")
	regexMap := map[string]string{}
	cols := []any{}
	for _, f := range model.Fields {
		regexMap[f.Name] = f.Regex
		cols = append(cols, f.Name)
	}

	colAttributes := map[int]string{}
	colAttributes[1] = "w-3/4"

	m := map[string]any{}
	headers := []string{"field", "value"}

	params := map[string]any{}
	params["item"] = client
	params["editable"] = map[string]string{}
	params["regex_map"] = regexMap
	m["headers"] = headers
	m["cells"] = c.MakeCells(cols, headers, params, "_client_show")
	m["col_attributes"] = colAttributes
	topVars := map[string]any{}
	topVars["name"] = client["name"]
	topVars["guid"] = guid
	send := map[string]any{}
	send["bottom"] = c.Template("table_show.html", m)
	send["top"] = c.Template("clients_top.html", topVars)

	c.SendContentInLayout("generic_top_bottom.html", send, 200)
}
