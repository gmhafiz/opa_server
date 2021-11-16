package opa

import (
	"fmt"
	"io/ioutil"
	"path"

	"github.com/open-policy-agent/opa/ast"
)

type Policy struct {
	Compiler *ast.Compiler
}

func New(filename string) (*Policy, error) {
	compiler, err := load(filename)

	return &Policy{Compiler: compiler}, err
}

func load(file string) (*ast.Compiler, error) {
	modules := map[string]string{}

	/* #nosec */
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("error load: %w", err)
	}

	modules[path.Base(file)] = string(content)

	return ast.CompileModules(modules)
}
