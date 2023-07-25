package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/andrewarrow/BuisnessPDF/pdfType"
	"github.com/andrewarrow/feedback/router"
	"github.com/rs/zerolog"
)

const SMALL_DATE = "January 2, 2006"

var logger zerolog.Logger

func generatePdf(c *router.Context, invoice map[string]any) {
	handler := pdfType.NewInvoice(&logger)
	jsonString := makeJsonForLib()

	userId := invoice["user_id"]
	clientId := invoice["client_id"]
	header := c.One("template", "where flavor=$1 and user_id=$2", "header", userId)
	headerText := header["text"].(string)

	client := c.One("client", "where id=$1", clientId)
	name := client["name"].(string)
	street1 := client["street1"].(string)
	city := client["city"].(string)
	state := client["state"].(string)
	zip := client["zip"].(string)
	country := client["country"].(string)

	var m map[string]any
	json.Unmarshal([]byte(jsonString), &m)
	senderAddress := m["senderAddress"].(map[string]any)
	senderAddress["companyName"] = headerText
	receiverAddress := m["receiverAddress"].(map[string]any)
	receiverAddress["fullForename"] = name
	receiverAddress["fullSurname"] = ""
	receiverAddress["nameTitle"] = ""
	address := map[string]string{}
	address["road"] = street1
	address["zipCode"] = zip
	address["cityName"] = city
	address["state"] = state
	address["country"] = country
	receiverAddress["address"] = address
	invoiceMeta := m["invoiceMeta"].(map[string]any)

	clientGuid := client["guid"].(string)
	number := invoice["number"].(int64)
	updatedAt := invoice["updated_at"].(int64)
	date := time.Unix(updatedAt, 0).Format(SMALL_DATE)
	tokens := strings.Split(clientGuid, "-")
	clientPrintId := strings.ToUpper(tokens[0])

	invoiceMeta["invoiceNumber"] = fmt.Sprintf("%d", number)
	invoiceMeta["invoiceDate"] = date
	invoiceMeta["customerNumber"] = clientPrintId

	newItems := []map[string]any{}
	items := invoice["items"].([]any)
	for i, item := range items {
		thing := item.(map[string]any)

		line := map[string]any{}
		line["positionNumber"] = fmt.Sprintf("%d", i+1)
		line["description"] = thing["text"]
		line["singlePrice"] = thing["amount"]
		line["quantity"] = 1
		line["taxRate"] = 0
		line["unit"] = ""
		line["currency"] = "$"
		newItems = append(newItems, line)
	}
	invoiceBody := m["InvoiceBody"].(map[string]any)
	invoiceBody["invoicedItems"] = newItems

	asBytes, _ := json.Marshal(m)
	jsonString = string(asBytes)

	handler.SetDataFromJson(jsonString)
	//i := pdfType.NewInvoice(&logger)
	//fmt.Println(i)
	pdf, err := handler.GeneratePDF()
	fmt.Println(err)
	var buffer bytes.Buffer
	pdf.Output(&buffer)

	contentType := "application/pdf"
	c.Writer.Header().Set("Content-Type", contentType)
	c.Writer.Write(buffer.Bytes())
}

func makeJsonForLib() string {
	return `{
  "senderAddress": {
    "fullForename": "",
    "fullSurname": "",
    "companyName": "",
    "nameTitle": "",
    "address": {
      "road": "",
      "houseNumber": "",
      "streetSupplement": "",
      "zipCode": "",
      "cityName": "",
      "state": "",
      "country": "",
      "countryCode": ""
    }
  },
  "receiverAddress": {
    "fullForename": "",
    "fullSurname": "",
    "companyName": "",
    "nameTitle": "",
    "address": {
      "road": "",
      "houseNumber": "",
      "streetSupplement": "",
      "zipCode": "",
      "cityName": "",
      "state": "",
      "country": "",
      "countryCode": ""
    }
  },
  "senderInfo": {
    "phone": "",
    "email": "",
    "web": "",
    "logoSvg": "",
    "iban": "",
    "bic": "",
    "taxNumber": "",
    "bankName": ""
  },
  "invoiceMeta": {
    "invoiceNumber": "",
    "invoiceDate": "",
    "customerNumber": "",
    "projectNumber": ""
  },
  "InvoiceBody": {
    "openingText": "",
    "headlineText": "Invoice",
    "serviceTimeText": "",
    "closingText": "Please pay within 7 days.",
    "ustNotice": "",
    "invoicedItems": []
  }
}`
}
