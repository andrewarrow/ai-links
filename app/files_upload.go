package app

import (
	"io"
	"net/http"

	"github.com/andrewarrow/feedback/router"
)

func handleFilePost(c *router.Context) {
	err := c.Request.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		router.SetFlash(c, err.Error())
		http.Redirect(c.Writer, c.Request, "/", 302)
		return
	}
	files := c.Request.MultipartForm.File["file"]

	for _, fileHeader := range files {
		name := fileHeader.Filename
		file, _ := fileHeader.Open()
		asBytes, _ := io.ReadAll(file)
		_ = name
		_ = asBytes
	}
	http.Redirect(c.Writer, c.Request, "/", 302)
}
