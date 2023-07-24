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
	pdf.SetFont("helvetica", "", 14)

	pdf.Text(20, 30, "Hello, this is some nice text!")

	var buffer bytes.Buffer
	pdf.Output(&buffer)

	contentType := "application/pdf"
	c.Writer.Header().Set("Content-Type", contentType)
	c.Writer.Write(buffer.Bytes())
}
