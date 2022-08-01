package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	idPath := "/" + path + "/{id}"
	for i := 0; i < typ.NumMethod(); i++ {
		mTyp := typ.Method(i)
		mVal := obj.Method(i)

		switch mTyp.Name {
		case "Create", "Insert":
			routeMap.Add(&route.RouteItem{Path: idPath, Method: "POST", HandleFunc: buildHandleFunc(mVal)})
		case "Modify", "Update":
			routeMap.Add(&route.RouteItem{Path: idPath, Method: "PUT", HandleFunc: buildHandleFunc(mVal)})
		case "Remove", "Delete":
			routeMap.Add(&route.RouteItem{Path: idPath, Method: "DELETE", HandleFunc: buildHandleFunc(mVal)})
		case "Get":
			routeMap.Add(&route.RouteItem{Path: idPath, Method: "GET", HandleFunc: buildHandleFunc(mVal)})
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

var HandleFuncTyp = reflect.TypeOf((*http.HandlerFunc)(nil)).Elem()

var NullValue = reflect.ValueOf(nil)

func buildHandleFunc(obj reflect.Value) http.HandlerFunc {
	typ := obj.Type()
	if typ.ConvertibleTo(HandleFuncTyp) {
		return obj.Interface().(http.HandlerFunc)
	} else {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				wVal := reflect.ValueOf(w)
				wTyp := wVal.Type()
				rVal := reflect.ValueOf(r)
				rTyp := rVal.Type()
				pathVals := mux.Vars(r)
				queryVals := r.URL.Query()
				headerVals := r.Header
				paramValues := make([]reflect.Value, typ.NumIn())
				for i := 0; i < typ.NumIn(); i++ {
					pTyp := typ.In(i)
					isPath := pTyp.AssignableTo(httptyp.RequestPathPtrType)
					isQuery := pTyp.AssignableTo(httptyp.RequestQueryPtrType)
					isHeader := pTyp.AssignableTo(httptyp.RequestHeaderPtrType)
					fmt.Println(isPath, isQuery, isHeader)
					if pTyp.ConvertibleTo(wTyp) {
						paramValues[i] = wVal
					} else if pTyp.ConvertibleTo(rTyp) {
						paramValues[i] = rVal
					} else if pTyp.AssignableTo(httptyp.RequestPathPtrType) {
						if len(pathVals) > 0 {
							for _, v := range pathVals {
								paramValues[i] = reflect.ValueOf(httptyp.ParseRequestPath(v))
								continue
							}
						} else {
							paramValues[i] = NullValue
						}
					} else if pTyp.AssignableTo(httptyp.RequestQueryPtrType) {
						if len(queryVals) > 0 {
							for _, v := range queryVals {
								paramValues[i] = reflect.ValueOf(v)
								continue
							}
						} else {
							paramValues[i] = NullValue
						}
					} else if pTyp.AssignableTo(httptyp.RequestHeaderPtrType) {
						if len(headerVals) > 0 {
							for _, v := range headerVals {
								paramValues[i] = reflect.ValueOf(v)
								continue
							}
						} else {
							paramValues[i] = NullValue
						}
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
			},
		)
	}
}
