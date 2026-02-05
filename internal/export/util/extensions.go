package util

import (
	"encoding/json"
	"maps"
)

// marshalWithExtensions marshals a struct with extensions inlined.
// This is a helper for custom MarshalJSON implementations.
//
// IMPORTANT: When calling this function, the caller MUST use a type alias
// to avoid infinite recursion. For example,
//
//	func (s *MyStruct) MarshalJSON() ([]byte, error) {
//	    type myStruct MyStruct  // Type alias prevents recursion
//	    return marshalWithExtensions(myStruct(*s), s.Extensions)
//	}
//
// Without the type alias, json.Marshal would recursively call MarshalJSON
// on the same type, causing infinite recursion. The type alias creates a
// new type that doesn't have the MarshalJSON method, allowing standard
// JSON marshaling to proceed.
func MarshalWithExtensions(v any, extensions map[string]any) ([]byte, error) {
	// Marshal the base struct
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	if len(extensions) == 0 {
		return data, nil
	}

	// Parse the JSON into a map
	var m map[string]any
	if unmarshalErr := json.Unmarshal(data, &m); unmarshalErr != nil {
		return nil, unmarshalErr
	}

	// Merge extensions into the map
	maps.Copy(m, extensions)

	// Marshal back to JSON
	return json.Marshal(m)
}
