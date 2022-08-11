package controller

import (
	"encoding/json"
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
	paramFunc := BuildParams(typ)
	if typ.ConvertibleTo(httptyp.HttpHandleFuncType) {
		return obj.Interface().(http.HandlerFunc)
	} else {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.Header().Add("Content-Type", "application/json")
				wVal := reflect.ValueOf(w)
				rVal := reflect.ValueOf(r)
				pathVals := mux.Vars(r)
				pathValsRV := reflect.ValueOf(httptyp.HttpRequestPathValues(pathVals))
				queryVals := r.URL.Query()
				queryValsRV := reflect.ValueOf(httptyp.HttpRequestQueryValues(queryVals))
				headerVals := r.Header
				headerValsRV := reflect.ValueOf(httptyp.HttpRequestHeaderValues(headerVals))
				trv := &TypeReflectValue{
					OBJ_HttpResponseWriter:      w,
					OBJ_HttpRequest:             r,
					OBJ_HttpRequestPathValues:   httptyp.HttpRequestPathValues(pathVals),
					OBJ_HttpRequestQueryValues:  httptyp.HttpRequestQueryValues(queryVals),
					OBJ_HttpRequestHeaderValues: httptyp.HttpRequestHeaderValues(headerVals),

					RV_HttpResponseWriter:      wVal,
					RV_HttpRequest:             rVal,
					RV_HttpRequestPathValues:   pathValsRV,
					RV_HttpRequestQueryValues:  queryValsRV,
					RV_HttpRequestHeaderValues: headerValsRV,
				}
				paramValues, err := paramFunc(trv)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte(err.Error()))
					return
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
			},
		)
	}
}

type Handler func(w http.ResponseWriter, r *http.Request)

func (s Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}
