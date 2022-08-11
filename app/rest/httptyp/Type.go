package httptyp

import (
	"net/http"
	"net/url"
	"reflect"
)

type HttpRequestPathValues map[string]string

var HttpRequestPathValuesType reflect.Type = reflect.TypeOf((*HttpRequestPathValues)(nil)).Elem()

func (s HttpRequestPathValues) Get(key string) (string, bool) {
	v, e := s[key]
	return v, e
}

type HttpRequestQueryValues url.Values

var HttpRequestQueryValuesType reflect.Type = reflect.TypeOf((*HttpRequestQueryValues)(nil)).Elem()

func (s HttpRequestQueryValues) Get(key string) (string, bool) {
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

type HttpRequestHeaderValues http.Header

var HttpRequestHeaderValuesType reflect.Type = reflect.TypeOf((*HttpRequestHeaderValues)(nil)).Elem()

func (s HttpRequestHeaderValues) Get(key string) (string, bool) {
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

type HttpRequestPath string

var HttpRequestPathPtrType reflect.Type = reflect.TypeOf((*HttpRequestPath)(nil))

func ParseHttpRequestPath(s string) *HttpRequestPath {
	hrp := HttpRequestPath(s)
	return &hrp
}

func (s *HttpRequestPath) String() string {
	if s == nil {
		return ""
	}
	return string(*s)
}

type HttpRequestQuery string

var HttpRequestQueryPtrType reflect.Type = reflect.TypeOf((*HttpRequestQuery)(nil))

func ParseHttpRequestQuery(s string) *HttpRequestQuery {
	hrq := HttpRequestQuery(s)
	return &hrq
}

func (s *HttpRequestQuery) String() string {
	if s == nil {
		return ""
	}
	return string(*s)
}

type HttpRequestHeader string

var HttpRequestHeaderPtrType reflect.Type = reflect.TypeOf((*HttpRequestHeader)(nil))

func ParseHttpRequestHeader(s string) *HttpRequestHeader {
	hrh := HttpRequestHeader(s)
	return &hrh
}

func (s *HttpRequestHeader) String() string {
	if s == nil {
		return ""
	}
	return string(*s)
}

var HttpHandleFuncType = reflect.TypeOf((*http.HandlerFunc)(nil)).Elem()
var HttpRequestType = reflect.TypeOf((*http.Request)(nil))

var NullValue = reflect.ValueOf(nil)
