package service

import (
	"errors"
	"leet-code/model"
	"net/http"
	"github.com/gin-gonic/gin"
)

var questions = []model.Question{
	{
		Id:			"1",
		Title:       "Question 1",
		Description: "Description for question 1",
		Level:       1,
		TestCase: map[string]string{
			"input":  "input_value_1",
			"output": "output_value_1",
		},
	},
	{
		Id:			"2",
		Title:       "Question 2",
		Description: "Description for question 2",
		Level:       2,
		TestCase: map[string]string{
			"input":  "input_value_2",
			"output": "output_value_2",
		},
	},
	// Add more questions as needed
}


func GetQuestions(context *gin.Context){
	context.IndentedJSON(http.StatusOK,questions)

}

func AddQuestion(context *gin.Context){
	var newQuestion model.Question

	if err :=context.BindJSON(&newQuestion); err != nil{
		return
	}
	questions = append(questions, newQuestion)

	context.IndentedJSON(http.StatusCreated , newQuestion)
}

func getQuestionById(id string) (*model.Question, error){
	for i, q := range(questions) {
		if q.Id == id{
			return &questions[i], nil
		}
	}
	return nil, errors.New("question not found")
}

func DeleteQuestion(context *gin.Context){
	id := context.Param("id")
	isPassed := false
	for i, q := range(questions){
		if q.Id == id {
			questions = append(questions[:i], questions[i+1:]...)
			isPassed = true
			break
		}
	}
	if !isPassed {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "question id not found"})
	} else {
		context.IndentedJSON(http.StatusOK, gin.H{"message": "question deleted successfully"})
	}
	
}

func UpdateQuestion(context *gin.Context){
	DeleteQuestion(context)
	AddQuestion(context)
}