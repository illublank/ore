package handler

import (
  "net/http"
)

// ResourceHandler todo
func ResourceHandler(filepath string) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, filepath)
  })
}
