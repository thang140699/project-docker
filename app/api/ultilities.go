package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ResponseBody struct {
	Message string      `json:"message"`
	Code    int         `json:"code"`
	Data    interface{} `json:"data,omitempty"`
}

func BindJSON(r *http.Request, obj interface{}) error {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return json.Unmarshal(b, obj)
}
