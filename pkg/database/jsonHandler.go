package database

import (
	"encoding/json"
)

// DBDataHandler is an interface for encoding/decoding a DBFile between object and a data format
type DBDataHandler interface {
	Encode(input interface{}) ([]byte, error)
	Decode(input []byte, output interface{}) error
}

// NewDBDataHandler is a factory function for the DBDataHandler interface
func NewDBDataHandler(dataType string) DBDataHandler {
	switch {
	case dataType == "json":
		return new(JSONDataHandler)
	default:
		return nil
	}
}

// JSONDataHandler encodes/decodes a dbFile to/from JSON format
type JSONDataHandler struct{}

// Encode marshals the given input DBFile to a JSON []byte
func (h *JSONDataHandler) Encode(input interface{}) ([]byte, error) {
	return json.Marshal(input)
}

// Decode unmarshals the given JSON []byte to a DBFile
func (h *JSONDataHandler) Decode(input []byte, output interface{}) error {
	return json.Unmarshal(input, output)
}
