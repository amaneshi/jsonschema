package jsonschema

import (
	"encoding/json"
	"errors"
	"strings"
)

const typeSeparator = ";"

type Type string

func (t Type) MarshalJSON() ([]byte, error) {
	str := string(t)
	if strings.Contains(str, typeSeparator) {
		// Split by ";" and filter out empty pieces
		parts := strings.Split(str, typeSeparator)
		nonEmptyParts := make([]string, 0, len(parts))
		for _, p := range parts {
			if p != "" {
				nonEmptyParts = append(nonEmptyParts, p)
			}
		}
		return json.Marshal(nonEmptyParts)
	}
	// No ";" found, marshal as a plain string
	return json.Marshal(str)
}

func (t *Type) UnmarshalJSON(data []byte) error {
	// Try to unmarshal as string first
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		*t = Type(str)
		return nil
	}

	// Otherwise try to unmarshal as []string
	var arr []string
	if err := json.Unmarshal(data, &arr); err == nil {
		// Join non-empty elements with ";"
		nonEmptyParts := make([]string, 0, len(arr))
		for _, p := range arr {
			if p != "" {
				nonEmptyParts = append(nonEmptyParts, p)
			}
		}
		*t = Type(strings.Join(nonEmptyParts, typeSeparator))
		return nil
	}

	return errors.New("type: unsupported JSON format")
}
