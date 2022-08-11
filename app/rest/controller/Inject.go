package controller

import (
	"encoding/json"
	"io"
	"net/http"
	"reflect"

	"github.com/illublank/ore/app/rest/httptyp"
)

type TypeReflectValue struct {
	OBJ_HttpResponseWriter      http.ResponseWriter
	OBJ_HttpRequest             *http.Request
	OBJ_HttpRequestPathValues   httptyp.HttpRequestPathValues
	OBJ_HttpRequestQueryValues  httptyp.HttpRequestQueryValues
	OBJ_HttpRequestHeaderValues httptyp.HttpRequestHeaderValues

	RV_HttpResponseWriter      reflect.Value
	RV_HttpRequest             reflect.Value
	RV_HttpRequestPathValues   reflect.Value
	RV_HttpRequestQueryValues  reflect.Value
	RV_HttpRequestHeaderValues reflect.Value
}

func GetNullValue(trv *TypeReflectValue) (reflect.Value, error) {
	return httptyp.NullValue, nil
}

func GetHttpResponseWriterRV(trv *TypeReflectValue) (reflect.Value, error) {
	return trv.RV_HttpResponseWriter, nil
}

func GetHttpRequestRV(trv *TypeReflectValue) (reflect.Value, error) {
	return trv.RV_HttpRequest, nil
}

func GetHttpRequestPathValuesRV(trv *TypeReflectValue) (reflect.Value, error) {
	return trv.RV_HttpRequestPathValues, nil
}

func GetHttpRequestQueryValuesRV(trv *TypeReflectValue) (reflect.Value, error) {
	return trv.RV_HttpRequestQueryValues, nil
}

func GetHttpRequestHeaderValuesRV(trv *TypeReflectValue) (reflect.Value, error) {
	return trv.RV_HttpRequestHeaderValues, nil
}

func BuildGetHttpRequestPathFunc(key string) func(*TypeReflectValue) (reflect.Value, error) {
	return func(trv *TypeReflectValue) (reflect.Value, error) {
		v, e := trv.OBJ_HttpRequestPathValues.Get(key)
		if e {
			return reflect.ValueOf(httptyp.ParseHttpRequestPath(v)), nil
		} else {
			return httptyp.NullValue, nil
		}
	}
}

func BuildGetHttpRequestQueryFunc(key string) func(*TypeReflectValue) (reflect.Value, error) {
	return func(trv *TypeReflectValue) (reflect.Value, error) {
		v, e := trv.OBJ_HttpRequestQueryValues.Get(key)
		if e {
			return reflect.ValueOf(httptyp.ParseHttpRequestQuery(v)), nil
		} else {
			return httptyp.NullValue, nil
		}
	}
}

func BuildGetHttpRequestHeaderFunc(key string) func(*TypeReflectValue) (reflect.Value, error) {
	return func(trv *TypeReflectValue) (reflect.Value, error) {
		v, e := trv.OBJ_HttpRequestHeaderValues.Get(key)
		if e {
			return reflect.ValueOf(httptyp.ParseHttpRequestHeader(v)), nil
		} else {
			return httptyp.NullValue, nil
		}
	}
}

func BuildGetHttpRequestBodyFunc(t reflect.Type) func(trv *TypeReflectValue) (reflect.Value, error) {
	return func(trv *TypeReflectValue) (reflect.Value, error) {
		bs, err := io.ReadAll(trv.OBJ_HttpRequest.Body)
		if err != nil {
			return httptyp.NullValue, err
		}
		newObj := reflect.New(t)
		err = json.Unmarshal(bs, newObj.Interface())
		if err != nil {
			return httptyp.NullValue, err
		}
		return newObj.Elem(), nil
	}
}

func BuildStructFunc(typ reflect.Type) func(*TypeReflectValue) (reflect.Value, error) {
	// structValue := reflect.New(pTyp).Elem()
	n := typ.NumField()
	funcs := make([]func(*TypeReflectValue) (reflect.Value, error), n)
	for i := 0; i < n; i++ {
		f := typ.Field(i)
		t := f.Type
		switch {
		case t.String() == "http.ResponseWriter":
			funcs[i] = GetHttpResponseWriterRV
		case httptyp.HttpRequestType.ConvertibleTo(t):
			funcs[i] = GetHttpRequestRV
		case t.AssignableTo(httptyp.HttpRequestPathValuesType):
			funcs[i] = GetHttpRequestPathValuesRV
		case t.AssignableTo(httptyp.HttpRequestQueryValuesType):
			funcs[i] = GetHttpRequestQueryValuesRV
		case t.AssignableTo(httptyp.HttpRequestHeaderValuesType):
			funcs[i] = GetHttpRequestHeaderValuesRV
		case t.AssignableTo(httptyp.HttpRequestPathPtrType):
			funcs[i] = BuildGetHttpRequestPathFunc(f.Name)
		case t.AssignableTo(httptyp.HttpRequestQueryPtrType):
			funcs[i] = BuildGetHttpRequestQueryFunc(f.Name)
		case t.AssignableTo(httptyp.HttpRequestHeaderPtrType):
			funcs[i] = BuildGetHttpRequestHeaderFunc(f.Name)
		default:
			funcs[i] = GetNullValue
		}
	}
	return func(trv *TypeReflectValue) (reflect.Value, error) {
		structValue := reflect.New(typ).Elem()
		for i := 0; i < len(funcs); i++ {
			rv, err := funcs[i](trv)
			if err != nil {
				return httptyp.NullValue, err
			}
			structValue.Field(i).Set(rv)
		}
		return structValue, nil
	}
}

func BuildParams(typ reflect.Type) func(*TypeReflectValue) ([]reflect.Value, error) {
	n := typ.NumIn()
	funcs := make([]func(*TypeReflectValue) (reflect.Value, error), n)
	for i := 0; i < n; i++ {
		t := typ.In(i)
		switch {
		case t.String() == "http.ResponseWriter":
			funcs[i] = GetHttpResponseWriterRV
		case httptyp.HttpRequestType.ConvertibleTo(t):
			funcs[i] = GetHttpRequestRV
		case t.AssignableTo(httptyp.HttpRequestPathValuesType):
			funcs[i] = GetHttpRequestPathValuesRV
		case t.AssignableTo(httptyp.HttpRequestQueryValuesType):
			funcs[i] = GetHttpRequestQueryValuesRV
		case t.AssignableTo(httptyp.HttpRequestHeaderValuesType):
			funcs[i] = GetHttpRequestHeaderValuesRV
		case t.AssignableTo(httptyp.HttpRequestPathPtrType):
			funcs[i] = BuildGetHttpRequestPathFunc("Id")
		case t.AssignableTo(httptyp.HttpRequestQueryPtrType):
			funcs[i] = BuildGetHttpRequestQueryFunc("Id")
		case t.AssignableTo(httptyp.HttpRequestHeaderPtrType):
			funcs[i] = BuildGetHttpRequestHeaderFunc("Id")
		case t.Kind() == reflect.Struct:
			funcs[i] = BuildStructFunc(t)
		default:
			funcs[i] = BuildGetHttpRequestBodyFunc(t)
		}
	}
	return func(trv *TypeReflectValue) ([]reflect.Value, error) {
		rvs := make([]reflect.Value, len(funcs))
		for i := 0; i < len(funcs); i++ {
			rv, err := funcs[i](trv)
			if err != nil {
				return nil, err
			}
			rvs[i] = rv
		}
		return rvs, nil
	}
}
