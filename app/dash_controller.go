package app

import (
	"strings"
	"time"

	"github.com/andrewarrow/feedback/router"
)

func HandleDash(c *router.Context, second, third string) {
	if NotLoggedIn(c) {
		return
	}
	if second == "" && third == "" && c.Method == "GET" {
		handleDashIndex(c)
		return
	}
	c.NotFound = true
}

func handleDashIndex(c *router.Context) {
	sql := `SELECT client_id, sum(hours) as hours,                                                                              DATE_TRUNC('month', created_at) AS created_at,                                               SUM(total) AS total                                                          FROM                                                                                        invoices 
	    WHERE
         user_id=$1
			GROUP BY                                                                                    client_id,                                                                              DATE_TRUNC('month', created_at)                                                     ORDER BY                                                                                    created_at desc, client_id;`
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
	lastMonth := int64(0)
	monthMap := map[int64]int64{}
	newList := []any{}
	for _, item := range list {
		monthTime := item["created_at"].(time.Time)
		total := item["total"].(int64)
		month := monthTime.Unix()
		monthMap[month] += total
		newList = append(newList, item)
		if lastMonth != month && lastMonth > 0 {
			totalItem := map[string]any{"created_at": monthTime,
				"total":  monthMap[month],
				"flavor": "month_total"}
			newList = append(newList, totalItem)
		}
		lastMonth = month
	}
	totalItem := map[string]any{"created_at": time.Unix(lastMonth, 0),
		"total":  monthMap[lastMonth],
		"flavor": "month_total"}
	newList = append(newList, totalItem)

	colAttributes := map[int]string{}
	//colAttributes[0] = "w-1/2"

	m := map[string]any{}
	headers := []string{"date", "client_id / industry", "hours / amount"}

	params := map[string]any{}
	params["client_map"] = clientMap
	m["headers"] = headers
	m["cells"] = c.MakeCells(newList, headers, params, "_dash")
	m["col_attributes"] = colAttributes

	topVars := map[string]any{}
	send := map[string]any{}
	send["bottom"] = c.Template("table_show.html", m)
	send["top"] = c.Template("dashs_list_top.html", topVars)
	c.SendContentInLayout("generic_top_bottom.html", send, 200)
}
