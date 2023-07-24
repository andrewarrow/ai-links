package app

import "github.com/andrewarrow/feedback/router"

func generatePdf(c *router.Context, invoice map[string]any) {
	value := "test"
	contentType := "text/plain"
	c.Writer.Header().Set("Content-Type", contentType)
	c.Writer.Write([]byte(value))
}
