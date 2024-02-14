package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"leet-code/share"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client
var collection *mongo.Collection
var idCounter int
var idCounterMutex sync.Mutex

type Parameters struct {
	Language   string
	ScriptFile string
	Args       []string
}

func init() {
	clientOptions := options.Client().ApplyURI("mongodb://mongo:mongo@localhost:27017/questions?authSource=admin&authMechanism=SCRAM-SHA-256")
	client, _ = mongo.Connect(context.Background(), clientOptions)
	err := client.Ping(context.Background(), nil)
	if err != nil {
		panic(err)
	}
	collection = client.Database("leet-code").Collection("questions")
}

type TemplateData struct {
	Inputs         []Input
	Code           template.HTML
	OutputIndex    int
	OutputType     string
	OutputIsArray  bool
	OutputArrayLen int
}

type Input struct {
	Name     string
	Type     string
	IsArray  bool
	ArrayLen int
}

func updatePythonCode(code string, testCase module.TestCase) {

	var inputs []Input
	var x int
	for _, tc := range testCase.Input {
		typ, isArray, len := typeInfo(tc.Value)
		if typ != "int" && typ != "float" && typ != "bool" {
			typ = ""
		}
		x += len
		inputs = append(inputs, Input{Name: tc.Name, Type: typ, IsArray: isArray, ArrayLen: len})
	}

	typ, isArray, len := typeInfo(testCase.Output)
	if typ != "int" && typ != "float" && typ != "bool" {
		typ = ""
	}

	data := TemplateData{
		Inputs:         inputs,
		Code:           template.HTML(code),
		OutputIndex:    int(testCase.Length),
		OutputType:     typ,
		OutputIsArray:  isArray,
		OutputArrayLen: len,
	}

	tmpl := template.New("try_arrays.tmpl").Funcs(template.FuncMap{"add": AddFunc})
	tmpl, err := tmpl.ParseFiles("../templates/try_arrays.tmpl") //python_template.tmpl")
	if err != nil {
		panic(err)
	}

	outputFile, err := os.Create("../temp/script.py")
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()

	err = tmpl.Execute(outputFile, data)
	if err != nil {
		panic(err)
	}
}
func AddFunc(a, b int) int {
	return a + b
}

func updateJsCode(code string, testCase module.TestCase) {

	var inputs []Input
	for _, tc := range testCase.Input {
		typ, isArray, len := typeInfo(tc.Value)
		if typ != "int" && typ != "float" && typ != "bool" {
			typ = ""
		} else {
			typ = strings.ToUpper(string(typ[0])) + typ[1:]
		}
		inputs = append(inputs, Input{Name: tc.Name, Type: typ, IsArray: isArray, ArrayLen: len})
	}

	typ, isArray, len := typeInfo(testCase.Output)
	if typ == "bool" {
		typ = "boolean"
	}
	if typ != "int" && typ != "float" && typ != "boolean" {
		typ = ""
	} else {
		typ = strings.ToUpper(string(typ[0])) + typ[1:]
	}

	data := TemplateData{
		Inputs:         inputs,
		Code:           template.HTML(code),
		OutputIndex:    int(testCase.Length),
		OutputType:     typ,
		OutputIsArray:  isArray,
		OutputArrayLen: len,
	}

	tmpl, err := template.ParseFiles("../templates/js_template.tmpl")
	if err != nil {
		panic(err)
	}

	outputFile, err := os.Create("../temp/script.js")
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()

	err = tmpl.Execute(outputFile, data)
	if err != nil {
		panic(err)
	}
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
	cursor, err := collection.Find(context.Background(), bson.D{}, findOptions)
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
	err := collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&question)
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

	question.ID = generateID()

	_, err := collection.InsertOne(context.Background(), question)
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

