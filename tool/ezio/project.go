package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"text/template"
)

// project project config
type project struct {
	Name     string
	Owner    string
	Path     string
	WithGRPC bool
	Here     bool
}

const (
	_tplTypeChangeLog = iota
	_tplTypeContributors
	_tplTypeReadme
	_tplTypeGomod
	_tplTypeServer
	_tplTypeHandler
	_tplTypeModelDefault
	_tplTypeModel
	_tplTypeProto
	_tplTypeAppJson
	_tplTypeMakefile
	_tplTypeDockerfile
	_tplTypeMain
)

var (
	p project
	// files type => path
	files = map[int]string{
		// init doc
		_tplTypeChangeLog:    "/CHANGELOG.md",
		_tplTypeContributors: "/CONTRIBUTORS.md",
		_tplTypeReadme:       "/README.md",
		// init project
		_tplTypeMain:         "/cmd/main.go",
		_tplTypeGomod:        "/go.mod",
		_tplTypeServer:       "/internal/server.go",
		_tplTypeHandler:      "/internal/handler.go",
		_tplTypeModelDefault: "/model/default.go",
		_tplTypeModel:        "/model/user.go",
		// init proto
		_tplTypeProto: "/proto/user.proto",
		// init config
		_tplTypeAppJson:    "/app.json",
		_tplTypeMakefile:   "/Makefile",
		_tplTypeDockerfile: "/Dockerfile",
	}
	// tpls type => content
	tpls = map[int]string{
		_tplTypeChangeLog:    _tplChangeLog,
		_tplTypeContributors: _tplContributors,
		_tplTypeReadme:       _tplReadme,
		_tplTypeMain:         _tplMain,
		_tplTypeGomod:        _tplGoMod,
		_tplTypeServer:       _tplServer,
		_tplTypeHandler:      _tplHandler,
		_tplTypeModelDefault: _tplModelDefault,
		_tplTypeModel:        _tplModel,
		_tplTypeProto:        _tplProto,
		_tplTypeAppJson:      _tplAppJson,
		_tplTypeMakefile:     _tplMakefile,
		_tplTypeDockerfile:   _tplDockerfile,
	}
)

func create() (err error) {
	if err = os.MkdirAll(p.Path, 0755); err != nil {
		return
	}
	for t, v := range files {
		i := strings.LastIndex(v, "/")
		if i > 0 {
			dir := v[:i]
			if err = os.MkdirAll(p.Path+dir, 0755); err != nil {
				return
			}
		}
		if err = write(p.Path+v, tpls[t]); err != nil {
			return
		}
	}
	if p.WithGRPC {
		if err = genpb(); err != nil {
			return
		}
	}
	return
}

func genpb() error {
	cmd := exec.Command("kratos", "tool", "protoc", p.Name+"/api/api.proto")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func write(name, tpl string) (err error) {
	data, err := parse(tpl)
	if err != nil {
		return
	}
	return ioutil.WriteFile(name, data, 0644)
}

func parse(s string) ([]byte, error) {
	t, err := template.New("").Parse(s)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	if err = t.Execute(&buf, p); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
