// Copyright © 2026 The gRPC Kit Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package projectmigrate

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"strconv"
	"strings"

	"golang.org/x/mod/modfile"
	"golang.org/x/mod/semver"
)

const maxDiagnosticGoFileSize = 8 << 20

var printfLoggerMethods = map[string]struct{}{
	"Tracef": {}, "Debugf": {}, "Infof": {}, "Printf": {}, "Warnf": {}, "Warningf": {}, "Errorf": {}, "Fatalf": {}, "Panicf": {},
}

// ScanManualActions reports known pkg v0.5.0 incompatibilities without
// modifying source. Paths already replaced by the managed plan are excluded.
func ScanManualActions(project Project, changedPaths map[string]struct{}) ([]Diagnostic, error) {
	actions, err := scanGoModActions(project)
	if err != nil {
		return nil, err
	}
	dedup := make(map[string]struct{})
	for _, action := range actions {
		dedup[diagnosticKey(action)] = struct{}{}
	}
	err = filepath.WalkDir(project.Root, func(currentPath string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			if entry.Name() == ".git" || entry.Name() == "vendor" {
				return filepath.SkipDir
			}
			return nil
		}
		if entry.Type()&fs.ModeSymlink != 0 || !entry.Type().IsRegular() || !strings.HasSuffix(entry.Name(), ".go") {
			return nil
		}
		relativePath, err := filepath.Rel(project.Root, currentPath)
		if err != nil {
			return err
		}
		relativePath = filepath.ToSlash(relativePath)
		if _, excluded := changedPaths[relativePath]; excluded {
			return nil
		}
		body, err := readRegularFile(currentPath, maxDiagnosticGoFileSize)
		if err != nil {
			return fmt.Errorf("read %s: %w", relativePath, err)
		}
		fileSet := token.NewFileSet()
		parsed, err := parser.ParseFile(fileSet, relativePath, body, parser.SkipObjectResolution)
		if err != nil {
			appendDiagnostic(&actions, dedup, Diagnostic{Code: "inspect_source", Path: relativePath, Message: "fix Go syntax before migration diagnostics: " + err.Error()})
			return nil
		}
		scanGoFile(fileSet, parsed, relativePath, &actions, dedup)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("scan manual actions: %w", err)
	}
	sortDiagnostics(actions)
	return actions, nil
}

func scanGoModActions(project Project) ([]Diagnostic, error) {
	path := filepath.Join(project.Root, "go.mod")
	body, err := readRegularFile(path, maxGoModSize)
	if err != nil {
		return nil, err
	}
	parsed, err := modfile.Parse("go.mod", body, nil)
	if err != nil {
		return nil, err
	}
	var actions []Diagnostic
	goVersion := ""
	if parsed.Go != nil {
		goVersion = parsed.Go.Version
	}
	if goVersion == "" || semver.Compare("v"+goVersion, "v"+currentTargetGoVersion) < 0 {
		actions = append(actions, Diagnostic{Code: "go_version", Path: "go.mod", Line: lineContaining(body, "go "), Message: fmt.Sprintf("set the Go directive to at least %s (current %q)", currentTargetGoVersion, goVersion)})
	}
	pkgVersion := ""
	for _, requirement := range parsed.Require {
		switch requirement.Mod.Path {
		case "github.com/grpc-kit/pkg":
			pkgVersion = requirement.Mod.Version
		case "github.com/sirupsen/logrus":
			actions = append(actions, Diagnostic{
				Code:    "logrus_dependency",
				Path:    "go.mod",
				Line:    lineContaining(body, "github.com/sirupsen/logrus"),
				Message: "remove the direct logrus requirement after migrating source usage, then run go mod tidy",
			})
		}
	}
	if pkgVersion != currentTargetPkgVersion {
		actions = append(actions, Diagnostic{Code: "pkg_version", Path: "go.mod", Line: lineContaining(body, "github.com/grpc-kit/pkg"), Message: fmt.Sprintf("require github.com/grpc-kit/pkg %s (current %q)", currentTargetPkgVersion, pkgVersion)})
	}
	return actions, nil
}

