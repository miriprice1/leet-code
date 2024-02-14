package module

type Parameter struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}

type TestCase struct {
	Input  []Parameter `json:"input"`
	Output interface{} `json:"output"`
	Length uint        `json:"length"` //number of inputs parameters.
}

type Question struct {
	ID          string     `json:"id" bson:"_id,omitempty"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Level       uint       `json:"level"`
	TestCases   []TestCase `json:"testCase"`
}

const PageSize  = 10