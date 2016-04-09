package common

import (
	"encoding/json"
	"fmt"
)

func ToJSON(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		fmt.Println("error:", err)
		return ""
	}

	return string(b)
}
