package convert

import (
	"bytes"
	"encoding/json"
	"fmt"

	jsonParser "github.com/tmax-cloud/hcljson/parser"
	hclprinter "github.com/tmax-cloud/hcljson/printer"
	// hclprinter "github.com/hashicorp/hcl/hcl/printer"
	// jsonParser "github.com/hashicorp/hcl/json/parser"
)

func JsonToHcl(input []byte, typeSchemaStr string) []byte {

	var typeSchema map[string]interface{}
	json.Unmarshal([]byte(typeSchemaStr), &typeSchema)

	bytes, err := convertJsonToHcl(input, typeSchema)
	if err != nil {
		fmt.Errorf("hclTojson() error. %s", err)
	}
	return bytes
}

func convertJsonToHcl(input []byte, typeSchema map[string]interface{}) ([]byte, error) {
	ast, err := jsonParser.Parse(input)
	if err != nil {
		return nil, fmt.Errorf("unable to parse JSON: %s", err)
	}
	var buf bytes.Buffer
	if err := hclprinter.Fprint(&buf, ast, typeSchema); err != nil {
		return nil, fmt.Errorf("Unable to print HCL: %s", err)
	}

	return buf.Bytes(), nil
}
