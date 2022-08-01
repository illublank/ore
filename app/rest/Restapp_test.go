package rest_test

import (
	"fmt"
	"go/types"
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

func (t *TestComponent) Get(id *httptyp.RequestPath, obj *Object) string {
	fmt.Println(obj)
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
