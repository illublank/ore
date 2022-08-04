package rest_test

import (
  "fmt"
  "go/types"
  "net/http"
  "testing"

  "github.com/illublank/go-common/config/mock"
  "github.com/illublank/go-common/typ/collection"
  "github.com/illublank/ore/app/rest"
  "github.com/illublank/ore/app/rest/controller"
  "github.com/illublank/ore/app/rest/httptyp"
)

type TestComponent struct {
}

type Object struct {
  A string
  B string
}

func (t *TestComponent) Get(id *httptyp.RequestPath, obj *Object, r *http.Request, w http.ResponseWriter, header httptyp.HeaderVals, query httptyp.QueryVals, pathVals httptyp.PathVals, p struct {
  Id *httptyp.RequestPath
}) string {
  x, _ := header.Get("Test-Key")
  h, _ := query.Get("ha-ha")
  h2, _ := query.Get("s.a")
  fmt.Println(id, obj, x, h, h2, p.Id)
  return id.String()
}

func TestRestapp(t *testing.T) {

  cfg := mock.NewMapConfig(collection.NewGoMap())

  app := rest.New(cfg, nil)

  app.HandleController(controller.NewAutoDIController("test", &TestComponent{}))

  app.SimpleRun()
}

func T(s string) types.Type {
  return nil
}
