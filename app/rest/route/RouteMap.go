package route

import (
  "net/http"
  "reflect"
  "runtime"
  "sort"
  "strings"
)

/*
	RouteMap {
		"${requestPath}":RouteHandler {
			"${requestMethod}": func(w http.ResponseWriter, r *http.Request)
		}
	}
*/
// RouteMethod todo
type RouteMethod string

// RouteItem todo
type RouteItem struct {
  Path       string
  Method     string
  HandleFunc func(http.ResponseWriter, *http.Request)
  AutoDIFunc interface{}
}

// RouteHandler todo
type RouteHandler map[string]func(w http.ResponseWriter, r *http.Request)

func (s RouteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  if handler, ok := s[r.Method]; ok {
    handler(w, r)
  } else {
    if r.Method == "OPTIONS" {
      allow := []string{}
      for k := range s {
        allow = append(allow, k)
      }
      sort.Strings(allow)
      w.Header().Set("Allow", strings.Join(allow, ", "))
      w.WriteHeader(http.StatusOK)
    } else {
      http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    }
  }
}

func (s RouteHandler) String() string {
  arr := []string{}
  for k, v := range s {
    arr = append(arr, "\""+k+"\": \""+runtime.FuncForPC(reflect.ValueOf(v).Pointer()).Name()+"\"")
  }
  return "{" + strings.Join(arr, ",") + "}"
}

// AllMethodHandler todo
type AllMethodHandler func(http.ResponseWriter, *http.Request)

func (s AllMethodHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  s(w, r)
}

func (s AllMethodHandler) String() string {
  return "{" + "\"ALL\": \"" + runtime.FuncForPC(reflect.ValueOf(s).Pointer()).Name() + "\"" + "}"
}

// RouteMap todo
type RouteMap map[string]http.Handler

func (s RouteMap) GetRouteHandler(path string) (RouteHandler, bool) {
  if h, exists := s[path]; !exists {
    return nil, false
  } else if rh, ok := h.(RouteHandler); ok {
    return rh, true
  } else {
    return nil, false
  }
}

func (s RouteMap) Add(item *RouteItem) RouteMap {

  if len(item.Method) != 0 {
    if rh, exists := s.GetRouteHandler(item.Path); !exists {
      s[item.Path] = RouteHandler{item.Method: item.HandleFunc}
    } else {
      rh[item.Method] = item.HandleFunc
    }
  } else {
    s[item.Path] = AllMethodHandler(item.HandleFunc)
  }

  return s
}

// NewRouteMap todo
func NewRouteMap(items ...*RouteItem) RouteMap {
  result := make(RouteMap)
  for _, item := range items {
    result.Add(item)
  }
  return result
}
