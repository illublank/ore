package httptyp

import (
	"reflect"
)

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

func (s *RequestQuery) String() string {
	return string(*s)
}

type RequestHeader string

var RequestHeaderPtrType reflect.Type = reflect.TypeOf((*RequestHeader)(nil))

func (s *RequestHeader) String() string {
	return string(*s)
}
