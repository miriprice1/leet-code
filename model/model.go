package model

type Question struct {
    ID          string            `json:"id" bson:"_id,omitempty"`
    Title       string            `json:"title"`
    Description string            `json:"description"`
    Level       uint              `json:"level"`
    TestCase    map[string]string `json:"testCase"`
}
