package main

import (
	"log"
	"time"

	"github.com/asynched/sse-chat-demo/http/controllers"
	"github.com/gofiber/fiber/v2"
)

func init() {
	log.SetPrefix("[SERVER] ")
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile | log.Lmsgprefix)
}

func main() {
	app := fiber.New(fiber.Config{DisableStartupMessage: true})

	log.Println("Initializing server")
	chatController := controllers.NewChatController()

	log.Println("Setting up routes")
	app.Get("/chat", chatController.SSE)
	app.Post("/chat", chatController.CreateMessage)

	log.Println("Serving static files")
	app.Static("/", "./static")

	go func() {
		log.Println("Server has started on port :8080")
		log.Println("Open in browser: http://localhost:8080")

		for {
			time.Sleep(10 * time.Second)
			log.Printf("Number of clients: %d\n", app.Server().GetOpenConnectionsCount())
		}
	}()

	log.Fatalf("Server crashed: %v\n", app.Listen(":8080"))
}
