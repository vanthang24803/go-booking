package router

import "github.com/gofiber/fiber/v2"

type RouteFunc func(router fiber.Router)

type RouteGroup struct {
	Router fiber.Router
}

func NewRouteGroup(router fiber.Router) *RouteGroup {
	return &RouteGroup{
		Router: router,
	}
}

func (rg *RouteGroup) Apply(routeFuncs ...RouteFunc) {
	for _, rf := range routeFuncs {
		rf(rg.Router)
	}
}

func InitRouter(router fiber.Router) {
	rg := NewRouteGroup(router)
	rg.Apply(
		AuthRouter,
		MeRouter,
	)
}
