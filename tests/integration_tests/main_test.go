package integration_tests

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/cucumber/godog"
	"github.com/cucumber/messages-go/v10"
)

const (
	delay   = 10 * time.Second
	address = "http://manager:50051"
)

type updateConfigTest struct {
	resp               *httptest.ResponseRecorder
	respBody           string
	responseStatusCode int
	newConf            string
}

func TestMain(m *testing.M) {
	log.Printf("wait %s for service availability...", delay)
	time.Sleep(delay)

	status := godog.RunWithOptions("integration", func(s *godog.Suite) {
		FeatureContext(s)
	}, godog.Options{
		Format:    "pretty",
		Paths:     []string{"features"},
		Randomize: 0,
	})

	if st := m.Run(); st > status {
		status = st
	}
	os.Exit(status)
}

func (test *updateConfigTest) theResponseCodeShouldBe(code int) error {
	if test.responseStatusCode != code {
		return fmt.Errorf("unexpected status code: %d != %d", test.responseStatusCode, code)
	}

	return nil
}

func (test *updateConfigTest) iSendRequestToWithData(method string, path string, body *messages.PickleStepArgument_PickleDocString) error {
	replacer := strings.NewReplacer("\n", "", "\t", "")
	cleanJson := replacer.Replace(body.Content)
	test.newConf = cleanJson

	req, err := http.NewRequest(method, address+path, bytes.NewReader([]byte(cleanJson)))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	test.responseStatusCode = resp.StatusCode

	if method == "POST" {
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		test.respBody = string(respBody)
	}

	return nil
}

func (test *updateConfigTest) iReceiveData(data *messages.PickleStepArgument_PickleDocString) error {
	cleanJson := strings.Replace(strings.Replace(data.Content, "\n", "", -1), " ", "", -1)
	if test.respBody != cleanJson {
		return fmt.Errorf("unexpected response data: %s != %s", test.respBody, cleanJson)
	}
	return nil
}

func (test *updateConfigTest) iSendRequestTo(addr string) error {
	resp, err := http.Get(addr)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	test.responseStatusCode = resp.StatusCode

	return nil
}

func FeatureContext(s *godog.Suite) {
	test := new(updateConfigTest)

	s.Step(`^I send (POST|PUT) request to "([^"]*)" with data:$`, test.iSendRequestToWithData)
	s.Step(`^The response code should be (\d+)$`, test.theResponseCodeShouldBe)
	s.Step(`^I receive data:$`, test.iReceiveData)
	s.Step(`^I send request to "([^"]*)"$`, test.iSendRequestTo)
}
