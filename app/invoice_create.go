package app

import (
	"net/http"

	"github.com/andrewarrow/feedback/router"
)

func handleInvoiceCreate(c *router.Context, guid string) {
	//c.ReadFormValuesIntoParams("")

	returnPath := "/sd/invoices"

	item := map[string]any{"text": "test", "amount": 1000}
	items := []map[string]any{item}

	client := c.One("client", "where guid=$1", guid)
	c.Params = map[string]any{}
	c.Params["user_id"] = c.User["id"]
	c.Params["client_id"] = client["id"]
	c.Params["total"] = 1000
	c.Params["items"] = items
	message := c.ValidateCreate("invoice")
	if message != "" {
		router.SetFlash(c, message)
		http.Redirect(c.Writer, c.Request, returnPath, 302)
		return
	}
	message = c.Insert("invoice")
	if message != "" {
		router.SetFlash(c, message)
		http.Redirect(c.Writer, c.Request, returnPath, 302)
		return
	}
	http.Redirect(c.Writer, c.Request, returnPath, 302)
}
