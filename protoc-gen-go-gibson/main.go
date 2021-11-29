package main

import (
	"embed"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"unicode"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"
)

// templateDir is the directory that createFile will look
// for the given template
//go:embed templates
var templateDir embed.FS

type method struct {
	Name          string
	NameLowerCase string
	Request       string
	Response      string
}

// service contains all the necessary information to generate the files
type serivce struct {
	ImportPath           string
	Package              string
	ServiceName          string
	ServiceNameLowerCase string
	Methods              []method
}

var generatedFiles = map[string]string{
	"./cmd/{svc}/main.go":    "templates/cmd/main.tmpl",
	"./pkg/server/server.go": "templates/server/server.tmpl",
}

func firstLetterToLowerCase(s string) string {
	l := byte(unicode.ToLower(rune(s[0])))

	return string(append([]byte{l}, []byte(s)[1:]...))
}

func main() {
	svc, err := parseProto()
	if err != nil {
		panic(fmt.Errorf("failed to parse proto file %w", err))
	}

	for f, t := range generatedFiles {
		if err := createFileFromTemplate(svc, f, t); err != nil {
			panic(fmt.Errorf("failed to create file %v: %w", f, err))
		}
	}

	// create the folder for go_out and grpc_out
	if err := os.Mkdir("./pkg/pb", 0775); err != nil {
		if errors.Is(err, os.ErrExist) {
			panic(err)
		}
	}
}

func parseProto() (svc serivce, err error) {
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		return svc, fmt.Errorf("io.ReadAll(): %w", err)
	}

	log.Println(string(input))

	req := pluginpb.CodeGeneratorRequest{}
	if err = proto.Unmarshal(input, &req); err != nil {
		return svc, fmt.Errorf("proto.Unmarshal(): %w", err)
	}

	opts := protogen.Options{}

	plugin, err := opts.New(&req)

	if err != nil {
		return svc, fmt.Errorf("protogen.Options.New(): %w", err)
	}

	// Files appear in topological order, so each file appears before any
	// file that imports it.
	f := plugin.Files[len(plugin.Files)-1]

	log.Println(f.GoPackageName)
	log.Println(f.GeneratedFilenamePrefix)
	log.Println(f.GoDescriptorIdent.String())
	log.Println(f.Desc.FullName().Name())
	log.Println(f.Desc.ParentFile().Path())
	log.Println(f.Desc.Path())
	log.Println(f.GoImportPath.String())
	log.Println(f.Desc.Name())

	svc.Package = string(f.GoPackageName)
	svc.ImportPath = filepath.Dir(string(f.GoImportPath))
	svc.ServiceName = strings.Title(f.Services[0].GoName)
	svc.ServiceNameLowerCase = firstLetterToLowerCase(svc.ServiceName)

	for _, m := range f.Services[0].Methods {
		var svcMethod method

		svcMethod, err = parseMethod(m)
		if err != nil {
			return
		}

		svc.Methods = append(svc.Methods, svcMethod)
	}

	return
}

func parseMethod(m *protogen.Method) (svcMethod method, err error) {

	svcMethod.Name = strings.Title(m.GoName)
	svcMethod.NameLowerCase = firstLetterToLowerCase(svcMethod.Name)
	svcMethod.Request = m.Input.GoIdent.GoName
	svcMethod.Response = m.Output.GoIdent.GoName

	return
}

func createFileFromTemplate(svc serivce, filePath, templatePath string) (err error) {
	filePath = strings.Replace(filePath, "{svc}", svc.ServiceNameLowerCase, -1)

	if err = os.MkdirAll(filepath.Dir(filePath), 0775); err != nil {
		return fmt.Errorf("os.MkdirAll() %w", err)
	}

	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("os.Create(): %w", err)
	}
	defer f.Close()

	b, err := templateDir.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("unable to read template file: %w", err)
	}

	t := template.New(templatePath)
	t, err = t.Parse(string(b))
	if err != nil {
		return fmt.Errorf("unable to parse template file: %w", err)
	}

	if err = t.Execute(f, svc); err != nil {
		return fmt.Errorf("error executing template %v:%w", t.Name(), err)
	}

	return
}
