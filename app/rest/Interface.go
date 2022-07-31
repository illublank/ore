package rest

import "github.com/illublank/salver/app/rest/route"

// Controller todo
type Controller interface {
  GetRouteMap() route.RouteMap
}
