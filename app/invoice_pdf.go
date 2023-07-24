package app

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/andrewarrow/feedback/router"
	"github.com/jung-kurt/gofpdf"
)

const SMALL_DATE = "January 2, 2006"

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
	street1 := client["street1"].(string)
	city := client["city"].(string)
	state := client["state"].(string)
	zip := client["zip"].(string)
	country := client["country"].(string)
	pdf.Text(20, 40, name)
	pdf.Text(20, 45, street1)
	pdf.Text(20, 50, city+", "+state+" "+zip)
	pdf.Text(20, 55, country)

	clientGuid := client["guid"].(string)
	number := invoice["number"].(int64)
	updatedAt := invoice["updated_at"].(int64)
	date := time.Unix(updatedAt, 0).Format(SMALL_DATE)
	tokens := strings.Split(clientGuid, "-")
	clientPrintId := strings.ToUpper(tokens[0])
	pdf.Text(20, 65, "Client ID: "+clientPrintId)
	pdf.Text(20, 70, fmt.Sprintf("Invoice ID: %d", number))
	pdf.Text(20, 75, fmt.Sprintf("Invoice Date: %s", date))

	pdf.SetFont("helvetica", "B", 16)
	pdf.Text(20, 90, "Invoice")

	cellWidth := 180
	cellHeight := 180
	text := "Amount"
	pdf.SetFont("helvetica", "B", 10)
	pdf.CellFormat(float64(cellWidth),
		float64(cellHeight), text, "", 1, "R", false, 0, "")

	pdf.Text(20, 100, "# Description")

	lineWidth := 170.0
	lineHeight := 1.0
	grayColor := 200
	pdf.SetDrawColor(grayColor, grayColor, grayColor)
	pdf.Line(20, 103, 20+lineWidth, 103+lineHeight)

	pdf.SetFont("helvetica", "", 10)
	items := invoice["items"].([]any)
	height := 110.0
	for i, item := range items {
		thing := item.(map[string]any)
		text := fmt.Sprintf("%d. %s", i+1, thing["text"])
		pdf.Text(20, 110+(float64(i)*5), text)
		text = Money(thing["amount"].(float64))
		pdf.Text(170, 110+(float64(i)*5), text)
		height += float64(i) * 5
	}

	pdf.Line(20, height, 20+lineWidth, height+lineHeight)
	total := invoice["total"].(int64)
	text = "Total: " + Money(float64(total))
	pdf.CellFormat(180, (-1*height)+10, text, "", 1, "R", false, 0, "")

	var buffer bytes.Buffer
	pdf.Output(&buffer)

	contentType := "application/pdf"
	c.Writer.Header().Set("Content-Type", contentType)
	c.Writer.Write(buffer.Bytes())
}

func Money(a float64) string {
	amount := fmt.Sprintf("%d", int64(a))
	if len(amount) < 3 {
		s := fmt.Sprintf("$00.%s USD", amount)
		return s
	}
	dollars := amount[0 : len(amount)-2]
	cents := amount[len(amount)-2:]
	s := fmt.Sprintf("$%s.%s USD", dollars, cents)
	return s
}
