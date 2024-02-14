package helper

import (
	"context"
	"fmt"
	"leet-code/server/structures"
	"leet-code/share"
	"strconv"
	"strings"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
)

var idCounter int
var idCounterMutex sync.Mutex

func GenerateID() string {
	idCounterMutex.Lock()
	defer idCounterMutex.Unlock()

	for {
		idCounter++
		generatedID := strconv.Itoa(idCounter)
		count, _ := structures.Collection.CountDocuments(context.Background(), bson.M{"_id": generatedID})
		if count == 0 {
			return generatedID
		}
	}
}

func TypeInfo(variable interface{}) (string, bool) {

	t := fmt.Sprintf("%T", variable)
	isArray := false

	if strings.Contains(t, "[") {
		isArray = true
	}
	if strings.Contains(t, "int") {
		return "int", isArray
	}
	if strings.Contains(t, "float") {
		return "float", isArray
	}
	return t, isArray

}

func GenerateArgsSlice(testCase module.TestCase)[]string{
	var args []string
	for _, arg := range testCase.Input {
		_, isArr := TypeInfo(arg.Value)
		if isArr {
			arg.Value = addCommas(arg.Value)
		}
		args = append(args, fmt.Sprintf("%v", arg.Value))
	}
	_, isArr := TypeInfo(testCase.Output)
	if isArr {
		testCase.Output = addCommas(testCase.Output)
	}
	args = append(args, fmt.Sprintf("%v", testCase.Output))

	return args

}

func addCommas(array interface{})string{
	valueWithCommas := fmt.Sprintf("%v", array)
	valueWithCommas = strings.ReplaceAll(valueWithCommas, " ", ",")
	return valueWithCommas
}
