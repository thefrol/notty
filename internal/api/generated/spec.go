// Package generated provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen/v2 version v2.0.0 DO NOT EDIT.
package generated

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"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",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %w", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	res := make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	resolvePath := PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		pathToFile := url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
