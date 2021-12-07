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
	Name           string
	NameLowerCase  string
	Request        string
	Response       string
	IsStreamServer bool
	IsStreamClient bool
}

// service contains all the necessary information to generate the files
type service struct {
	Name          string
	NameLowerCase string
	Methods       []method
}

type protoDefinition struct {
	Package    string
	ImportPath string
	GoPackage  string
	services   []service
}

func (p protoDefinition) dump() {
	log.Printf("Package: %s", p.Package)
	log.Printf("ImportPath: %s", p.ImportPath)
	log.Printf("GoPackage: %s", p.GoPackage)

	for _, s := range p.services {
		log.Printf("Name: %s", s.Name)
		log.Printf("NameLowerCase: %s", s.NameLowerCase)
		for _, m := range s.Methods {
			log.Printf("Name: %s", m.Name)
			log.Printf("NameLowerCase: %s", m.NameLowerCase)
			log.Printf("Request: %s", m.Request)
			log.Printf("Response: %s", m.Response)
			log.Printf("IsStreamServer: %t", m.IsStreamServer)
			log.Printf("IsStreamClient: %t", m.IsStreamClient)
			log.Printf("-------------------------------------------------------")
		}
	}
}

var generatedFiles = map[string]string{
	"./Makefile":                "templates/Makefile.tmpl",
	"./go.mod":                  "templates/go.mod.tmpl",
	"./main.go":                 "templates/main.go.tmpl",
	"./cmd/client.go":           "templates/cmd/client.go.tmpl",
	"./cmd/root.go":             "templates/cmd/root.go.tmpl",
	"./cmd/run.go":              "templates/cmd/run.go.tmpl",
	"./cmd/server.go":           "templates/cmd/server.go.tmpl",
	"./common/logger/logger.go": "templates/common/logger/logger.go.tmpl",
	"./common/logger/zap.go":    "templates/common/logger/zap.go.tmpl",
	"./common/common.go":        "templates/common/common.go.tmpl",
	"./service/service.go":      "templates/service/service.go.tmpl",
	"./service/grpc/crl.go":     "templates/service/grpc/crl.go.tmpl",
	"./service/grpc/server.go":  "templates/service/grpc/server.go.tmpl",
	"./service/grpc/tls.go":     "templates/service/grpc/tls.go.tmpl",
	"./service/http/server.go":  "templates/service/http/server.go.tmpl",
	//"./cmd/{svc}/main.go":    "templates/cmd/main.tmpl",
	//"./pkg/server/server.go": "templates/server/server.tmpl",
}

func firstLetterToLowerCase(s string) string {
	l := byte(unicode.ToLower(rune(s[0])))

	return string(append([]byte{l}, []byte(s)[1:]...))
}

func main() {
	//_p, err := parseProto()
	_p, err := parseProto()
	if err != nil {
		panic(fmt.Errorf("failed to parse proto file %w", err))
	}

	//_p.dump()
	for f, t := range generatedFiles {
		if err := createFileFromTemplate(_p.Package, _p.services[0], f, t); err != nil {
			panic(fmt.Errorf("failed to create file %v: %w", f, err))
		}
	}

	// create the folder for go_out and grpc_out
	if err := os.MkdirAll("./pkg/pb", 0775); err != nil {
		if errors.Is(err, os.ErrExist) {
			panic(err)
		}
	}

}

func parseProto() (protoDefinition, error) {
	var err error
	p := protoDefinition{}

	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		return p, fmt.Errorf("io.ReadAll(): %w", err)
	}

	req := pluginpb.CodeGeneratorRequest{}
	if err = proto.Unmarshal(input, &req); err != nil {
		return p, fmt.Errorf("proto.Unmarshal(): %w", err)
	}

	opts := protogen.Options{}

	plugin, err := opts.New(&req)

	if err != nil {
		return p, fmt.Errorf("protogen.Options.New(): %w", err)
	}

	// Files appear in topological order, so each file appears before any
	// file that imports it.
	f := plugin.Files[len(plugin.Files)-1]

	p.Package = string(*f.Proto.Package)
	p.GoPackage = string(f.GoPackageName)
	p.ImportPath = filepath.Dir(string(f.GoImportPath))

	for _, s := range f.Services {
		svc := service{}
		svc.Name = strings.Title(s.GoName)
		svc.NameLowerCase = firstLetterToLowerCase(svc.Name)

		for _, m := range f.Services[0].Methods {
			var svcMethod method
			svcMethod, err = parseMethod(m)

			if err != nil {
				return p, err
			}
			svc.Methods = append(svc.Methods, svcMethod)
		}
		p.services = append(p.services, svc)
	}

	return p, err
}

func parseMethod(m *protogen.Method) (svcMethod method, err error) {

	svcMethod.Name = strings.Title(m.GoName)
	svcMethod.NameLowerCase = firstLetterToLowerCase(svcMethod.Name)
	svcMethod.Request = m.Input.GoIdent.GoName
	svcMethod.Response = m.Output.GoIdent.GoName
	svcMethod.IsStreamClient = m.Desc.IsStreamingClient()
	svcMethod.IsStreamServer = m.Desc.IsStreamingServer()

	return
}

func createFileFromTemplate(packageName string, svc service, filePath, templatePath string) (err error) {
	filePath = strings.Replace(filePath, "{svc}", svc.NameLowerCase, -1)

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

	funcMap := template.FuncMap{
		"ToLower": strings.ToLower,
	}

	t := template.New(templatePath).Funcs(funcMap)
	t, err = t.Parse(string(b))
	if err != nil {
		return fmt.Errorf("unable to parse template file: %w", err)
	}

	data := struct {
		PackageName string
		Service     service
	}{PackageName: packageName, Service: svc}

	if err = t.Execute(f, data); err != nil {
		return fmt.Errorf("error executing template %v:%w", t.Name(), err)
	}

	return
}
