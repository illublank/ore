package listener

import (
	"bytes"
	"fmt"
	"net"

	"github.com/illublank/go-common/syn"
)

type ListenerHandlers struct {
  BeforeAccept *syn.WaitFuncs[func(net.Listener)]
  AfterAccept  *syn.WaitFuncs[func(net.Listener, net.Conn, error)]
  BeforeClose *syn.WaitFuncs[func(net.Listener)]
  AfterClose *syn.WaitFuncs[func(net.Listener, error)]
}

func NewListenerHandlers() *ListenerHandlers {
  return &ListenerHandlers{
    BeforeAccept: syn.NewWaitFuncs[func(net.Listener)](),
    AfterAccept: syn.NewWaitFuncs[func(net.Listener, net.Conn, error)](),
    BeforeClose: syn.NewWaitFuncs[func(net.Listener)](),
    AfterClose: syn.NewWaitFuncs[func(net.Listener, error)](),
  }
}

type ListenerWrapper struct {
  net.Listener
  Original     net.Listener
  Handlers *ListenerHandlers
}

func NewListenerWrapper(original net.Listener, handlers *ListenerHandlers) *ListenerWrapper {
  return &ListenerWrapper{
    Original: original,
    Handlers: handlers,
  }
}

func NewDefaultListenerWrapper(address string, handlers *ListenerHandlers) (*ListenerWrapper, error) {
  if address == "" {
    address = ":http"
  }
  ln, err := net.Listen("tcp", address)
  if err != nil {
    return nil, err
  }
  return NewListenerWrapper(ln, handlers), nil
}

// Accept waits for and returns the next connection to the listener.
func (s *ListenerWrapper) Accept() (net.Conn, error) {
  fmt.Println("Accept")
  s.Handlers.BeforeAccept.ForEach(func(f func(net.Listener)) {
    f(s.Original)
  })
  conn, err := s.Original.Accept()

  conn = NewTraceConn(conn)
  s.Handlers.AfterAccept.ForEach(func(f func(net.Listener,net.Conn,error)) {
    f(s.Original, conn, err)
  })
  return conn, err
}

// Close closes the listener.
// Any blocked Accept operations will be unblocked and return errors.
func (s *ListenerWrapper) Close() error {
  fmt.Println("Close")
  s.Handlers.BeforeClose.ForEach(func(f func(net.Listener)) {
    f(s.Original)
  })
  err := s.Original.Close()
  s.Handlers.AfterClose.ForEach(func(f func(net.Listener, error)) {
    f(s.Original, err)
  })
  return err
}

// Addr returns the listener's network address.
func (s *ListenerWrapper) Addr() net.Addr {
  return s.Original.Addr()
}


type TraceConn struct {
  net.Conn
  wbs *bytes.Buffer
}

func NewTraceConn(c net.Conn) *TraceConn {
  return &TraceConn{
    Conn: c,
    wbs: &bytes.Buffer{},
  }
}

func (c *TraceConn) Write(b []byte) (n int, err error) {
  c.wbs.Write(b)
  return c.Conn.Write(b)
}

func(c *TraceConn) Close() error {
  fmt.Println("TraceConn.Close", c.wbs.Len(), c.wbs.String())
  return c.Conn.Close()
}
