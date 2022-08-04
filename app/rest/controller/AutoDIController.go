package controller

import (
  "encoding/json"
  "fmt"
  "net/http"
  "reflect"

  "github.com/gorilla/mux"
  "github.com/illublank/ore/app/rest"
  "github.com/illublank/ore/app/rest/httptyp"
  "github.com/illublank/ore/app/rest/route"
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
  idPath := "/" + path
  for i := 0; i < typ.NumMethod(); i++ {
    mTyp := typ.Method(i)
    mVal := obj.Method(i)

    switch mTyp.Name {
    case "Create", "Insert":
      routeMap.Add(&route.RouteItem{Path: idPath, Method: "POST", HandleFunc: buildHandleFunc(mVal)})
    case "Modify", "Update":
      routeMap.Add(&route.RouteItem{Path: idPath, Method: "PUT", HandleFunc: buildHandleFunc(mVal)})
    case "Remove", "Delete":
      routeMap.Add(&route.RouteItem{Path: idPath + "/{Id}", Method: "DELETE", HandleFunc: buildHandleFunc(mVal)})
    case "Get":
      routeMap.Add(&route.RouteItem{Path: idPath + "/{Id}", Method: "GET", HandleFunc: buildHandleFunc(mVal)})
    case "List":
      routeMap.Add(&route.RouteItem{Path: idPath, Method: "PROPFIND", HandleFunc: buildHandleFunc(mVal)})
    }
  }

  return &AutoDIController{
    routeMap: routeMap,
  }
}

func (s *AutoDIController) GetRouteMap() route.RouteMap {
  return s.routeMap
}

func buildHandleFunc(obj reflect.Value) http.HandlerFunc {
  typ := obj.Type()
  if typ.ConvertibleTo(HandleFuncType) {
    return obj.Interface().(http.HandlerFunc)
  } else {
    return http.HandlerFunc(
      func(w http.ResponseWriter, r *http.Request) {

        wVal := reflect.ValueOf(w)
        // wTyp := wVal.Type()
        // fmt.Println(wTyp)
        rVal := reflect.ValueOf(r)
        // rTyp := rVal.Type()
        pathVals := mux.Vars(r)
        pathValsRV := reflect.ValueOf(httptyp.PathVals(pathVals))
        queryVals := r.URL.Query()
        queryValsRV := reflect.ValueOf(httptyp.QueryVals(queryVals))
        headerVals := r.Header
        headerValsRV := reflect.ValueOf(httptyp.HeaderVals(headerVals))
        trv := &TypeReflectValue{
          OBJ_HttpResponseWriter: w,
          RV_HttpResponseWriter:  wVal,
          OBJ_HttpRequest:        r,
          RV_HttpRequest:         rVal,
          OBJ_HttpPathValues:     httptyp.PathVals(pathVals),
          RV_HttpPathValues:      pathValsRV,
          OBJ_HttpQueryValues:    httptyp.QueryVals(queryVals),
          RV_HttpQueryValues:     queryValsRV,
          OBJ_HttpHeaderValues:   httptyp.HeaderVals(headerVals),
          RV_HttpHeaderValues:    headerValsRV,
        }
        paramValues := make([]reflect.Value, typ.NumIn())
        for i := 0; i < typ.NumIn(); i++ {
          pTyp := typ.In(i)
          fmt.Println(pTyp, pTyp.Kind())
          trv.InType = pTyp
          rv, err := tMap(pTyp, trv)
          if err != nil {
            w.WriteHeader(http.StatusBadRequest)
            w.Write([]byte(err.Error()))
          }
          paramValues[i] = rv
        }
        fmt.Println(paramValues)
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
      },
    )
  }
}
