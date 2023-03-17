package api

import (
  "encoding/json"
  "net/http"
)

var jsonContentType = []string{"application/json;charset=utf-8"}

const (
  MESSAGE_INTERNAL_SERVER_ERROR = "Internal Server Error"
  MESSAGE_FILE_EXCEED_LIMIT     = "File Exceeded Limit"
)

func WriteJSON(w http.ResponseWriter, code int, obj interface{}) error {
  w.WriteHeader(code)
  writeContentType(w, jsonContentType)
  jsonBytes, err := json.Marshal(obj)
  if err != nil {
    return err
  }
  w.Write(jsonBytes)
  return nil
}
func writeContentType(w http.ResponseWriter, value []string) {
  header := w.Header()
  if val := header["Content-Type"]; len(val) == 0 {
    header["Content-type"] = value
  }
}