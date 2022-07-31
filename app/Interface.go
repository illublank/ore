package app

import "github.com/illublank/go-common/log"

// App todo
type App interface {
  Run(log.Level) error
  SimpleRun() error
}
