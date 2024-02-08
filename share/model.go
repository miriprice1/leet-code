package model

type parameter struct {
	Name		string			 `json:"name"`
	Value		interface{}		 `json:"value"`
}

type TestCase struct {
	Input       []parameter      `json:"input"`
    Output      interface{}      `json:"output"`
	Length		uint			 `json:"length"`//number of inputs parameters.
}

type Question struct {
	ID          string            `json:"id" bson:"_id,omitempty"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Level       uint              `json:"level"`
	TestCases    []TestCase        `json:"testCase"`
}
