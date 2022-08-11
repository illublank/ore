package rest_test

import (
	"fmt"
	"go/types"
	"net"
	"net/http"
	"testing"

	"github.com/illublank/go-common/config/mock"
	"github.com/illublank/go-common/log"
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

func (t *TestComponent) Get(obj *Object, id *httptyp.HttpRequestPath, r *http.Request, w http.ResponseWriter, header httptyp.HttpRequestHeaderValues, query httptyp.HttpRequestQueryValues, pathVals httptyp.HttpRequestPathValues, p struct {
	Id *httptyp.HttpRequestPath
}) string {
	x, _ := header.Get("Test-Key")
	h, _ := query.Get("ha-ha")
	h2, _ := query.Get("s.a")
	fmt.Println(id, x, h, h2, p.Id, obj)
	return id.String()
}

func TestRestapp(t *testing.T) {

	cfg := mock.NewMapConfig(collection.NewGoMap())

	app := rest.New(cfg)

	app.HandleController(controller.NewAutoDIController("test", &TestComponent{}))

	app.Handlers.BeforeAccept.Add(func(net.Listener) {

		// for i := 0; i < 5; i++ {
		//   fmt.Println("BeforeAccept", i)
		//   time.Sleep(time.Second)
		// }
	})
	remoteCount := 0
	app.Handlers.AfterAccept.Add(func(l net.Listener, c net.Conn, e error) {
		remoteCount++
		fmt.Println("remote", remoteCount, c.RemoteAddr())
	})
	// app.Handlers.BeforeAccept.Add(func(net.Listener) {

	//   for i := 0; i < 10; i++ {
	//     fmt.Println("BeforeAccept", i)
	//     time.Sleep(time.Second)
	//   }
	// })

	app.Run(log.Debug)
}

func T(s string) types.Type {
	return nil
}
