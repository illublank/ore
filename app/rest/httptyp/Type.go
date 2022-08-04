package httptyp

import (
  "net/http"
  "net/url"
  "reflect"
)

type PathVals map[string]string

var PathValsType reflect.Type = reflect.TypeOf((*PathVals)(nil)).Elem()

func (s PathVals) Get(key string) (string, bool) {
  v, e := s[key]
  return v, e
}

type QueryVals url.Values

var QueryValsType reflect.Type = reflect.TypeOf((*QueryVals)(nil)).Elem()

func (s QueryVals) Get(key string) (string, bool) {
  v, e := s[key]
  if e {
    if len(v) > 0 {
      return v[0], true
    } else {
      return "", true
    }
  }
  return "", false
}

type HeaderVals http.Header

var HeaderValsType reflect.Type = reflect.TypeOf((*HeaderVals)(nil)).Elem()

func (s HeaderVals) Get(key string) (string, bool) {
  v, e := s[key]
  if e {
    if len(v) > 0 {
      return v[0], true
    } else {
      return "", true
    }
  }
  return "", false
}

type RequestPath string

var RequestPathPtrType reflect.Type = reflect.TypeOf((*RequestPath)(nil))

func ParseRequestPath(s string) *RequestPath {
  rp := RequestPath(s)
  return &rp
}

func (s *RequestPath) String() string {
  if s == nil {
    return ""
  }
  return string(*s)
}

type RequestQuery string

var RequestQueryPtrType reflect.Type = reflect.TypeOf((*RequestQuery)(nil))

func ParseRequestQuery(s string) *RequestQuery {
  rp := RequestQuery(s)
  return &rp
}

func (s *RequestQuery) String() string {
  return string(*s)
}

type RequestHeader string

var RequestHeaderPtrType reflect.Type = reflect.TypeOf((*RequestHeader)(nil))

func ParseRequestHeader(s string) *RequestHeader {
  rp := RequestHeader(s)
  return &rp
}

func (s *RequestHeader) String() string {
  return string(*s)
}
