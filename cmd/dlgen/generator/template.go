package generator

import (
	"fmt"
	"go/types"
	"path"
	"strings"
	"text/template"
)

type ServiceDefinition struct {
	Object    types.Object
	Interface *types.Interface
	Func      *types.Func
}

type DataloaderDefinition struct {
	Object types.Object
	Struct *types.Struct
}

type TemplateParams struct {
	FetchSignature *types.Signature
	Dataloader     DataloaderDefinition
	Service        ServiceDefinition
}

func (t TemplateParams) ProviderType() string {
	dl := t.Dataloader.Object.Name()
	return fmt.Sprintf("type %sProvider func(con context.Context) %s", dl, dl)
}

func (t TemplateParams) Provider() string {
	return t.Dataloader.Object.Name() + "Provider"
}

func (t TemplateParams) ServiceName() string {
	return fmt.Sprintf("%s.%s", t.Service.Object.Pkg().Name(), t.Service.Object.Name())
}

func (t TemplateParams) ServiceMethod() string {
	return t.Service.Func.Name()
}

func (t TemplateParams) DataloaderName() string {
	return t.Dataloader.Object.Name()
}

func (t TemplateParams) Fetch() string {
	res := t.FetchSignature.Results()
	var replaceTargets []string

	for i := 0; i < res.Len(); i++ {
		curr := res.At(i).Type()
		for {
			switch t := curr.(type) {
			case *types.Slice:
				curr = t.Elem()
				continue
			case *types.Pointer:
				curr = t.Elem()
				continue
			default:
				break
			}
			break
		}

		str, ok := curr.(*types.Named)
		if !ok || str.Obj().Pkg() == nil {
			continue
		}

		pkgPath := str.Obj().Pkg().Path()
		pkgName := str.Obj().Pkg().Name()
		s := strings.ReplaceAll(pkgPath, pkgName, "")
		replaceTargets = append(replaceTargets, s)
	}

	target := t.FetchSignature.String()
	for _, replace := range replaceTargets {
		target = strings.ReplaceAll(target, replace, "")
	}

	return target
}

func (t TemplateParams) FetchKeys() string {
	res := t.FetchSignature.Params().At(0).Name()
	return res
}

func (t TemplateParams) OutputType() string {
	orig := t.FetchSignature.Results().At(0).Type()
	pkg, _ := t.getFetchNamedType()
	if pkg == nil {
		return ""
	}

	pkgPath := path.Dir(pkg.Path())
	pkgPath = strings.ReplaceAll(orig.String(), pkgPath, "")
	return strings.ReplaceAll(pkgPath, "/", "")
}

func (t TemplateParams) InputType() string {
	sign, ok := t.Service.Func.Type().(*types.Signature)
	if !ok {
		return "// Not able to find your mapping"
	}

	// Assume that the first returned value is the target, usually a map of some kind
	ret := sign.Results().At(0)
	pkg, n, ok := t.fetchBaseType(ret.Type())
	if !ok {
		return "// Not able to find your mapping"
	}
	pkgPath := path.Dir(pkg.Path())
	pkgPath = strings.ReplaceAll(n.String(), pkgPath, "")
	return strings.ReplaceAll(pkgPath, "/", "")
}

func (t TemplateParams) Imports() []string {
	var imports []string

	// Import the package for the GraphQL model
	fetchResults := t.FetchSignature.Results()
	for i := 0; i < t.FetchSignature.Results().Len(); i++ {
		res := fetchResults.At(i).Type()
		pkg, _, ok := t.fetchBaseType(res)
		if !ok {
			continue
		}
		imports = append(imports, pkg.Path())
	}

	imports = append(imports, t.Service.Object.Pkg().Path())
	return imports
}

func (t TemplateParams) UnwrapSlice(t1 string, times int) string {
	if times == 0 {
		return t1
	}

	if t1[0:2] == "[]" {
		return t.UnwrapSlice(t1[2:], times - 1)
	}

	return t1
}

func (t TemplateParams) BaseStruct(t1 string) string {
	if t1[0:2] == "[]" {
		return t.BaseStruct(t1[2:])
	}
	if t1[0:1] == "*" {
		return t1[1:]
	}
	return t1
}

