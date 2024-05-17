package handler

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/dev-blax/dalali-plus/internal/model"
	"github.com/dev-blax/dalali-plus/internal/service"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client
var messagesCollection *mongo.Collection

func InitMongoDB() {
	var err error
	mongoClient, err = mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}

	messagesCollection = mongoClient.Database("chatapp").Collection("messages")
}

func ChatHandler(wsService *service.WebSocketService) fiber.Handler {
	return websocket.New(func(c *websocket.Conn) {
		wsService.HandleConnections(c)
	})
}

func GetMessages(c *fiber.Ctx) error {
	var messages []model.Message
	cursor, err := messagesCollection.Find(context.TODO(), bson.D{})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	if err = cursor.All(context.TODO(), &messages); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.JSON(messages)
}

func SaveMessage(c *fiber.Ctx, wsService *service.WebSocketService) {
	var message model.Message
	if err := c.BodyParser(&message); err != nil {
		c.Status(fiber.StatusBadRequest).SendString(err.Error())
		return
	}

	message.Time = time.Now().Unix()

	_, err := messagesCollection.InsertOne(context.TODO(), message)
	if err != nil {
		c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		return
	}

	msgBytes, err := json.Marshal(message)
	if err != nil {
		c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		return
	}

	wsService.Broadcast <- msgBytes
	c.SendStatus(fiber.StatusOK)
}
