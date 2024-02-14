package helper

import (
	"encoding/json"
	"fmt"

	"leet-code/share"
)

func ConvertUint8ToQuestion(body []uint8) *module.Question {
	var questions module.Question
	err := json.Unmarshal(body, &questions)
	if err != nil {
		fmt.Println("Error unmarshaling data:", err)
		return nil
	}
	return &questions
}

func ConvertUint8ToQuestions(body []uint8) *[]module.Question {
	var questions []module.Question
	err := json.Unmarshal(body, &questions)
	if err != nil {
		fmt.Println("Error unmarshaling data:", err)
		return nil
	}
	return &questions
}

func InputsStringNames(ex module.TestCase) string {
	params := ex.Input[0].Name

	if ex.Length > 1 {
		for i := 1; i < int(ex.Length); i++ {
			params = fmt.Sprintf("%v,%v", params, ex.Input[i].Name)
		}
	}
	return params
}

func ParseInput(input string) interface{} {
	// Try to parse the input string as JSON
	var value interface{}
	err := json.Unmarshal([]byte(input), &value)
	if err == nil {
		return value
	}

	// If parsing as JSON fails, treat the input as a string
	return input
}