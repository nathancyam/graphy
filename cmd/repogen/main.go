package main

import (
	"fmt"
	"go/types"
	"golang.org/x/tools/go/packages"
	"gopkg.in/yaml.v2"
	"graphy/cmd/repogen/repogen"
	"log"
	"os"
	"path"
	"strings"
)

// Cypher repository generator tool. This tool should be run with the following command:
//
// > go run graphy/cmd/repogen/main.go path/to/repo.yaml
//
// This tool creates a basic repository for your methods which are based from an given interface
// as provided in the YAML file.
func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: return type")
		fmt.Println("repogen path/to/repository.yaml")
		os.Exit(1)
	}

	d, under := getRepositoryInterfaceData()
	rr := make([]repogen.RepositoryMethod, under.NumMethods())
	var importLookup = make(map[string]string)
	var modelLookup = make(map[string]string)

	totalMethods := under.NumMethods()
	if totalMethods != len(d.Methods) {
		log.Fatal(fmt.Errorf("configuration yaml file must implement all methods provided on the interface, Expected: %d, Got: %d", totalMethods, len(d.Methods)))
	}

	for i := 0; i < under.NumMethods(); i++ {
		fn := under.Method(i)
		signature := fn.Type().(*types.Signature)
		methodDef, ok := d.FindMethodByName(fn.Name())
		if !ok {
			log.Fatal(fmt.Errorf("a method definition is missing from the yaml file %s %v", fn.Name(), methodDef))
		}

		params := signature.Params()
		cypherArgs, ok := getCypherParams(params, methodDef.Cypher)
		if !ok {
			log.Fatal(fmt.Errorf("cypher query provided does not align with repository parameters"))
		}

		modelType, basePkg, queryType := getModelType(signature.Results())
		if !methodDef.Dataloader {
			modelLookup[modelType] = modelType
		}
		importLookup[basePkg] = basePkg

		var fullImports []string
		fullImports = append(fullImports, basePkg)

		returnType := signature.Results().String()
		for _, imp := range fullImports {
			splits := path.Dir(imp)
			returnType = strings.ReplaceAll(returnType, splits, "")
			returnType = strings.ReplaceAll(returnType, "/", "")
		}

		rr[i] = repogen.RepositoryMethod{
			Definition:   methodDef,
			Params:       params.String(),
			ParamTuple:   params,
			Return:       returnType,
			MethodType:   queryType,
			CypherParams: cypherArgs,
			ModelType:    modelType,
		}
	}

	var models []string
	for _, m := range modelLookup {
		models = append(models, m)
	}

	var imports []string
	for _, i := range importLookup {
		imports = append(imports, i)
	}

	s, err := repogen.RunTemplate(repogen.TemplateParams{
		RepositoryDefinition: d,
		Methods:              rr,
		Models:               models,
		Imports:              imports,
	})
	if err != nil {
		log.Fatal(fmt.Errorf("failed to generate code from template"))
	}

	fmt.Println(s)
}

func getCypherParams(params *types.Tuple, cypherQuery string) ([]string, bool) {
	var validParams []string
	// Skip the first one, assuming that it is a context
	for i := 1; i < params.Len(); i++ {
		paramName := params.At(i).Name()
		if !strings.Contains(cypherQuery, "$"+paramName) {
			return nil, false
		}
		validParams = append(validParams, paramName)
	}

	return validParams, true
}

func getRepositoryInterfaceData() (repogen.RepositoryDefinition, *types.Interface) {
	yamlPath := os.Args[1]
	b, err := os.ReadFile(yamlPath)
	if err != nil {
		log.Fatal(fmt.Errorf("could not read YAML file from path: %s, %v", yamlPath, err))
	}

	var d repogen.RepositoryDefinition
	if err := yaml.Unmarshal(b, &d); err != nil {
		log.Fatal(fmt.Errorf("could not parse YAML file, most likely malformed: %v", err))
	}

	srcTypePkg, srcType := splitSourceType(d.Implements)

	pkgs, err := packages.Load(&packages.Config{
		Mode: packages.NeedTypes | packages.NeedImports,
	}, srcTypePkg)

	if err != nil {
		log.Fatal(fmt.Errorf("could not load package metadata for interface source package: %v", err))
	}

	targetPkg := pkgs[0]
	obj := targetPkg.Types.Scope().Lookup(srcType)
	if obj == nil {
		log.Fatal(fmt.Errorf("repository interface could not be found: %s", srcType))
	}

	oo, ok := obj.(*types.TypeName)
	if !ok {
		log.Fatal(fmt.Errorf("incorrect type for interface src type provided"))
	}

	under, ok := oo.Type().Underlying().(*types.Interface)
	if !ok {
		log.Fatal(fmt.Errorf("the implements key in YAML configuration must be an interface type"))
	}

	return d, under
}

func getModelType(tuple *types.Tuple) (string, string, repogen.MethodType) {
	var o string
	var m repogen.MethodType

	ret := tuple.At(0)
	switch t := ret.Type().(type) {
	case *types.Slice:
		m = repogen.MethodSlice
		o = t.Elem().String()
	case *types.Map:
		m = repogen.MethodDataloader
		o = t.Elem().String()
	case *types.Pointer:
		m = repogen.MethodSingular
		o = t.Elem().String()
	default:
		return "", "", repogen.MethodSingular
	}

	idx := strings.LastIndex(o, ".")
	pkg := strings.ReplaceAll(o[0:idx], "[]", "")

	if strings.Contains(o, "[]") {
		return `[]` + path.Base(o), pkg, m
	}

	return path.Base(o), pkg, m
}

func splitSourceType(sourceType string) (string, string) {
	idx := strings.LastIndexByte(sourceType, '.')
	if idx == -1 {
		log.Fatal("")
	}
	sourceTypePackage := sourceType[0:idx]
	sourceTypeName := sourceType[idx+1:]
	return sourceTypePackage, sourceTypeName
}
