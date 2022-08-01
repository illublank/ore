package rest

import "github.com/illublank/ore/app/rest/route"

// Controller todo
type Controller interface {
	GetRouteMap() route.RouteMap
}
