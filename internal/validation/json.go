package validation

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/xeipuuv/gojsonschema"
)

const (
	PACK_SCHEMA = "file://api/pack_schema.json"
)

func ValidateAgainstSchema(schemaPath string, validatee json.RawMessage) bool {
	schemaLoader := gojsonschema.NewReferenceLoader(schemaPath)
	validateeLoader := gojsonschema.NewBytesLoader(validatee)
	res, err := gojsonschema.Validate(schemaLoader, validateeLoader)
	if err != nil {
		if err != io.EOF {
			fmt.Fprintf(os.Stderr, "Validating json against schema went wrong: %v", err)
		}
		return false
	}
	return res.Valid()
}
