package utils

import "encoding/json"

// MapToJSON converts a map to a json string
func MapToJSON(objmap map[string]string) string {
	b, _ := json.Marshal(objmap)
	return string(b)
}
