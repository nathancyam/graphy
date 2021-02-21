package generator

import (
	"bytes"
	"errors"
	"fmt"
	"go/types"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/imports"
	"io/ioutil"
	"log"
	"strings"
)

func Generate(dlPkgStr string, svcPkgStr string, destination string) error {
	dlObj, err := getObject(dlPkgStr)
	if err != nil {
		return err
	}

	svcPackageStr, svcMethodStr, err := parseSvcString(svcPkgStr)
	if err != nil {
		return err
	}

	svcObj, err := getObject(svcPackageStr)
	if err != nil {
		return err
	}

	dlStruct, ok := dlObj.Type().Underlying().(*types.Struct)
	if !ok {
		return fmt.Errorf("dataloader package referred to must be a struct")
	}

	svcInterface, ok := svcObj.Type().Underlying().(*types.Interface)
	if !ok {
		return errors.New("dataloader package referred must be an interface")
	}

	var svcMethod *types.Func
	for i := 0; i < svcInterface.NumMethods(); i++ {
		method := svcInterface.Method(i)
		if method.Name() == svcMethodStr {
			svcMethod = method
			break
		}
	}

	if svcMethod == nil {
		return fmt.Errorf("the method %s could not be found in the service interface: %s", svcMethodStr, svcPkgStr)
	}

	var fetchField *types.Var
	for i := 0; i < dlStruct.NumFields(); i++ {
		field := dlStruct.Field(i)
		if field.Name() == "fetch" {
			fetchField = field
			break
		}
	}

	if fetchField == nil {
		return errors.New(`dataloader must a "fetch" field defined`)
	}

	fetchSig := fetchField.Type().(*types.Signature)

	//fetchField.
	params := TemplateParams{
		FetchSignature:    fetchSig,
		Dataloader: DataloaderDefinition{
			Object: dlObj,
			Struct: dlStruct,
		},
		Service: ServiceDefinition{
			Object: svcObj,
			Interface: svcInterface,
			Func: svcMethod,
		},
	}

	buf := bytes.NewBuffer(nil)
	if err = tpl.Execute(buf, params); err != nil {
		return err
	}
	src, err := imports.Process(destination, buf.Bytes(), nil)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(destination, src, 0644); err != nil {
		return err
	}

	return nil
}

func parseSvcString(svcString string) (string, string, error) {
	if !strings.Contains(svcString, ":") {
		return "", "", errors.New("service method must be provided with a colon")
	}

	idx := strings.LastIndex(svcString, ":")
	return svcString[0:idx], svcString[idx+1:], nil
}

func getObject(packageStr string) (types.Object, error) {
	srcPkg, srcType := splitSourceType(packageStr)

	loadedPkgs, err := packages.Load(&packages.Config{
		Mode: packages.NeedTypes | packages.NeedImports,
	}, srcPkg)
	if err != nil {
		return nil, err
	}

	targetPkg := loadedPkgs[0]
	dlObj := targetPkg.Types.Scope().Lookup(srcType)
	if dlObj == nil {
		return nil, errors.New("package path does not exist")
	}
	return dlObj, nil
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
