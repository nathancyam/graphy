package generator

import (
	"fmt"
	"golang.org/x/tools/go/packages"
	"os"
	"regexp"
	"strings"
)

type goType struct {
	Modifiers  string
	ImportPath string
	ImportName string
	Name       string
}

var partsRe = regexp.MustCompile(`^([\[\]\*]*)(.*?)(\.\w*)?$`)

func Generate(_ string, t string) error {
	t1, err := parseType(t)
	if err != nil {
		return err
	}

	if err = tpl.Execute(os.Stdout, struct{
		Type string
	}{
		Type: fmt.Sprintf("%s%s.%s", t1.Modifiers, t1.ImportName, t1.Name),
	}); err != nil {
		return err
	}

	return nil
}

func parseType(v string) (*goType, error) {
	parts := partsRe.FindStringSubmatch(v)
	if len(parts) != 4 {
		return nil, fmt.Errorf("type must be in the form []*github.com/import/path.Name")
	}

	t := &goType{
		Modifiers:  parts[1],
		ImportPath: parts[2],
		Name:       strings.TrimPrefix(parts[3], "."),
	}

	if t.Name == "" {
		t.Name = t.ImportPath
		t.ImportPath = ""
	}

	if t.ImportPath != "" {
		p, err := packages.Load(&packages.Config{Mode: packages.NeedName}, t.ImportPath)
		if err != nil {
			return nil, err
		}
		if len(p) != 1 {
			return nil, fmt.Errorf("not found")
		}

		t.ImportName = p[0].Name
	}

	return t, nil
}