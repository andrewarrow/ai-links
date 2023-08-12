package app

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/andrewarrow/feedback/router"
)

func handleInvoiceSave(c *router.Context, guid string) {
	invoice := c.One("invoice", "where guid=$1", guid)
	list := invoice["items"].([]any)
	params := []string{""}
	for i, _ := range list {
		p := fmt.Sprintf("text_%d", i)
		params = append(params, p)
		p = fmt.Sprintf("amount_%d", i)
		params = append(params, p)
	}
	c.ReadFormValuesIntoParams(params...)
	items := []map[string]any{}
	total := 0
	for i, _ := range list {
		p := fmt.Sprintf("text_%d", i)
		text := c.Params[p]
		p = fmt.Sprintf("amount_%d", i)
		amount := c.Params[p].(string)
		amountInt, _ := strconv.Atoi(amount)
		if amountInt == 0 {
			continue
		}
		total += amountInt
		item := map[string]any{"text": text, "amount": amountInt}
		items = append(items, item)
	}
	c.Params = map[string]any{}
	c.Params["items"] = items
	c.Params["total"] = total
	c.Update("invoice", "where guid=", guid)
	returnPath := "/sd/invoices"
	http.Redirect(c.Writer, c.Request, returnPath, 302)
}
