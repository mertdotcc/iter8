package base

import (
	"io"
	"net/http"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"

	"sigs.k8s.io/yaml"

	"github.com/Masterminds/sprig"
	"github.com/iter8-tools/iter8/base/log"
)

// MajorMinor is the minor version of Iter8
// set this manually whenever the major or minor version changes
var MajorMinor = "v0.13"

// Version is the semantic version of Iter8 (with the `v` prefix)
// Version is intended to be set using LDFLAGS at build time
var Version = "v0.13.0"

const (
	toYAMLString = "toYaml"
)

// int64Pointer takes an int64 as input, creates a new variable with the input value, and returns a pointer to the variable
func int64Pointer(i int64) *int64 {
	return &i
}

// intPointer takes an int as input, creates a new variable with the input value, and returns a pointer to the variable
func intPointer(i int) *int {
	return &i
}

// float32Pointer takes an float32 as input, creates a new variable with the input value, and returns a pointer to the variable
func float32Pointer(f float32) *float32 {
	return &f
}

// float64Pointer takes an float64 as input, creates a new variable with the input value, and returns a pointer to the variable
func float64Pointer(f float64) *float64 {
	return &f
}

// StringPointer takes string as input, creates a new variable with the input value, and returns a pointer to the variable
func StringPointer(s string) *string {
	return &s
}

// BoolPointer takes bool as input, creates a new variable with the input value, and returns a pointer to the variable
func BoolPointer(b bool) *bool {
	return &b
}

// CompletePath is a helper function for converting file paths, specified relative to the caller of this function, into absolute ones.
// CompletePath is useful in tests and enables deriving the absolute path of experiment YAML files.
func CompletePath(prefix string, suffix string) string {
	_, filename, _, _ := runtime.Caller(1) // one step up the call stack
	return filepath.Join(filepath.Dir(filename), prefix, suffix)
}

// getTextTemplateFromURL gets template from URL
func getTextTemplateFromURL(providerURL string) (*template.Template, error) {
	// fetch b from url
	// #nosec
	resp, err := http.Get(providerURL)
	if err != nil {
		log.Logger.Error(err)
		return nil, err
	}
	// read responseBody
	// get the doubly templated metrics spec
	defer func() {
		_ = resp.Body.Close()
	}()
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	tpl, err := CreateTemplate(string(responseBody))
	if err != nil {
		log.Logger.Error(err)
		return nil, err
	}

	return tpl, nil
}

// CreateTemplate creates a template from a string
func CreateTemplate(tplString string) (*template.Template, error) {
	return template.New("provider template").Funcs(FuncMapWithToYAML()).Parse(string(tplString))
}

// FuncMapWithToYAML return sprig text function map with a toYaml function
func FuncMapWithToYAML() template.FuncMap {
	f := sprig.TxtFuncMap()
	f[toYAMLString] = ToYAML

	return f
}

// ToYAML takes an interface, marshals it to yaml, and returns a string. It will
// always return a string, even on marshal error (empty string).
//
// This is designed to be called from a template.
func ToYAML(v interface{}) string {
	data, err := yaml.Marshal(v)
	if err != nil {
		// Swallow errors inside of a template.
		return ""
	}
	return strings.TrimSuffix(string(data), "\n")
}
