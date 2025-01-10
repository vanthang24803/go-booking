package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/may20xx/booking/config"
	"github.com/may20xx/booking/internal/api/middleware/interceptor"
	"github.com/may20xx/booking/internal/api/router"
	"github.com/may20xx/booking/internal/database"
	"github.com/may20xx/booking/internal/utils"
	"github.com/may20xx/booking/pkg/log"
)

func init() {

	app := fiber.New()

	config := config.GetConfig()

	app.Use(interceptor.Logging())
	app.Use(interceptor.Error())

	err := database.InitDatabase()

	if err != nil {
		log.Msg.Panic(err)
	}

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(utils.NewAppError(fiber.StatusOK, "Hello World!"))
	})

	v1 := app.Group("/api/v1")

	router.InitRouter(v1)

	app.Use(interceptor.RouteNotMatch())

	log.Msg.Infof("Server is running on port %s 🚀", config.Port)

	app.Listen(":" + config.Port)

	defer database.CloseDatabase()
}