func scanGoFile(fileSet *token.FileSet, file *ast.File, relativePath string, actions *[]Diagnostic, dedup map[string]struct{}) {
	logrusAliases := make(map[string]struct{})
	sdAliases := make(map[string]struct{})
	for _, imported := range file.Imports {
		importPath, err := strconv.Unquote(imported.Path.Value)
		if err != nil {
			continue
		}
		alias := filepath.Base(importPath)
		if imported.Name != nil {
			alias = imported.Name.Name
		}
		switch importPath {
		case "github.com/sirupsen/logrus":
			logrusAliases[alias] = struct{}{}
			appendDiagnostic(actions, dedup, Diagnostic{Code: "logrus_import", Path: relativePath, Line: fileSet.Position(imported.Pos()).Line, Message: "replace direct logrus usage with log/slog"})
		case "github.com/grpc-kit/pkg/sd":
			sdAliases[alias] = struct{}{}
		}
	}

	ast.Inspect(file, func(node ast.Node) bool {
		switch current := node.(type) {
		case *ast.StarExpr:
			if selector, ok := current.X.(*ast.SelectorExpr); ok && selector.Sel.Name == "Entry" && isAlias(selector.X, logrusAliases) {
				appendDiagnostic(actions, dedup, Diagnostic{Code: "logrus_type", Path: relativePath, Line: fileSet.Position(current.Pos()).Line, Message: "replace *logrus.Entry with *slog.Logger"})
			}
		case *ast.SelectorExpr:
			if isAlias(current.X, logrusAliases) && (current.Sel.Name == "Hook" || current.Sel.Name == "Formatter" || current.Sel.Name == "Fields") {
				appendDiagnostic(actions, dedup, Diagnostic{Code: "logrus_custom", Path: relativePath, Line: fileSet.Position(current.Pos()).Line, Message: "manually redesign logrus Hook/Formatter/Fields behavior for slog"})
			}
			if current.Sel.Name == "Data" && looksLikeLoggerExpression(current.X) {
				appendDiagnostic(actions, dedup, Diagnostic{Code: "logrus_custom", Path: relativePath, Line: fileSet.Position(current.Pos()).Line, Message: "replace direct logrus Entry.Data access with explicit slog attributes"})
			}
		case *ast.CallExpr:
			scanCall(fileSet, current, relativePath, logrusAliases, sdAliases, actions, dedup)
		case *ast.FuncDecl:
			if current.Name.Name == "Init" && receiverTypeName(current) == "IndependentCfg" && !firstParameterIsContext(current.Type.Params) {
				appendDiagnostic(actions, dedup, Diagnostic{Code: "context_signature", Path: relativePath, Line: fileSet.Position(current.Pos()).Line, Message: "change IndependentCfg.Init to accept context.Context as its first parameter"})
			}
		}
		return true
	})
}

