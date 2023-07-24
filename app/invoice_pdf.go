package app

import (
	"bytes"

	"github.com/andrewarrow/feedback/router"
	"github.com/jung-kurt/gofpdf"
)

func generatePdf(c *router.Context, invoice map[string]any) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	//fontPath := "path/to/your/nice_font.ttf"
	//pdf.AddFont("CustomFont", "", fontPath)
	pdf.SetFont("helvetica", "", 10)

	userId := invoice["user_id"]
	clientId := invoice["client_id"]
	header := c.One("template", "where flavor=$1 and user_id=$2", "header", userId)
	headerText := header["text"].(string)
	pdf.Text(20, 30, headerText)

	client := c.One("client", "where id=$1", clientId)
	name := client["name"].(string)
	pdf.Text(20, 40, name)

	var buffer bytes.Buffer
	pdf.Output(&buffer)

	contentType := "application/pdf"
	c.Writer.Header().Set("Content-Type", contentType)
	c.Writer.Write(buffer.Bytes())
}
