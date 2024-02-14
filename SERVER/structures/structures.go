package structures

import (
	"html/template"

	"go.mongodb.org/mongo-driver/mongo"
)

type TemplateData struct {
	Inputs         []Input
	Code           template.HTML
	OutputIndex    int
	OutputType     string
	OutputIsArray  bool
}

type Input struct {
	Name     string
	Type     string
	IsArray  bool
}

type YamlParameters struct {
	Language   string
	ScriptFile string
	Args       []string
}

var Client *mongo.Client
var Collection *mongo.Collection