package main

import (
	"leet-code/client/display"
	"leet-code/client/http"
)

func main() {
	exit := false
	for !exit {
		task := display.WhatToDo()
		switch task {
		case "solve":
			httpRequest.SolveQ()
		case "add":
			httpRequest.AddQ()
		case "edit":
			httpRequest.EditQ()
		case "delete":
			httpRequest.DeleteQ()
		case "exit":
			println("Bey....")
			exit = true
		}
	}
}

