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
	message = c.Update("user", "where id=", id)
	if message != "" {
		router.SetFlash(c, message)
		http.Redirect(c.Writer, c.Request, returnPath+"/"+id, 302)
		return
	}
	http.Redirect(c.Writer, c.Request, returnPath, 302)
}
