package controllers

import (
	"bufio"
	"fmt"
	"log"

	"github.com/asynched/sse-chat-demo/sync/channels"
	"github.com/gofiber/fiber/v2"
)

type ChatController struct {
	broadcaster *channels.Broadcaster[string]
}

func (controller *ChatController) SSE(c *fiber.Ctx) error {
	log.Printf("New client connected: %s\n", c.IP())

	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")

	c.Context().Response.SetBodyStreamWriter(func(w *bufio.Writer) {
		messages := controller.broadcaster.Subscribe()
		defer controller.broadcaster.Unsubscribe(messages)

		w.WriteString("\n\n")
		w.Flush()

		for {
			message := <-messages

			if _, err := w.WriteString(fmt.Sprintf("data: %s\n\n", message)); err != nil {
				log.Printf("Client disconnected, error: '%v'\n", err)
				break
			}

			if err := w.Flush(); err != nil {
				log.Printf("Client disconnected, error: '%v'\n", err)
				break
			}
		}
	})

	return nil
}

func (controller *ChatController) CreateMessage(c *fiber.Ctx) error {
	var message struct {
		Content string `json:"content"`
	}

	if err := c.BodyParser(&message); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	if message.Content == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	controller.broadcaster.Broadcast(message.Content)

	return nil
}

func NewChatController() *ChatController {
	return &ChatController{
		broadcaster: channels.NewBroadcaster[string](),
	}
}
