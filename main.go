package main

import (
	"github.com/gin-gonic/gin"
	"leet-code/service"
)


func main(){
	server := gin.Default()
	server.GET("/getQuestion", service.GetQuestions)
	server.POST("/addQuestion/", service.AddQuestion)
	server.DELETE("/deleteQuestion/:id",service.DeleteQuestion)
	server.PUT("updateQuestion/:id",service.UpdateQuestion)
	server.Run("localhost:8089")
}