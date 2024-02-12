package main

import (
	"fmt"
	"net/url"
	"net/http"

	"github.com/charmbracelet/huh"
)

func main() {

	resData := httpGetRequest("http://localhost:8080/questions")
	questions := *convertUint8ToQuestions(resData)

	var options []huh.Option[string]

	// Loop through questions and add titles as options
	for _, q := range questions {
		options = append(options, huh.Option[string]{
			Key:   q.Title,
			Value: q.ID,
		})
	}

	id, language := chooseQuestionAndLanguage(options)//on pagination I will need to split this function

	myurl := fmt.Sprintf("http://localhost:8080/questions/%s", id)

	selectedData := httpGetRequest(myurl)
	selectedQuestion := *convertUint8ToQuestion(selectedData)

	printQuestionDetails(selectedQuestion)

	code := displayAnswerInterface(language, selectedQuestion.TestCases[0])

	data := url.Values{}
	data.Set("code",code)
	data.Set("language",language)
	data.Set("question", string(selectedData))

	response, err := http.PostForm("http://localhost:8080/runtest", data)//change run test to another name!
	if err != nil{
		println(err.Error())
	}
	
	println(response.Body)

	// //add execute command for any test case.
	// for _,tCase := range(selectedQuestion.TestCases){
	// 	code = fmt.Sprintf("%v\nsolution(%v)",code,inputsStringValues(tCase))
	// }

}

// func runServer(serverReady chan<- struct{}) {
// 	cmd := exec.Command("sh", "-c", "cd ../server && go run main.go service.go model.go")
// 	cmd.Stdout = os.Stdout
// 	cmd.Stderr = os.Stderr
// 	err := cmd.Start()
// 	if err != nil {
// 		fmt.Println("Error executing server:", err)
// 		close(serverReady) // Signal that the server failed to start
// 		return
// 	}
// 	fmt.Println("Server is starting...")
// 	close(serverReady) // Signal that the server has started successfully
// 	time.Sleep(15 * time.Second)
// }
// "os"
// 	"os/exec"
//	"time"
