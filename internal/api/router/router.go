package router

import "github.com/gofiber/fiber/v2"

type RouteFunc func(router fiber.Router)

type RouteGroup struct {
	Router fiber.Router
}

func newRouteGroup(router fiber.Router) *RouteGroup {
	return &RouteGroup{
		Router: router,
	}
}

func (rg *RouteGroup) apply(routeFuncs ...RouteFunc) {
	for _, rf := range routeFuncs {
		rf(rg.Router)
	}
}

func InitRouter(router fiber.Router) {
	rg := newRouteGroup(router)
	rg.apply(
		AuthRouter,
		MeRouter,
		CatalogRouter,
	)
}
