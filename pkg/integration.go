package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"time"

	"github.com/aws/aws-lambda-go/events"
)

// Run to call sam invoke on the system for intergration testing.
// The template is the generated cloudfromation template for this service.
// The handler is the lambda handler function you wish to run the itergration test for (named as in the template).
// A JSON string with requestData data to process - this can also be empty.
// stdErrRegex is needed because we only pull errors we are looking for out of stdErr which can have
// can be multiple unrelated errors present.
func Run(template, handler, requestData string, stdErrRegx string) (*events.APIGatewayProxyResponse, error) {
	var data events.APIGatewayProxyResponse
	testTime := time.Now().Format(time.RFC3339)
	tmp := fmt.Sprintf("tmp-%v", testTime)
	requestEvent := fmt.Sprintf("%v/req-%v.json", tmp, testTime)

	if err := os.Mkdir(tmp, os.ModePerm); err != nil {
		return nil, err
	}
	defer os.RemoveAll(tmp)

	if err := ioutil.WriteFile(requestEvent, []byte(requestData), 0644); err != nil {
		return nil, err
	}

	cmd := exec.Command("sam", "local", "invoke", handler, "--template", template, "--event", requestEvent)
	cmd.Env = append(os.Environ(), "LC_ALL=en_US.utf-8", "LANG=en_US.utf-8")
	var stdOut, stdErr bytes.Buffer
	cmd.Stdout = &stdOut
	cmd.Stderr = &stdErr

	if err := cmd.Run(); err != nil {
		return nil, err
	}

	regx := regexp.MustCompile(stdErrRegx)
	if len(regx.FindStringSubmatch(stdErr.String())) > 0 {
		if err := json.Unmarshal(stdErr.Bytes(), &data); err != nil {
			return nil, err
		}
		return &data, nil
	}

	if err := json.Unmarshal(stdOut.Bytes(), &data); err != nil {
		return nil, err
	}

	return &data, nil
}
