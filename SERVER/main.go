package main

import (
	"context"
	httphandler "leet-code/server/httpHandler"
	"leet-code/server/structures"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {

    router := gin.Default()

    // Set up routes using functions from the 'service' package
    router.GET("/questions", httphandler.GetAllQuestions)
    router.GET("/questions/:id", httphandler.GetQuestionByID)
    router.POST("/questions", httphandler.AddQuestion)
    router.PUT("/questions/:id", httphandler.UpdateQuestion)
    router.DELETE("/questions/:id", httphandler.DeleteQuestion)

    router.POST("/runtest", httphandler.RunTest )

    router.Run(":8080")
}

func init() {
	clientOptions := options.Client().ApplyURI("mongodb://mongo:mongo@localhost:27017/questions?authSource=admin&authMechanism=SCRAM-SHA-256")
	structures.Client, _ = mongo.Connect(context.Background(), clientOptions)
	err := structures.Client.Ping(context.Background(), nil)
	if err != nil {
		panic(err)
	}
	structures.Collection = structures.Client.Database("leet-code").Collection("questions")
}