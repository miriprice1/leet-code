package main

import (
	"bytes"
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
	Param1     interface{}
	Param2     interface{}
	Output     interface{}
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

type PythonTemplateData struct {
	Inputs      []Input
	Code        template.HTML
	OutputIndex int
	OutputType  string
}

type Input struct {
	Name string
	Type string
}

func updatePythonCode(code string, testCase model.TestCase) {

	var inputs []Input
	for _, tc := range testCase.Input {
		typ := typeName(tc.Value)
		if typ != "int" && typ != "float" && typ != "bool"{
			typ = ""
		}
		inputs = append(inputs, Input{Name: tc.Name, Type: typ})
	}

	typ := typeName(testCase.Output)
	if typ != "int" && typ != "float" && typ != "bool"{
		typ = ""
	}

	data := PythonTemplateData{
		Inputs:      inputs,
		Code:        template.HTML(code),
		OutputIndex: int(testCase.Length) + 1,
		OutputType:  typ,
	}

	tmpl, err := template.ParseFiles("../templates/python_template.tmpl")
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

func updateJsCode(code string, testCase model.TestCase) {

	var inputs []Input
	for _, tc := range testCase.Input {
		typ := typeName(tc.Value)
		if typ != "int" && typ != "float" && typ != "bool"{
			typ = ""
		} else {
			typ = strings.ToUpper(string(typ[0])) + typ[1:]
		}
		inputs = append(inputs, Input{Name: tc.Name, Type: typ})
	}

	typ := typeName(testCase.Output)
	if typ != "int" && typ != "float" && typ != "bool"{
		typ = ""
	}else{
		typ = strings.ToUpper(string(typ[0])) + typ[1:]
	}

	data := PythonTemplateData{
		Inputs:      inputs,
		Code:        template.HTML(code),
		OutputIndex: int(testCase.Length),
		OutputType:  typ,
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
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

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

	var questions []model.Question
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

func runTest(c *gin.Context) {

	code := c.PostForm("code")
	language := c.PostForm("language")
	questionJSON := c.PostForm("question")

	// Parse the question JSON into a struct
	var question model.Question
	if err := json.Unmarshal([]byte(questionJSON), &question); err != nil {
		println(err.Error())
	}
	println("*******************************************************************************")
	if language == "python"{
		updatePythonCode(code, question.TestCases[0])
	}else{
		updateJsCode(code, question.TestCases[0])
	}
	//code = updateCode(code, language, question.TestCases[0])
	//fmt.Println("code  =", code)
	//createScriptFile(code, language)
	createDokerfile(language)
	buildDockerImage(code)
	time.Sleep(10 * time.Second)

	isSuccecc := runJobOnK8s(language)

	// Send the variables as a response
	c.JSON(http.StatusOK, gin.H{
		"code": isSuccecc,
	})
}

//**********************************************************
//**********************************************************
//**********************************************************
//**********************************************************
//**********************************************************

func runJobOnK8s(language string) int {
	var scriptMap = map[string]string{
		"js":     "script.js",
		"python": "script.py",
	}
	script := scriptMap[language]

	if language == "js" {
		language = "node"
	}

	//for other question must fix this map.
	params := Parameters{
		Language:   language,
		ScriptFile: script,
		Param1:     5,
		Param2:     20,
		Output:     25,
	}

	yamlTemplate, err := ioutil.ReadFile("../templates/job.yaml")
	if err != nil {
		fmt.Println("Error reading template file:", err)
		os.Exit(1)
	}

	tmpl, err := template.New("job").Parse(string(yamlTemplate))
	if err != nil {
		fmt.Println("Error parsing template:", err)
		os.Exit(1)
	}
	var filledYAMLFile bytes.Buffer

	err = tmpl.Execute(&filledYAMLFile, params)
	if err != nil {
		fmt.Println("Error filling in template:", err)
		os.Exit(1)
	}

	_err := ioutil.WriteFile("job.yaml", filledYAMLFile.Bytes(), 0644)
	if _err != nil {
		fmt.Println("Error creating temporary file:", err)
		os.Exit(1)
	}
	// tmpfile, err := ioutil.TempFile("", "job-*.yaml")
	// if err != nil {
	// 	fmt.Println("Error creating temporary file:", err)
	// 	os.Exit(1)
	// }
	// defer os.Remove(tmpfile.Name()) // Clean up temporary file

	// if _, err := tmpfile.Write(filledYAMLFile.Bytes()); err != nil {
	// 	fmt.Println("Error writing to temporary file:", err)
	// 	os.Exit(1)
	// }

	cmd := exec.Command("kubectl", "apply", "-f", "job.yaml")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("Error applying YAML:", err)
		os.Exit(1)
	}

	exit := 0
	cmd = exec.Command("kubectl", "wait", "job/function-test-job", "--for=condition=complete", "--timeout=30s")
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error waiting for job to complete:", err)
		exit = 1
	}
	// time.Sleep(7 * time.Second)
	// cmd = exec.Command("kubectl", "logs", "-l", "job-name=function-test-job")
	// logs, err := cmd.Output()
	// if err != nil {
	// 	fmt.Println("Error getting logs:", err)
	// 	os.Exit(1)
	// }

	// job, err := clientset.BatchV1().Jobs(namespace).Get(context.Background(), "your-job-name", metav1.GetOptions{})
	// if err != nil {
	// 	panic(err.Error())
	// }

	// if job.Status.Succeeded > 0 {
	// 	fmt.Println("Job completed successfully")
	// } else {
	// 	fmt.Println("Job is still running or failed")
	// }

	cmd = exec.Command("kubectl", "delete", "job", "function-test-job")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("Error getting logs:", err)
		os.Exit(1)
	}

	// Check if the logs are empty
	// if len(strings.TrimSpace(string(logs))) == 0 {
	// 	fmt.Println("Logs are empty, job completed successfully.")
	// 	return 0
	// } else {
	// 	fmt.Println("Logs are not empty, job may have failed.")
	// 	return 1
	// }
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

	fmt.Println("Content successfully written to", filePath)
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
	fmt.Println("Docker image build successfuly.")
}

func typeName(variable interface{}) string {

	t := fmt.Sprintf("%T", variable)

	if strings.Contains(t, "int") {
		return "int"
	}
	if strings.Contains(t, "float") {
		return "float"
	}
	return t

}
