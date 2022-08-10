package rest

import (
  "context"
  "net"
  "net/http"
  "net/url"
  "os"
  "os/signal"
  "strings"
  "syscall"

  "github.com/gorilla/mux"

  "github.com/illublank/go-common/config"
  "github.com/illublank/go-common/log"
  "github.com/illublank/ore/app"
  "github.com/illublank/ore/app/rest/handler"
  "github.com/illublank/ore/app/rest/listener"
)

// Restapp todo
type Restapp struct {
  app.App
  Logger   log.Logger
  Handlers *listener.ListenerHandlers
  Router   *mux.Router
  Server   *http.Server
  Address  string
  doneChan chan error
  ctx      context.Context
}

// New todo
func New(config config.Config) *Restapp {
  logger := log.NewCommonLogger("app")

  doneChan := make(chan error)

  ctx := context.Background()

  addr := config.GetString("WEBAPP_LISTEN_ADDRESS", ":8080")
  router := mux.NewRouter()
  // router.StrictSlash(true)
  // file, _ := exec.LookPath(os.Args[0])
  // path, _ := filepath.Abs(file)
  // currPath := filepath.Dir(path)
  // logger.DebugF("currPath:%v", currPath)
  // router.Handle("/", &handler.DebugHandler{Logger: logger, OrginalHandler: FileResource(currPath + "/static/index.html")})
  // router.Handle("/favicon.ico", &handler.DebugHandler{Logger: logger, OrginalHandler: FileResource(currPath + "/static/favicon.ico")})
  // router.Handle("/static/{_dummy:.*}", &handler.DebugHandler{Logger: logger, OrginalHandler: http.StripPrefix("/static/", http.FileServer(http.Dir(currPath+"/static/")))})

  return &Restapp{
    Logger:   logger,
    Handlers: listener.NewListenerHandlers(),
    Router:   router,
    Server: &http.Server{
      Addr:    addr,
      Handler: router,
      BaseContext: func(l net.Listener) context.Context {
        return ctx
      },
    },
    Address: addr,

    doneChan: doneChan,
    ctx:      ctx,
  }
}

// Handle todo
func (s *Restapp) Handle(p string, h http.Handler) *Restapp {
  s.Router.Name(p).Path(p).Handler(h)
  return s
}

// HandleFunc todo
func (s *Restapp) HandleFunc(p string, f func(http.ResponseWriter, *http.Request)) *Restapp {
  s.Router.Name(p).Path(p).HandlerFunc(f)
  return s
}

// HandleController todo
func (s *Restapp) HandleController(c Controller) *Restapp {
  for k, v := range c.GetRouteMap() {
    s.Logger.Infof("registed request path: {%v} {%v}", k, v)
    s.Handle(k, &handler.DebugHandler{Logger: s.Logger, OrginalHandler: v})
  }
  return s
}

// RegisterOnShutdown todo
func (s *Restapp) RegisterOnShutdown(f func()) {
  s.Server.RegisterOnShutdown(f)
}

// Run todo
func (s *Restapp) Run(level log.Level) error {
  s.Logger.Infof("start with address: {%v}", s.Address)
  s.Logger.SetLevel(level)

  stopChan := make(chan os.Signal, 1)
  signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

  go func() {
    sig := <-stopChan
    s.Logger.Debugf("received sig:{%v}", sig)
    // Shutdown会让监听断开，即协程里的server.ListenAndServe()将往后执行。
    // Shutdown按协议说的是graceful，Close是immediately（强杀）。
    // s.server.Close()
    if err := s.Server.Shutdown(s.ctx); err == nil || err == http.ErrServerClosed {
      s.Logger.Debug("Shutdown ok")
      s.doneChan <- nil
    } else {
      s.Logger.Errorf("Shutdown err: {%v}", err)
      s.doneChan <- err
    }
  }()

  // http 端口监听
  go func() {
    lnw, err := listener.NewDefaultListenerWrapper(s.Address, s.Handlers)
    if err != nil {
      s.Logger.Errorf("Listen and serve close err: {%v}", err)
      return
    }
    if err := s.Server.Serve(lnw); err == nil || err == http.ErrServerClosed {
      s.Logger.Debug("Listen and serve close ok")
    } else {
      s.Logger.Errorf("Listen and serve close err: {%v}", err)
    }
  }()

  // 可以考虑其他端口监听
  //

  // 退出信号
  err := <-s.doneChan
  if err != nil {
    s.Logger.Errorf("server shutdown err: {%v}", err)
  } else {
    s.Logger.Infof("server shutdown graceful")
  }
  return err
}

// SimpleRun todo
func (s *Restapp) SimpleRun() error {
  return s.Run(log.Info)
}

// StaticResource todo
func StaticResource(prefixs []string, h http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    for _, prefix := range prefixs {
      if p := strings.TrimPrefix(r.URL.Path, prefix); len(p) < len(r.URL.Path) {
        r2 := new(http.Request)
        *r2 = *r
        r2.URL = new(url.URL)
        *r2.URL = *r.URL
        r2.URL.Path = p
        h.ServeHTTP(w, r2)
      }
    }
    http.NotFound(w, r)
  })
}
