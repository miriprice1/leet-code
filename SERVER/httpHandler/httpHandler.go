package httphandler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"leet-code/server/deployment"
	"leet-code/server/files"
	"leet-code/server/helper"
	"leet-code/server/structures"
	"leet-code/share"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client
var Collection *mongo.Collection

func GetAllQuestions(c *gin.Context) {

	// Parse page number and page size from query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", fmt.Sprintf("%v", module.PageSize)))

	// Calculate skip value based on page number and page size
	skip := (page - 1) * pageSize

	// Define options for pagination
	findOptions := options.Find()
	findOptions.SetLimit(int64(pageSize))
	findOptions.SetSkip(int64(skip))

	// Query MongoDB for questions with pagination
	cursor, err := structures.Collection.Find(context.Background(), bson.D{}, findOptions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(context.Background())

	var questions []module.Question
	//this is indices that all the data return,
	if err := cursor.All(context.Background(), &questions); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return paginated questions
	c.JSON(http.StatusOK, questions)
}

// GetQuestionByID retrieves a question by its ID from the database.
func GetQuestionByID(c *gin.Context) {
	id := c.Param("id")

	var question module.Question
	err := structures.Collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&question)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Question not found"})
		return
	}

	c.JSON(http.StatusOK, question)
}

// AddQuestion adds a new question to the database.
func AddQuestion(c *gin.Context) {
	var question module.Question
	if err := c.ShouldBindJSON(&question); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	question.ID = helper.GenerateID()

	_, err := structures.Collection.InsertOne(context.Background(), question)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"massage": "data inserted successfuly"})
}

// UpdateQuestion updates a question in the database by its ID.
func UpdateQuestion(c *gin.Context) {
	id := c.Param("id")

	var updatedQuestion module.Question
	if err := c.ShouldBindJSON(&updatedQuestion); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	filter := bson.M{"_id": id}
	update := bson.M{"$set": updatedQuestion}

	_, err := structures.Collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Question updated successfully"})
}

// DeleteQuestion deletes a question from the database by its ID.
func DeleteQuestion(c *gin.Context) {
	id := c.Param("id")

	filter := bson.M{"_id": id}
	_, err := structures.Collection.DeleteOne(context.Background(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Question deleted successfully"})
}

func RunTest(c *gin.Context) {

	code := c.PostForm("code")
	language := c.PostForm("language")
	questionJSON := c.PostForm("question")

	// Parse the question JSON into a struct
	var question module.Question
	if err := json.Unmarshal([]byte(questionJSON), &question); err != nil {
		println(err.Error())
	}
	if language == "python" {
		files.CreatePythonCode(code, question.TestCases[0])
	} else {
		files.CreateJsCode(code, question.TestCases[0])
	}
	files.CreateDokerfile(language)
	deployment.BuildDockerImage(code)
	var isSuccecc bool
	for _, tc := range question.TestCases {
		isSuccecc = deployment.BuildAndRunJob(language, tc)
		if !isSuccecc {
			break
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": isSuccecc,
	})
}
