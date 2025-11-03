package utils

import (
	"encoding/json"
	"fmt"
)

func PrettyPrint(v any) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println(string(b))
}
