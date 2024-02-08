package main

import (
	"fmt"

	"github.com/charmbracelet/huh"
)

func main() {

	data :=httpGetRequest("http://localhost:8080/questions")
	questions :=*convertUint8ToQuestions(data)

	var options [] huh.Option[string]

	// Loop through questions and add titles as options
	for _, q := range questions {
		options = append(options, huh.Option[string]{
			Key: q.Title,
			Value: q.ID, 
		})
	}
	
	id,language := chooseQuestionAndLanguage(options)

	url := fmt.Sprintf("http://localhost:8080/questions/%s",id)

	selectedData := httpGetRequest(url)
	selectedQuestion := *convertUint8ToQuestion(selectedData)

	printQuestionDetails(selectedQuestion)

	code := displayAnswerInterface(language,selectedQuestion.TestCases[0])

	//add execute command for any test case.
	for _,tCase := range(selectedQuestion.TestCases){
		code = fmt.Sprintf("%v\nsolution(%v)",code,inputsStringValues(tCase))
	} 
	
	fmt.Println("code  =",code)

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
