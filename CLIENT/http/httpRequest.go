package httpRequest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"leet-code/client/display"
	"leet-code/client/helper"
	"leet-code/share"

	"github.com/charmbracelet/huh"
)

func EditQ() {
	pageNumber := 1
	id := "-1"
	for id == "-1" {//Pagination
		resData := httpGetRequest(fmt.Sprintf("http://localhost:8080/questions?page=%v",pageNumber))
		questions := *helper.ConvertUint8ToQuestions(resData)
		pageNumber++

		id = display.ChooseQuestion(questions)
	}
	myurl := fmt.Sprintf("http://localhost:8080/questions/%s", id)

	selectedData := httpGetRequest(myurl)
	selectedQuestion := *helper.ConvertUint8ToQuestion(selectedData)

	//display the question to the user as json :(
	//I will change it 😔
	jsonData, err := json.Marshal(selectedQuestion)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	obj := string(jsonData)

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewText().
				Title("Edit this question:").
				Value(&obj),
		),
	)

	err = form.Run()
	if err != nil {
		fmt.Println("Error:", err)
	}

	url := fmt.Sprintf("http://localhost:8080/questions/%s", id)
	buffer := bytes.NewBufferString(obj)
	httpRequest("PUT",url,buffer)
}

func DeleteQ() {
	pageNumber := 1
	id := "-1"
	for id == "-1" {
		resData := httpGetRequest(fmt.Sprintf("http://localhost:8080/questions?page=%v",pageNumber))
		questions := *helper.ConvertUint8ToQuestions(resData)
		pageNumber++

		id = display.ChooseQuestion(questions)
	}
	url := fmt.Sprintf("http://localhost:8080/questions/%s", id)
	httpRequest("DELETE",url,bytes.NewBuffer(nil))

}
//Doing http requests with just print response
func httpRequest(methodType string , url string , buffer *bytes.Buffer){
	req, err := http.NewRequest(methodType, url, buffer)
	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
	}

	// Send the request using the default client
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error sending HTTP request:", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
	}

	fmt.Println("Response from server:", string(body))
}

func AddQ() {
	var q module.Question

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Insert title:").
				Prompt("*").
				Value(&q.Title),
			huh.NewInput().Title("Insert description:").Prompt("*").Value(&q.Description),
			huh.NewSelect[uint]().Title("Select level").Options(
				huh.NewOption("easy", uint(1)),
				huh.NewOption("medium", uint(2)),
				huh.NewOption("hard", uint(3)),
			).Value(&q.Level),
		),
	)
	err := form.Run()
	if err != nil {
		log.Fatal(err)
	}

	q.TestCases = CaptureTestCases()

	//Convert to json for safe forwarding
	jsonData, err := json.Marshal(q)
	if err != nil {
		fmt.Println("Error marshaling question:", err)
		return
	}
	url := "http://localhost:8080/questions"
	httpRequest("POST",url,bytes.NewBuffer(jsonData))

}

func CaptureTestCases()[]module.TestCase{
	var testCases []module.TestCase
	var len uint
	println("enter intputs length:")
	fmt.Scan(&len)
	addAnother := true
	isFirstTime := true
	for addAnother {
		var tc module.TestCase

		var params []module.Parameter

		

		tc.Length = len

		for i := 0; i < int(len); i++ {
			var param module.Parameter
			var val string
			if isFirstTime{
				fmt.Println("Enter input name:")
				fmt.Scan(&param.Name)
			}
			fmt.Println("Enter input value:")
			fmt.Scan(&val)
			param.Value = helper.ParseInputToInterface(val)

			// To verify that all inputs name are the same
			if !isFirstTime{
				param.Name = testCases[0].Input[i].Name
			}

			params = append(params, param)
		}

		tc.Input = params

		var val string
		println("Enter output:")
		fmt.Scan(&val)
		tc.Output = helper.ParseInputToInterface(val)

		testCases = append(testCases, tc)

		//Option to add more test case
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewConfirm().
					Title("Are you want to add test case?").
					Affirmative("Yes!").
					Negative("No.").
					Value(&addAnother),
			),
		)
		err := form.Run()
		if err != nil {
			log.Fatal(err)
		}
		isFirstTime = false
	}
	return testCases
}

func SolveQ() {
	pageNumber := 1
	id := "-1"
	//Pagination
	for id == "-1" {
		resData := httpGetRequest(fmt.Sprintf("http://localhost:8080/questions?page=%v",pageNumber))
		questions := *helper.ConvertUint8ToQuestions(resData)
		pageNumber++

		id = display.ChooseQuestion(questions)
	}

	language := display.ChooseLanguage()
	myurl := fmt.Sprintf("http://localhost:8080/questions/%s", id)

	selectedData := httpGetRequest(myurl)
	selectedQuestion := *helper.ConvertUint8ToQuestion(selectedData)

	display.PrintQuestionDetails(selectedQuestion)

	code := display.DisplayAnswerInterface(language, selectedQuestion.TestCases[0])

	data := url.Values{}
	data.Set("code", code)
	data.Set("language", language)
	data.Set("question", string(selectedData))

	//Send Post request
	//The reason is defined separately (and not in the HttpRequest function) It is because of the type of the request body
	resp, err := http.PostForm("http://localhost:8080/runtest", data)
	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
	}

	fmt.Println("Response from server:", string(body))

}

func httpGetRequest(url string) []uint8 {
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
