package httphandler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"leet-code/server/deployment"
	"leet-code/server/files"
	"leet-code/server/helper"
	"leet-code/server/structures"
	module "leet-code/share"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetAllQuestions(c *gin.Context) {

	// Parse page number and page size from query parameters or costant
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
	if err := cursor.All(context.Background(), &questions); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return paginated questions
	c.JSON(http.StatusOK, questions)
}

func GetQuestionByID(c *gin.Context) {
	id := c.Param("id")

	var question module.Question
	//MongoDB query
	err := structures.Collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&question)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Question not found"})
		return
	}

	c.JSON(http.StatusOK, question)
}

func AddQuestion(c *gin.Context) {
	var question module.Question
	if err := c.ShouldBindJSON(&question); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//Generating unique question ID
	question.ID = helper.GenerateID()

	//MongoDB query
	_, err := structures.Collection.InsertOne(context.Background(), question)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"massage": "data inserted successfuly"})
}

func UpdateQuestion(c *gin.Context) {
	id := c.Param("id")

	var updatedQuestion module.Question
	//Parsing updated question from the context
	if err := c.ShouldBindJSON(&updatedQuestion); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	filter := bson.M{"_id": id}
	update := bson.M{"$set": updatedQuestion}

	//MongoDB query
	_, err := structures.Collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Question updated successfully"})
}

func DeleteQuestion(c *gin.Context) {
	id := c.Param("id")

	filter := bson.M{"_id": id}
	//MongoDB query
	_, err := structures.Collection.DeleteOne(context.Background(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Question deleted successfully"})
}

func RunTest(c *gin.Context) {

	//Parse the data from the context
	code := c.PostForm("code")
	language := c.PostForm("language")
	questionJSON := c.PostForm("question")

	// Parse the question JSON into a struct
	var question module.Question
	if err := json.Unmarshal([]byte(questionJSON), &question); err != nil {
		println("Error:", err.Error())
	}

	//Create temporary folder
	os.Mkdir("../temp", 0755)
	//Remove the temporary folder
	defer os.RemoveAll("../temp")

	//Generate scriptFile
	if language == "python" {
		files.CreatePythonCode(code, question.TestCases[0])
	} else {
		files.CreateJsCode(code, question.TestCases[0])
	}

	files.CreateDokerfile(language)
	deployment.BuildDockerImage(code)

	//Check all the testcases
	var isSuccecc bool
	for _, tc := range question.TestCases {
		isSuccecc = deployment.BuildAndRunJob(language, tc)
		// To finish in the first that failed
		if !isSuccecc {
			break
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"Is success? ": isSuccecc,
	})

}