func scanCall(fileSet *token.FileSet, call *ast.CallExpr, relativePath string, logrusAliases, sdAliases map[string]struct{}, actions *[]Diagnostic, dedup map[string]struct{}) {
	selector, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return
	}
	line := fileSet.Position(call.Pos()).Line
	name := selector.Sel.Name
	if _, ok := printfLoggerMethods[name]; ok && looksLikeLoggerExpression(selector.X) {
		appendDiagnostic(actions, dedup, Diagnostic{Code: "printf_logger", Path: relativePath, Line: line, Message: fmt.Sprintf("replace %s with a context-aware slog call and structured attributes", name)})
	}
	if (name == "WithFields" || name == "WithField" || name == "WithError") && (looksLikeLoggerExpression(selector.X) || isAlias(selector.X, logrusAliases)) {
		appendDiagnostic(actions, dedup, Diagnostic{Code: "logrus_custom", Path: relativePath, Line: line, Message: fmt.Sprintf("manually translate %s to slog attributes", name)})
	}
	switch {
	case name == "NewMicroservice" && len(call.Args) == 1:
		appendDiagnostic(actions, dedup, Diagnostic{Code: "context_call", Path: relativePath, Line: line, Message: "pass context.Context as the first NewMicroservice argument"})
	case name == "Init" && len(call.Args) == 0:
		appendDiagnostic(actions, dedup, Diagnostic{Code: "context_call", Path: relativePath, Line: line, Message: "verify this Init call and pass context.Context when it targets grpc-kit configuration"})
	case name == "Deregister" && len(call.Args) == 0:
		appendDiagnostic(actions, dedup, Diagnostic{Code: "context_call", Path: relativePath, Line: line, Message: "pass context.Context to Deregister"})
	case name == "Register" && isAlias(selector.X, sdAliases) && len(call.Args) == 5:
		appendDiagnostic(actions, dedup, Diagnostic{Code: "context_call", Path: relativePath, Line: line, Message: "pass context.Context as the first sd.Register argument"})
	case name == "HTTPHandlerFrontend" && len(call.Args) == 2:
		appendDiagnostic(actions, dedup, Diagnostic{Code: "context_call", Path: relativePath, Line: line, Message: "pass context.Context as the first HTTPHandlerFrontend argument"})
	case name == "StartBackground" && len(call.Args) == 0:
		appendDiagnostic(actions, dedup, Diagnostic{Code: "context_call", Path: relativePath, Line: line, Message: "pass context.Context to StartBackground"})
	case name == "WithLogger" && len(call.Args) == 3:
		appendDiagnostic(actions, dedup, Diagnostic{Code: "errs_with_logger", Path: relativePath, Line: line, Message: "pass context.Context before the slog logger to errs.Status.WithLogger"})
	}
}

func looksLikeLoggerExpression(expression ast.Expr) bool {
	switch current := expression.(type) {
	case *ast.Ident:
		name := strings.ToLower(current.Name)
		return strings.Contains(name, "logger") || name == "log" || strings.Contains(name, "entry")
	case *ast.SelectorExpr:
		name := strings.ToLower(current.Sel.Name)
		return strings.Contains(name, "logger") || name == "log" || strings.Contains(name, "entry")
	case *ast.CallExpr:
		selector, ok := current.Fun.(*ast.SelectorExpr)
		if !ok {
			return false
		}
		name := selector.Sel.Name
		return name == "WithFields" || name == "WithField" || name == "WithError" || strings.Contains(strings.ToLower(name), "logger")
	default:
		return false
	}
}

func isAlias(expression ast.Expr, aliases map[string]struct{}) bool {
	identifier, ok := expression.(*ast.Ident)
	if !ok {
		return false
	}
	_, ok = aliases[identifier.Name]
	return ok
}

func receiverTypeName(function *ast.FuncDecl) string {
	if function.Recv == nil || len(function.Recv.List) == 0 {
		return ""
	}
	typeExpression := function.Recv.List[0].Type
	if pointer, ok := typeExpression.(*ast.StarExpr); ok {
		typeExpression = pointer.X
	}
	if identifier, ok := typeExpression.(*ast.Ident); ok {
		return identifier.Name
	}
	return ""
}

func firstParameterIsContext(parameters *ast.FieldList) bool {
	if parameters == nil || len(parameters.List) == 0 {
		return false
	}
	selector, ok := parameters.List[0].Type.(*ast.SelectorExpr)
	return ok && selector.Sel.Name == "Context"
}

func appendDiagnostic(actions *[]Diagnostic, dedup map[string]struct{}, action Diagnostic) {
	key := diagnosticKey(action)
	if _, exists := dedup[key]; exists {
		return
	}
	dedup[key] = struct{}{}
	*actions = append(*actions, action)
}

func diagnosticKey(diagnostic Diagnostic) string {
	return fmt.Sprintf("%s:%d:%s", diagnostic.Path, diagnostic.Line, diagnostic.Code)
}

func lineContaining(body []byte, needle string) int {
	for index, line := range strings.Split(string(body), "\n") {
		if strings.Contains(line, needle) {
			return index + 1
		}
	}
	return 1
}
