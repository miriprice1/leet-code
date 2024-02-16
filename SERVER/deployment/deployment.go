package deployment

import (
	"fmt"
	"leet-code/server/files"
	"leet-code/server/helper"
	"leet-code/server/structures"
	module "leet-code/share"
	"os"
	"os/exec"
)

func BuildAndRunJob(language string, testCase module.TestCase) bool {

	var scriptNameMap = map[string]string{
		"js":     "script.js",
		"python": "script.py",
	}
	//Choose the correct commands for running
	scriptName := scriptNameMap[language]

	if language == "js" {
		language = "node"
	}

	args := helper.GenerateArgsSlice(testCase)

	params := structures.YamlParameters{
		Language:   language,
		ScriptFile: scriptName,
		Args:       args,
	}

	files.CreateYamlFile(params)

	runningResult := runJobOnK8s()

	return runningResult
}

func runJobOnK8s() bool {

	cmd := exec.Command("kubectl", "apply", "-f", "../temp/job.yaml")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("Error applying YAML:", err)
		os.Exit(1)
	}

	exit := true
	cmd = exec.Command("kubectl", "wait", "job/function-test-job", "--for=condition=complete", "--timeout=20s")
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error waiting for job to complete:", err)
		exit = false

	}

	cmd = exec.Command("kubectl", "delete", "job", "function-test-job")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("Error getting logs:", err)
		os.Exit(1)
	}
	return exit
}

func BuildDockerImage(code string) {
	cmd := exec.Command("sh", "-c", "cd ./temp && docker build -t test .")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		fmt.Println("Error building docker image", err)
		return
	}
	defer fmt.Println("Docker image build successfuly.")
}
