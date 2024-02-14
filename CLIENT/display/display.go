package display

import (
	"fmt"
	"log"

	"github.com/charmbracelet/huh"
	"leet-code/client/helper"
	"leet-code/share"
)

func ChooseLanguage() string {

	var language string
	form := huh.NewForm(
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
	return language
}

func PrintQuestionDetails(q module.Question) {
	fmt.Printf("Title: %v \n", q.Title)
	fmt.Printf("Desription: %v \n", q.Description)
	fmt.Printf("Level: %v \n", q.Level)

	tc := q.TestCases[0] //the first example
	input := fmt.Sprintf("\"%v\" = %v", tc.Input[0].Name, tc.Input[0].Value)

	if tc.Length > 1 {
		for i := 1; i < int(tc.Length); i++ {
			input = fmt.Sprintf("%v , \"%v\" = %v ", input, tc.Input[i].Name, tc.Input[i].Value)
		}
	}

	fmt.Printf("\bExample:\n\t-intput: %v\n\t-output: %v\n", input, tc.Output)

}

func ChooseQuestion(questions []module.Question) string {

	var options []huh.Option[string]

	// Loop through questions and add titles as options
	for _, q := range questions {
		options = append(options, huh.Option[string]{
			Key:   q.Title,
			Value: q.ID,
		})
	}
	//Pagination
	if len(options) == module.PageSize {
		options = append(options, huh.Option[string]{
			Key:   "-more-",
			Value: "-1",
		})
	}

	var id string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Choose your question").
				Options(options...).
				Value(&id),
		),
	)

	err := form.Run()
	if err != nil {
		log.Fatal(err)
	}
	return id
}

func DisplayAnswerInterface(language string, testCase module.TestCase) string {

	var code string
	params := helper.InputsStringNames(testCase)

	switch language {
	case "python":
		code = fmt.Sprintf("def solution(%v):\n\t", params)
	case "js":
		code = fmt.Sprintf("function solution(%v){\n\t\n};", params)
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

func WhatToDo() string {
	var id string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("What do you want to do?").
				Options(
					huh.NewOption("Solve question", "solve"),
					huh.NewOption("Add question", "add"),
					huh.NewOption("Edit question", "edit"),
					huh.NewOption("Delete question", "delete"),
					huh.NewOption("Exit", "exit"),
				).
				Value(&id),
		),
	)
	err := form.Run()
	if err != nil {
		log.Fatal(err)
	}
	return id
}
