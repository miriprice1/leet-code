package main 

import "github.com/gin-gonic/gin"

func main() {

    router := gin.Default()

    // Set up routes using functions from the 'service' package
    router.GET("/questions", GetAllQuestions)
    router.GET("/questions/:id", GetQuestionByID)
    router.POST("/questions", AddQuestion)
    router.PUT("/questions/:id", UpdateQuestion)
    router.DELETE("/questions/:id", DeleteQuestion)

    router.POST("/runtest", runTest )

    router.Run(":8080")
}
