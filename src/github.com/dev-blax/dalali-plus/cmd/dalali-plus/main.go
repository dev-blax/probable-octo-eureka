package main

import (
	"github.com/dev-blax/dalali-plus/internal/handler"
	"github.com/dev-blax/dalali-plus/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New()
	app.Use(cors.New())
	handler.InitMongoDB()
	wsService := service.NewWebSocketService()

	go wsService.HandleMessages()

	app.Use("/ws", handler.ChatHandler(wsService))
	app.Get("/", handler.Home)
	app.Get("/messages", handler.GetMessages)
	app.Post("/messages", func(c *fiber.Ctx) error {
		handler.SaveMessage(c, wsService)
		return nil
	})

	app.Listen(":4000")
}
