package repogen

import (
	"bytes"
	_ "embed"
	"go/types"
	"strings"
	"text/template"
)

type MethodType string

const (
	MethodSlice      MethodType = "SLICE"
	MethodSingular              = "SINGULAR"
	MethodDataloader            = "DATALOADER"
)

type MethodDefinition struct {
	Name       string `yaml:"name"`
	Cypher     string `yaml:"cypher"`
	Output     string `yaml:"output"`
	Dataloader bool   `yaml:"dataloader"`
}

type RepositoryMethod struct {
	CypherParams []string
	Definition   MethodDefinition
	ModelType    string
	Params       string
	Return       string
	ParamTuple   *types.Tuple
	MethodType   MethodType
}

func (r RepositoryMethod) IsSingular() bool {
	return r.MethodType == MethodSingular
}

func (r RepositoryMethod) IsSlice() bool {
	return r.MethodType == MethodSlice
}

func (r RepositoryMethod) IsDataloader() bool {
	return r.MethodType == MethodDataloader
}

func (r RepositoryMethod) GetDataloaderMap() string {
	return "map[string][]" + r.GetBaseModel()
}

func (r RepositoryMethod) GetDataloaderSlice() string {
	return r.ParamTuple.At(1).Name()
}

func (r RepositoryMethod) GetBaseModel() string {
	if strings.Contains(r.ModelType, "[]") {
		idx := strings.LastIndex(r.ModelType, "]")
		return r.ModelType[idx+1:]
	}
	return r.ModelType
}

type RepositoryDefinition struct {
	Package    string             `yaml:"package"`
	Name       string             `yaml:"name"`
	Implements string             `yaml:"implements"`
	Methods    []MethodDefinition `yaml:",flow"`
}

func (r RepositoryDefinition) FindMethodByName(name string) (MethodDefinition, bool) {
	for _, m := range r.Methods {
		if m.Name == name {
			return m, true
		}
		continue
	}
	return MethodDefinition{}, false
}

type TemplateParams struct {
	RepositoryDefinition RepositoryDefinition
	Methods              []RepositoryMethod
	Models               []string
	Imports              []string
}

func (t TemplateParams) RepoName() string {
	return t.RepositoryDefinition.Name
}

func (t TemplateParams) Package() string {
	return t.RepositoryDefinition.Package
}

var (
	//go:embed repository_template.tmpl
	template1 string
	tmpl      = template.Must(template.New("").Parse(template1))
)

func RunTemplate(t TemplateParams) (string, error) {
	b := bytes.NewBuffer(nil)
	if err := tmpl.Execute(b, t); err != nil {
		return "", err
	}

	return b.String(), nil
}
