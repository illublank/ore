package controller

import (
  "encoding/json"
  "io"
  "net/http"
  "reflect"

  "github.com/illublank/ore/app/rest/httptyp"
)

var HandleFuncType = reflect.TypeOf((*http.HandlerFunc)(nil)).Elem()
var ResponseType = reflect.TypeOf(http.ResponseWriter(nil))
var RequestType = reflect.TypeOf((*http.Request)(nil))

var NullValue = reflect.ValueOf(nil)

type TypeReflectValue struct {
  InType                 reflect.Type
  OBJ_HttpResponseWriter http.ResponseWriter
  RV_HttpResponseWriter  reflect.Value
  OBJ_HttpRequest        *http.Request
  RV_HttpRequest         reflect.Value
  OBJ_HttpPathValues     httptyp.PathVals
  RV_HttpPathValues      reflect.Value
  OBJ_HttpQueryValues    httptyp.QueryVals
  RV_HttpQueryValues     reflect.Value
  OBJ_HttpHeaderValues   httptyp.HeaderVals
  RV_HttpHeaderValues    reflect.Value
}

type TypeTuple struct {
  Con func(reflect.Type) bool
  Val func(*TypeReflectValue) reflect.Value
}

func Lookup(t reflect.Type, trv *TypeReflectValue) (reflect.Value, error) {
  switch {
  case t.String() == "http.ResponseWriter":
    return trv.RV_HttpResponseWriter, nil
  case RequestType.ConvertibleTo(t):
    return trv.RV_HttpRequest, nil
  case t.AssignableTo(httptyp.PathValsType):
    return trv.RV_HttpPathValues, nil
  case t.AssignableTo(httptyp.RequestPathPtrType):
    v, e := trv.OBJ_HttpPathValues.Get("Id")
    if e {
      return reflect.ValueOf(httptyp.ParseRequestPath(v)), nil
    } else {
      return NullValue, nil
    }
  case t.AssignableTo(httptyp.QueryValsType):
    return trv.RV_HttpQueryValues, nil
  case t.AssignableTo(httptyp.HeaderValsType):
    return trv.RV_HttpHeaderValues, nil
  case t.Kind() == reflect.Struct:
    pTyp := trv.InType
    structValue := reflect.New(pTyp).Elem()
    for i := 0; i < pTyp.NumField(); i++ {
      f := pTyp.Field(i)
      if trv.RV_HttpResponseWriter.CanConvert(f.Type) {
        structValue.Field(i).Set(trv.RV_HttpResponseWriter)
      } else if trv.RV_HttpRequest.CanConvert(f.Type) {
        structValue.Field(i).Set(trv.RV_HttpRequest)
      } else if f.Type.AssignableTo(httptyp.RequestPathPtrType) {
        v, e := trv.OBJ_HttpPathValues.Get(f.Name)
        if e {
          structValue.Field(i).Set(reflect.ValueOf(httptyp.ParseRequestPath(v)))
        } else {
          structValue.Field(i).Set(NullValue)
        }
      } else if f.Type.AssignableTo(httptyp.RequestQueryPtrType) {
        v, e := trv.OBJ_HttpQueryValues.Get(f.Name)
        if e {
          structValue.Field(i).Set(reflect.ValueOf(httptyp.ParseRequestQuery(v)))
        } else {
          structValue.Field(i).Set(NullValue)
        }
      } else if f.Type.AssignableTo(httptyp.RequestHeaderPtrType) {
        v, e := trv.OBJ_HttpHeaderValues.Get(f.Name)
        if e {
          structValue.Field(i).Set(reflect.ValueOf(httptyp.ParseRequestHeader(v)))
        } else {
          structValue.Field(i).Set(NullValue)
        }
      } else {
        structValue.Field(i).Set(NullValue)
      }
    }
    return structValue, nil
  default:
    bs, err := io.ReadAll(trv.OBJ_HttpRequest.Body)
    if err != nil {
      return NullValue, err
    }
    newObj := reflect.New(t)
    err = json.Unmarshal(bs, newObj.Interface())
    if err != nil {
      return NullValue, err
    }
    return newObj.Elem(), nil
  }
}
