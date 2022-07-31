package handler

import (
  "fmt"
  "net/http"

  "github.com/illublank/go-common/log"
)

// DebugHandler todo
type DebugHandler struct {
  http.Handler
  Logger         log.Logger
  OrginalHandler http.Handler
}

func (s *DebugHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  s.Logger.Debug(fmt.Sprintf("%v %v %v %v", r.Method, r.URL.Path, r.URL.RawQuery, r.Header))
  s.OrginalHandler.ServeHTTP(w, r)
}
