package files

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"leet-code/server/helper"
	"leet-code/server/structures"
	module "leet-code/share"
	"os"
	"strings"
)

func CreatePythonCode(code string, testCase module.TestCase) {

	var inputs []structures.TemplateInput
	//Generate Inputs struct for sending to template
	for _, tc := range testCase.Input {
		typ, isArray := helper.TypeInfo(tc.Value)
		if typ != "int" && typ != "float" && typ != "bool" {
			typ = "string"
		}
		inputs = append(inputs, structures.TemplateInput{Name: tc.Name, Type: typ, IsArray: isArray})
	}

	typ, isArray := helper.TypeInfo(testCase.Output)
	if typ != "int" && typ != "float" && typ != "bool" {
		typ = "string"
	}

	data := structures.TemplateData{
		Inputs:        inputs,
		Code:          template.HTML(code),
		OutputIndex:   int(testCase.Length),
		OutputType:    typ,
		OutputIsArray: isArray,
	}

	//Craete script file
	tmplPath := "../templates/python_template.tmpl"
	filePath := "../temp/script.py"
	createFileFromTemplate(tmplPath, filePath, data)

}

func createFileFromTemplate(tmplPath string, outputPath string, data structures.TemplateData) {
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		panic(err)
	}

	outputFile, err := os.Create(outputPath)
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()

	//Build the file according to the parameters
	err = tmpl.Execute(outputFile, data)
	if err != nil {
		panic(err)
	}
}

func CreateJsCode(code string, testCase module.TestCase) {

	var inputs []structures.TemplateInput
	//Generate Inputs struct for sending to template
	for _, tc := range testCase.Input {
		typ, isArray := helper.TypeInfo(tc.Value)
		if typ != "int" && typ != "float" { //} && typ != "bool" {
			typ = "string"
		} else {
			typ = strings.ToUpper(string(typ[0])) + typ[1:]
		}
		inputs = append(inputs, structures.TemplateInput{Name: tc.Name, Type: typ, IsArray: isArray})
	}

	typ, isArray := helper.TypeInfo(testCase.Output)
	// if typ == "bool" {
	// 	typ = "boolean"
	// }
	if typ != "int" && typ != "float" { // && typ != "boolean" {
		typ = "string"
	} else {
		typ = strings.ToUpper(string(typ[0])) + typ[1:]
	}

	data := structures.TemplateData{
		Inputs:        inputs,
		Code:          template.HTML(code),
		OutputIndex:   int(testCase.Length),
		OutputType:    typ,
		OutputIsArray: isArray,
	}

	//Create script file
	tmplPath := "../templates/js_template.tmpl"
	filePath := "../temp/script.js"
	createFileFromTemplate(tmplPath, filePath, data)

}

func CreateDokerfile(language string) {
	filePath := "../temp/Dockerfile"
	var code string
	//Define other dockerfile for each language
	switch language {
	case "python":
		code = "FROM python:3.9\nWORKDIR /app\nCOPY script.py /app/script.py\nCMD [\"python\", \"script.py\"]"
	case "js":
		code = "FROM node:14\nWORKDIR /app\nCOPY script.js .\nCMD [\"node\", \"script.js\"]"
	}
	// Write content to the file
	err := ioutil.WriteFile(filePath, []byte(code), 0644)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	defer fmt.Println("Content successfully written to", filePath)
}

func CreateYamlFile(params structures.YamlParameters) {
	tmpl, err := template.ParseFiles("../templates/job.tmpl")
	if err != nil {
		panic(err)
	}

	outputFile, err := os.Create("../temp/job.yaml")
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()

	//Build the file according to the parameters
	err = tmpl.Execute(outputFile, params)
	if err != nil {
		panic(err)
	}
}