func (t TemplateParams) Mapping() string {
	_, fetchType := t.getFetchNamedType()
	if fetchType == nil {
		return "// Not able to find your mapping"
	}

	sign, ok := t.Service.Func.Type().(*types.Signature)
	if !ok {
		return "// Not able to find your mapping"
	}

	// Assume that the first returned value is the target, usually a map of some kind
	ret := sign.Results().At(0)
	_, n, ok := t.fetchBaseType(ret.Type())
	if !ok {
		return "// Not able to find your mapping"
	}

	coreStruct, ok := n.Underlying().(*types.Struct)
	if !ok {
		return "// Not able to find your mapping"
	}
	encodedStruct, ok := fetchType.Underlying().(*types.Struct)
	if !ok {
		return "// Not able to find your mapping"
	}

	var encodedFields = make(map[string]string, encodedStruct.NumFields())
	var coreFields = make(map[string]string, coreStruct.NumFields())

	for i := 0; i < encodedStruct.NumFields(); i++ {
		field := encodedStruct.Field(i)
		encodedFields[field.Name()] = field.Type().String()
	}

	for i := 0; i < coreStruct.NumFields(); i++ {
		field := coreStruct.Field(i)
		coreFields[field.Name()] = field.Type().String()
	}

	structCallable := true
	for k, v := range encodedFields {
		if !structCallable {
			break
		}

		coreVal, ok := coreFields[k]
		if !ok || v != coreVal {
			structCallable = false
			break
		}
	}

	if structCallable {
		return fmt.Sprintf("return %s(%s), nil", t.BaseStruct(t.OutputType()), "m")
	}

	return fmt.Sprintf(`// No direct mapping with model structs found. You will have to manually enter the decode code yourself
	return %s{}, nil`, t.BaseStruct(t.OutputType()))
}

func (t TemplateParams) fetchBaseType(curr types.Type) (*types.Package, *types.Named, bool) {
	for {
		switch t := curr.(type) {
		case *types.Slice:
			curr = t.Elem()
			continue
		case *types.Pointer:
			curr = t.Elem()
			continue
		case *types.Map:
			curr = t.Elem()
			continue
		default:
			break
		}
		break
	}

	named, ok := curr.(*types.Named)
	if !ok || named.Obj().Pkg() == nil {
		return nil, nil, false
	}

	return named.Obj().Pkg(), named, true
}

func (t TemplateParams) getFetchNamedType() (*types.Package, *types.Named) {
	curr := t.FetchSignature.Results().At(0).Type()

	pkg, named, ok := t.fetchBaseType(curr)
	if !ok {
		return nil, nil
	}
	return pkg, named
}

var tpl = template.Must(template.New("generated").Parse(`
// Code generated by dlgen, DO NOT EDIT.

package dataloader

import (
	"context"
	"errors"
	"time"
	{{range $i, $import := .Imports -}}
	"{{$import}}"
    {{end -}}
)

{{.ProviderType}}

func NewGradeDLoader(svc {{.ServiceName}}) {{.Provider}} {
	return func(ctx context.Context) {{.DataloaderName}} {
		return {{.DataloaderName}} {
			fetch: {{.Fetch}} {
				// TODO: Generate the service here
				m, err := svc.{{.ServiceMethod}}(ctx, {{.FetchKeys}})

				if err != nil {
					// TODO: Make error message better
					return nil, []error{errors.New("")}
				}

				var out = make({{.OutputType}}, len({{.FetchKeys}}))
				for index, keyID := range {{.FetchKeys}} {
					resolvedVal := m[keyID]
					var decodedItems = make({{.UnwrapSlice .OutputType 1}}, len(resolvedVal))
					for i, rawItem := range resolvedVal {
						decoded, _ := decode(rawItem)
						decodedItems[i] = &decoded
					}

					out[index] = decodedItems
				}

				return out, nil
			},
			wait:     1 * time.Millisecond,
			maxBatch: 100,
		}
	}
}

func decode(m {{.InputType}}) ({{.BaseStruct .OutputType}}, error) {
    {{.Mapping}}
}`))
