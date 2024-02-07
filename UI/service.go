package main

import(
	"fmt"
	"log"
	"net/http"
	"encoding/json"
	"io/ioutil"
	
	"leet-code/share"
	"github.com/charmbracelet/huh"
)

//in the functions replace the return nil with the return the errors
func convertUint8ToQuestion(body []uint8) *model.Question{
	var questions model.Question
	err := json.Unmarshal(body, &questions)
	if err != nil {
		fmt.Println("Error unmarshaling data:", err)
		return nil
	}
	return &questions
}

func convertUint8ToQuestions(body []uint8) *[]model.Question{
	var questions []model.Question
	err := json.Unmarshal(body, &questions)
	if err != nil {
		fmt.Println("Error unmarshaling data:", err)
		return nil
	}
	return &questions
}


func httpGetRequest(url string)[]uint8{
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error sending GET request:", err)
		return nil
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil
	}
	return body
}


func displayAnswerInterface(language string)string{
	var code string

	switch language {
	case "python":
		code = "def solution(){\n\t\n}\n\nsolution()"//If there is time left - add a field of the function name and replace it with a "solution"
	case "js":
		code = "function solution(){\n\t\n}\n\nsolution();"
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewText().
				Title("Write your solution: ").
				Value(&code),
		),
	)
	err := form.Run()
	if err != nil {
		log.Fatal(err)
	}
	return code
}

func chooseQuestionAndLanguage(options [] huh.Option[string])(string,string){
	var id string
	var language string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Choose your question").
				Options(options...).
				Value(&id),
		),
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Choose programing language:").
				Options(
					huh.NewOption("Python", "python"),
					huh.NewOption("Java script", "js"),
				).
				Value(&language),
		),
	)

	err := form.Run()
	if err != nil {
		log.Fatal(err)
	}
	return id,language
}
