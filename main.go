package main

import (
    "github.com/gin-gonic/gin"
    "leet-code/service"
)

func main() {
    router := gin.Default()

    // Set up routes using functions from the 'service' package
    router.GET("/questions", service.GetAllQuestions)
    router.GET("/questions/:id", service.GetQuestionByID)
    router.POST("/questions", service.AddQuestion)
    router.PUT("/questions/:id", service.UpdateQuestion)
    router.DELETE("/questions/:id", service.DeleteQuestion)

    router.Run(":8080")
}
