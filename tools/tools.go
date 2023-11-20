//go:build tools
// +build tools

// This file is not intended to be compiled.  Because some of these imports are
// not actual go packages, we use a build constraint at the top of this file to
// prevent tools from inspecting the imports.

package tools

import (
	_ "github.com/client9/misspell/cmd/misspell"
	_ "github.com/gogo/protobuf/protoc-gen-gogoslick"
	_ "golang.org/x/lint/golint"
	_ "golang.org/x/tools/cmd/goimports"
	_ "honnef.co/go/tools/cmd/staticcheck"
)