func runTest(c *gin.Context) {

	code := c.PostForm("code")
	language := c.PostForm("language")
	questionJSON := c.PostForm("question")

	// Parse the question JSON into a struct
	var question module.Question
	if err := json.Unmarshal([]byte(questionJSON), &question); err != nil {
		println(err.Error())
	}
	if language == "python" {
		updatePythonCode(code, question.TestCases[0])
	} else {
		updateJsCode(code, question.TestCases[0])
	}
	createDokerfile(language)
	buildDockerImage(code)
	time.Sleep(8 * time.Second)
	var isSuccecc bool //may be trigerr
	for _, tc := range question.TestCases {
		isSuccecc = runJobOnK8s(language, tc)
		if !isSuccecc {
			break
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": isSuccecc,
	})
}

//**********************************************************
//**********************************************************
//**********************************************************
//**********************************************************
//**********************************************************

func runJobOnK8s(language string, testCase module.TestCase) bool {
	var scriptMap = map[string]string{
		"js":     "script.js",
		"python": "script.py",
	}
	script := scriptMap[language]

	if language == "js" {
		language = "node"
	}
	var args []string

	for _, arg := range testCase.Input {
		_, isArr, _ := typeInfo(arg.Value)
		if isArr {
			valueWithCommas := fmt.Sprintf("%v", arg.Value)
			fmt.Println("Before replacement:", valueWithCommas)

			// Replace spaces with commas
			valueWithCommas = strings.ReplaceAll(valueWithCommas, " ", ",")
			fmt.Println("After replacement:", valueWithCommas)

			// Update testCase.Output with the modified string
			arg.Value = valueWithCommas

			fmt.Println("Updated output:", arg.Value)
		}
		args = append(args, fmt.Sprintf("%v", arg.Value))
	}
	_, isArr, _ := typeInfo(testCase.Output)
	if isArr {
		valueWithCommas := fmt.Sprintf("%v", testCase.Output)
		valueWithCommas = strings.ReplaceAll(valueWithCommas, " ", ",")
		testCase.Output = valueWithCommas
	}
	args = append(args, fmt.Sprintf("%v", testCase.Output))

	params := Parameters{
		Language:   language,
		ScriptFile: script,
		Args:       args,
	}
	/*******************/

	tmpl, err := template.ParseFiles("../templates/job.tmpl")
	if err != nil {
		panic(err)
	}

	outputFile, err := os.Create("../temp/job.yaml")
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()

	err = tmpl.Execute(outputFile, params)
	if err != nil {
		panic(err)
	}

	cmd := exec.Command("kubectl", "apply", "-f", "../temp/job.yaml")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("Error applying YAML:", err)
		os.Exit(1)
	}

	exit := true
	cmd = exec.Command("kubectl", "wait", "job/function-test-job", "--for=condition=complete", "--timeout=30s")
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error waiting for job to complete:", err)
		exit = false
	}

	cmd = exec.Command("kubectl", "delete", "job", "function-test-job")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("Error getting logs:", err)
		os.Exit(1)
	}

	//defer os.RemoveAll("../temp")

	return exit
}

func createDokerfile(language string) {
	filePath := "../temp/Dockerfile"
	var code string
	switch language {
	case "python":
		code = "FROM python:3.9\nWORKDIR /app\nCOPY script.py /app/script.py\nCMD [\"python\", \"script.py\"]"
	case "js":
		code = "FROM node:14\nWORKDIR /app\nCOPY script.js .\nCMD [\"node\", \"script.js\"]"
	}
	// Write content to the file
	err := ioutil.WriteFile(filePath, []byte(code), 0644)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	defer fmt.Println("Content successfully written to", filePath)
}

func buildDockerImage(code string) {
	cmd := exec.Command("sh", "-c", "cd ../temp && docker build -t test .")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		fmt.Println("Error building docker image", err)
		return
	}
	defer fmt.Println("Docker image build successfuly.")
}

func typeInfo(variable interface{}) (string, bool, int) {

	t := fmt.Sprintf("%T", variable)
	isArray := false
	length := 0

	if strings.Contains(t, "[") {
		isArray = true
		println(t)
		length = len(variable.([]interface{}))
	}
	if strings.Contains(t, "int") {
		return "int", isArray, length
	}
	if strings.Contains(t, "float") {
		return "float", isArray, length
	}
	return t, isArray, length

}
