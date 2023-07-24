package app

import (
	"net/http"
	"os"

	"github.com/andrewarrow/feedback/router"
	"github.com/andrewarrow/feedback/util"
	"golang.org/x/crypto/bcrypt"
)

func HandleUsers(c *router.Context, second, third string) {
	if second == "" && third == "" && c.Method == "POST" {
		handleCreateUser(c)
		return
	}
	c.NotFound = true
}

func HashPassword(password string) string {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(hashedPassword)
}

func handleCreateUser(c *router.Context) {
	c.ReadFormValuesIntoParams("username", "password")
	message := c.ValidateCreate("user")
	returnPath := "/" + c.Router.Prefix + "/sessions/new"
	if message != "" {
		router.SetFlash(c, message)
		http.Redirect(c.Writer, c.Request, returnPath, 302)
		return
	}

	c.Params["password"] = HashPassword(c.Params["password"].(string))
	message = c.Insert("user")
	if message != "" {
		router.SetFlash(c, message)
		http.Redirect(c.Writer, c.Request, returnPath, 302)
		return
	}

	row := c.SelectOne("user", "where username=$1", []any{c.Params["username"]})
	guid := util.PseudoUuid()
	c.Params = map[string]any{"guid": guid, "user_id": row["id"].(int64)}
	c.Insert("cookie_token")
	router.SetUser(c, guid, os.Getenv("COOKIE_DOMAIN"))
	returnPath = "/" + c.Router.Prefix

	http.Redirect(c.Writer, c.Request, returnPath, 302)
	return
}
