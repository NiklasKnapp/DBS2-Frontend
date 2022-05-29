package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

var myClient = &http.Client{Timeout: 10 * time.Second}

func GetJson(url string, target interface{}) error {
	r, err := myClient.Get(url)
	if err != nil {
		return fmt.Errorf("no response from request: %v", err)
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}
