package main

import (
	"context"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"leet-code/share"
)

var client *mongo.Client
var collection *mongo.Collection
var idCounter int
var idCounterMutex sync.Mutex

func init() {
	clientOptions := options.Client().ApplyURI("mongodb://mongo:mongo@localhost:27017/questions?authSource=admin&authMechanism=SCRAM-SHA-256")
	client, _ = mongo.Connect(context.Background(), clientOptions)
	err := client.Ping(context.Background(), nil)
	if err != nil {
		panic(err)
	}
	collection = client.Database("leet-code").Collection("questions")
}

func generateID() string {
    idCounterMutex.Lock()
    defer idCounterMutex.Unlock()

    for {
        idCounter++
        generatedID := strconv.Itoa(idCounter)
        count, _ := collection.CountDocuments(context.Background(), bson.M{"_id": generatedID})
        if count == 0 {
            return generatedID
        }
    }
}

// GetAllQuestions retrieves all questions from the database.
func GetAllQuestions(c *gin.Context) {
	cursor, err := collection.Find(context.Background(), bson.D{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(context.Background())

	var questions []model.Question
	if err := cursor.All(context.Background(), &questions); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, questions)
}

// GetQuestionByID retrieves a question by its ID from the database.
func GetQuestionByID(c *gin.Context) {
	id := c.Param("id")

	var question model.Question
	err := collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&question)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Question not found"})
		return
	}

	c.JSON(http.StatusOK, question)
}

// AddQuestion adds a new question to the database.
func AddQuestion(c *gin.Context) {
	var question model.Question
	if err := c.ShouldBindJSON(&question); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	question.ID = generateID()

	_, err := collection.InsertOne(context.Background(), question)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, question)
}

// UpdateQuestion updates a question in the database by its ID.
func UpdateQuestion(c *gin.Context) {
	id := c.Param("id")

	var updatedQuestion model.Question
	if err := c.ShouldBindJSON(&updatedQuestion); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	filter := bson.M{"_id": id}
	update := bson.M{"$set": updatedQuestion}

	_, err := collection.UpdateOne(context.Background(), filter, update)
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
	_, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Question deleted successfully"})
}
