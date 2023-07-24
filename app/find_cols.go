package app

import "github.com/andrewarrow/feedback/router"

func GetCols(c *router.Context, modelString string) ([]string, map[string]string) {
	model := c.FindModel(modelString)
	cols := []string{}
	editable := map[string]string{}
	for _, f := range model.Fields {
		if f.Flavor == "editable" {
			editable[f.Name] = "string"
		}
		cols = append(cols, f.Name)
	}
	return cols, editable
}
