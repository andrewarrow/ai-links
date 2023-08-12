package app

import (
	"strings"

	"github.com/andrewarrow/feedback/router"
	"github.com/andrewarrow/feedback/util"
)

func HandleDash(c *router.Context, second, third string) {
	if second == "" && third == "" && c.Method == "GET" {
		handleDashIndex(c)
		return
	}
	c.NotFound = true
}

func handleDashIndex(c *router.Context) {
	sql := `SELECT client_id,                                                                              DATE_TRUNC('month', created_at) AS month,                                               SUM(total) AS total                                                          FROM                                                                                        invoices 
	    WHERE
         user_id=$1
			GROUP BY                                                                                    client_id,                                                                              DATE_TRUNC('month', created_at)                                                     ORDER BY                                                                                    month desc, client_id;`
	list := c.FreeFormSelect(sql, c.User["id"])
	ids := []any{}
	for _, item := range list {
		id := item["client_id"].(int64)
		ids = append(ids, id)
	}
	clientMap := c.WhereIn("client", ids)
	for _, v := range clientMap {
		tokens := strings.Split(v["guid"].(string), "-")
		v["short"] = strings.ToUpper(tokens[0])
	}

	colAttributes := map[int]string{}
	//colAttributes[0] = "w-1/2"

	m := map[string]any{}
	headers := []string{"date", "client", "amount"}

	params := map[string]any{}
	params["client_map"] = clientMap
	m["headers"] = headers
	m["cells"] = c.MakeCells(util.ToAnyArray(list), headers, params, "_dash")
	m["col_attributes"] = colAttributes

	topVars := map[string]any{}
	send := map[string]any{}
	send["bottom"] = c.Template("table_show.html", m)
	send["top"] = c.Template("dashs_list_top.html", topVars)
	c.SendContentInLayout("generic_top_bottom.html", send, 200)
}
