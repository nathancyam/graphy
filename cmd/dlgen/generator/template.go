package generator

import "html/template"

var tpl = template.Must(template.New("generated").Parse(`
m := o.(map[string]interface{})
id := m["id"].(string)
items := m["items"].([]interface{})

var rs = make([]{{.Type}}, len(items))
for index, i := range items {
	var r {{.Type}}
	n, ok := i.(neo4j.Node)
	if !ok {
		return nil, err
	}

	if err := mapstructure.Decode(n.Props(), &r); err != nil {
		return nil, err
	}
	rs[index] = r
}

res[id] = rs
`))