package model

type Question struct {
	Id			string	`json:"id"`
	Title       string	`json:"title"`
	Description string	`json:"description"`
	Level       uint	`json:"level"`
	TestCase    map[string]string	`json:"testCase"`
}

