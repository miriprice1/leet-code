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
	code = updateCode(code, language, question.TestCases[0])
	//fmt.Println("code  =", code)
	createScriptFile(code, language)
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

func updateCode(code string, language string, tc model.TestCase) string {

	var param string
	switch language {
	case "python":
		for i := 0; i < int(tc.Length); i++ {
			tN, isArr := typeName(tc.Input[i].Value)
			if isArr { //change!!
				param = fmt.Sprintf("%v\n\t%v = @(sys.argv[%v]) if len(sys.argv) > %v else None", param, tc.Input[i].Name, i+1, i+1) //tc.Input[i].Value
			} else {
				param = fmt.Sprintf("%v\n\t%v = %v(sys.argv[%v]) if len(sys.argv) > %v else None", param, tc.Input[i].Name, tN, i+1, i+1) //tc.Input[i].Value
			}
		}
		tN, isArr := typeName(tc.Output)
		if isArr {
			param = fmt.Sprintf("%v\n\toutput = %v(sys.argv[%v]) if len(sys.argv) > %v else None", param, tN, tc.Length+1, tc.Length+1)
		} else {
			param = fmt.Sprintf("%v\n\toutput = %v(sys.argv[%v]) if len(sys.argv) > %v else None", param, tN, tc.Length+1, tc.Length+1)
		}
		code = fmt.Sprintf("import sys\n%v\nif __name__ == \"__main__\":%v\n\tif solution(%v) == output:\n\t\texit(0)\n\telse:\n\t\texit(1)", code, param, inputsStringNames(tc))
	case "js":
		for i := 0; i < int(tc.Length); i++ {
			tN, isArr := typeName(tc.Input[i].Value)
			tN = strings.ToUpper(string(tN[0])) + tN[1:]
			if isArr { //change!!
				//
				if tN == "Int" || tN == "Float" {
					param = fmt.Sprintf("%v\nconst %v = parse%v(process.argv[%v])", param, tc.Input[i].Name, tN, i+2) //,tc.Input[i].Value
				} else {
					param = fmt.Sprintf("%v\nconst %v = process.argv[%v]", param, tc.Input[i].Name, i+2) //X
				}
			} else {
				if tN == "Int" || tN == "Float" {
					param = fmt.Sprintf("%v\nconst %v = parse%v(process.argv[%v])", param, tc.Input[i].Name, tN, i+2) //V
				} else {
					param = fmt.Sprintf("%v\nconst %v = process.argv[%v]", param, tc.Input[i].Name, i+2) //V
				}
			}
		}
		tN, isArr := typeName(tc.Output)
		tN = strings.ToUpper(string(tN[0])) + tN[1:]
		if isArr { //change!!
			//
			if tN == "Int" || tN == "Float" {
				param = fmt.Sprintf("%v\nconst output = parse%v(process.argv[%v])", param, tN,tc.Length+2) //V
			} else {
				param = fmt.Sprintf("%v\nconst output = parse%v(process.argv[%v])", param, tN,tc.Length+2) //V
			}
		} else {
			if tN == "Int" || tN == "Float" {
				param = fmt.Sprintf("%v\nconst output = parse%v(process.argv[%v])", param, tN,tc.Length+2) //V
			} else {
				param = fmt.Sprintf("%v\nconst output = process.argv[%v]", param, tc.Length+2) //V
			}
		}
		code = fmt.Sprintf("%v\n%v\nif (solution(%v) == output) {\n\tprocess.exit(0);\n} else {\n\tprocess.exit(1);\n}", code, param, inputsStringNames(tc))
	}
	return code
}

func inputsStringNames(ex model.TestCase) string {
	params := ex.Input[0].Name

	if ex.Length > 1 {
		for i := 1; i < int(ex.Length); i++ {
			params = fmt.Sprintf("%v,%v", params, ex.Input[i].Name)
		}
	}
	return params
}

func createScriptFile(code string, language string) {
	var filePath string
	switch language {
	case "python":
		filePath = "../temp/script.py"
	case "js":
		filePath = "../temp/script.js"
	}

	// Write content to the file
	err := ioutil.WriteFile(filePath, []byte(code), 0644)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	fmt.Println("Content successfully written to", filePath)

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

func typeName(variable interface{}) (string, bool) {

	t := fmt.Sprintf("%T", variable)
	isArray := strings.Contains(t, "[")

	if strings.Contains(t, "int") {
		return "int", isArray
	}
	if strings.Contains(t, "float") {
		return "float", isArray
	}
	return t, isArray

}
