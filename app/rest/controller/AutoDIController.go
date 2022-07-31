package controller

import (
  "encoding/json"
  "io/ioutil"
  "net/http"
  "reflect"

  "github.com/illublank/salver/app/rest"
  "github.com/illublank/salver/app/rest/route"
)

type AutoDIController struct {
  rest.Controller
  routeMap route.RouteMap
}

func NewAutoDIController(path string, repo interface{}) *AutoDIController {
  // typ := reflect.TypeOf(repo)
  obj := reflect.ValueOf(repo)
  typ := obj.Type()
  routeMap := make(route.RouteMap)
  for i := 0; i < typ.NumMethod(); i++ {
    mTyp := typ.Method(i)
    mVal := obj.Method(i)

    switch mTyp.Name {
    case "Create", "Insert":
      routeMap.Add(&route.RouteItem{Path: path, Method: "POST", HandleFunc: buildHandleFunc(mVal)})
    case "Modify", "Update":
      routeMap.Add(&route.RouteItem{Path: path, Method: "PUT", HandleFunc: buildHandleFunc(mVal)})
    case "Remove", "Delete":
      routeMap.Add(&route.RouteItem{Path: path, Method: "DELETE", HandleFunc: buildHandleFunc(mVal)})
    case "Get":
      routeMap.Add(&route.RouteItem{Path: path, Method: "GET", HandleFunc: buildHandleFunc(mVal)})
    case "List":
      routeMap.Add(&route.RouteItem{Path: path, Method: "PROPFIND", HandleFunc: buildHandleFunc(mVal)})
    }
  }

  return &AutoDIController{
    routeMap: routeMap,
  }
}

var handleFuncTyp = reflect.TypeOf(func(w http.ResponseWriter, r *http.Request) {})

func buildHandleFunc(obj reflect.Value) func(w http.ResponseWriter, r *http.Request) {
  typ := obj.Type()
  if typ.ConvertibleTo(handleFuncTyp) {
    return obj.Interface().(func(w http.ResponseWriter, r *http.Request))
  } else {
    return func(w http.ResponseWriter, r *http.Request) {
      wVal := reflect.ValueOf(w)
      wTyp := wVal.Type()
      rVal := reflect.ValueOf(r)
      rTyp := rVal.Type()

      paramValues := make([]reflect.Value, typ.NumIn())
      for i := 0; i < typ.NumIn(); i++ {
        pTyp := typ.In(i)
        if pTyp.ConvertibleTo(wTyp) {
          paramValues[i] = wVal
        } else if pTyp.ConvertibleTo(rTyp) {
          paramValues[i] = rVal
        } else {
          bs, err := ioutil.ReadAll(r.Body)
          if err != nil {
            w.WriteHeader(http.StatusBadRequest)
            w.Write([]byte(err.Error()))
            return
          }
          newObj := reflect.New(pTyp)
          err = json.Unmarshal(bs, newObj.Interface())
          if err != nil {
            w.WriteHeader(http.StatusBadRequest)
            w.Write([]byte(err.Error()))
            return
          }
          paramValues[i] = newObj.Elem()
        }
      }
      retVals := obj.Call(paramValues)

      resObjIdx := -1
      for i := 0; i < len(retVals); i++ {
        if err, ok := retVals[i].Interface().(error); ok && err != nil {
          w.WriteHeader(http.StatusInternalServerError)
          w.Write([]byte(err.Error()))
          return
        } else if resObjIdx == -1 {
          resObjIdx = i
        }
      }
      if resObjIdx != -1 {
        resObjVal := retVals[resObjIdx]
        bs, _ := json.Marshal(resObjVal.Interface())
        w.Write(bs)
      }
    }
  }
}
