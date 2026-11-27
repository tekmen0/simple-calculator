package httpapi

import (
	"bytes"
	"encoding/json"
	"io"
)

type operands struct {
	left  float64
	right float64
}

func decodeOperands(body []byte) (operands, bool) {
	// Decoding into a struct would accept duplicate keys and case-insensitive field names.
	decoder := json.NewDecoder(bytes.NewReader(body))
	if token, err := decoder.Token(); err != nil || token != json.Delim('{') {
		return operands{}, false
	}

	values := make(map[string]float64, 2)
	for decoder.More() {
		token, err := decoder.Token()
		field, _ := token.(string)
		if err != nil || (field != "left" && field != "right") {
			return operands{}, false
		}
		if _, exists := values[field]; exists {
			return operands{}, false
		}
		var value *float64
		if err := decoder.Decode(&value); err != nil || value == nil {
			return operands{}, false
		}
		values[field] = *value
	}
	if token, err := decoder.Token(); err != nil || token != json.Delim('}') || len(values) != 2 {
		return operands{}, false
	}
	if _, err := decoder.Token(); err != io.EOF {
		return operands{}, false
	}
	return operands{left: values["left"], right: values["right"]}, true
}
