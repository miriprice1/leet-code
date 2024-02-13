package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

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

//try to combining this functions
func inputsStringNames(ex model.TestCase)string{
	params := ex.Input[0].Name

	if ex.Length>1{
		for i := 1; i < int(ex.Length); i++ {
			params = fmt.Sprintf("%v,%v",params,ex.Input[i].Name)
		}
	}
	return params
}

// func inputsStringValues(ex model.TestCase)interface{}{
// 	params := ex.Input[0].Value

// 	if ex.Length>1{
// 		for i := 1; i < int(ex.Length); i++ {
// 			params = fmt.Sprintf("%v,%v",params,ex.Input[i].Value)
// 		}
// 	}
// 	return params
// }



func displayAnswerInterface(language string, testCase model.TestCase)string{
	
	var code string
	params := inputsStringNames(testCase)

	switch language {
	case "python":
		code = fmt.Sprintf("def solution(%v):\n\t",params)//If there is time left - add a field of the function name and replace it with a "solution"
	case "js":
		code = fmt.Sprintf("function solution(%v){\n\t\n};",params)
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

func chooseQuestionAndLanguage(options []huh.Option[string])(string,string){
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

func printQuestionDetails(q model.Question){
	fmt.Printf("Title: %v \n",q.Title)
	fmt.Printf("Desription: %v \n",q.Description)
	fmt.Printf("Level: %v \n",q.Level)

	tc := q.TestCases[0] //the first example
	input := fmt.Sprintf("\"%v\" = %v",tc.Input[0].Name,tc.Input[0].Value)

	if tc.Length>1{
		for i := 1; i < int(tc.Length); i++ {
			input = fmt.Sprintf("%v , \"%v\" = %v ",input,tc.Input[i].Name,tc.Input[i].Value)
		}
	}

	fmt.Printf("\bExample:\n\t-intput: %v\n\t-output: %v\n",input,tc.Output)

}