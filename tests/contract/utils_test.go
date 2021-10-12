package contract

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/cucumber/godog"
)

type expectResponse struct {
	Contract     string       `json:"contract"`
	Expectations expectations `json:"expectations"`
}

type expectations struct {
	Pending []string `json:"pending"`
	Failed  []string `json:"failed"`
}

type contractNameKey struct{}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Before(before)
	ctx.After(after)
	MapSteps(ctx)
}

const initialize_contract_url = "http://%s/init?__contract__script__=%s"
const expect_contract_url = "http://%s/expect"

func before(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
	contractTestConfig, err := newContractTestConfig()
	newCtx := context.WithValue(ctx, contractTestConfigKey{}, contractTestConfig)
	if err != nil {
		return ctx, err
	}
	commonState := newCommonState(contractTestConfig)
	accessState := newAccessState(commonState.pubNub)

	newCtx = context.WithValue(newCtx, commonStateKey{}, commonState)
	newCtx = context.WithValue(newCtx, accessStateKey{}, accessState)

	if !contractTestConfig.serverMock {
		return newCtx, nil
	}

	contractName := ""
	for _, tag := range sc.Tags {
		if strings.Contains(tag.Name, "contract") {
			contractName = strings.SplitN(tag.Name, "=", 2)[1]
		}
	}
	newCtx = context.WithValue(newCtx, contractNameKey{}, contractName)

	if len(contractName) != 0 {

		_, err := http.Get(fmt.Sprintf(initialize_contract_url, contractTestConfig.hostPort, contractName))
		if err != nil {
			return newCtx, err
		}
	}
	return newCtx, nil
}

func after(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
	contractTestConfig := ctx.Value(contractTestConfigKey{}).(contractTestConfig)
	if !contractTestConfig.serverMock {
		return ctx, nil
	}

	contractName := ctx.Value(contractNameKey{}).(string)

	if len(contractName) != 0 {
		resp, err := http.Get(fmt.Sprintf(expect_contract_url, contractTestConfig.hostPort))
		if err != nil {
			return ctx, err
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return ctx, err
		}
		bodyUnmarshalled := expectResponse{}
		json.Unmarshal(body, &bodyUnmarshalled)

		if len(bodyUnmarshalled.Expectations.Failed) != 0 ||
			len(bodyUnmarshalled.Expectations.Pending) != 0 {
			failureReason, _ := json.Marshal(bodyUnmarshalled)
			return ctx, errors.New("Failed or pending expectations: " + string(failureReason))
		}
	}

	return ctx, nil
}

func getAccessState(ctx context.Context) *accessState {
	return ctx.Value(accessStateKey{}).(*accessState)
}

func getCommonState(ctx context.Context) *commonState {
	return ctx.Value(commonStateKey{}).(*commonState)
}

type contractTestConfigKey struct{}

type contractTestConfig struct {
	publishKey   string
	subscribeKey string
	secretKey    string
	serverMock   bool
	hostPort     string
	secure       bool
}

func newContractTestConfig() (contractTestConfig, error) {
	var serverMock bool
	var secure bool
	var err error
	serverMock, err = getenvBoolWithDefault("SERVER_MOCK", true)
	if err != nil {
		return contractTestConfig{}, nil
	}

	secure, err = getenvBoolWithDefault("SECURE", false)
	return contractTestConfig{
		publishKey:   getenvWithDefault("PUBLISH_KEY", "pubKey"),
		subscribeKey: getenvWithDefault("SUBSCRIBE_KEY", "subKey"),
		secretKey:    getenvWithDefault("SECRET_KEY", "secKey"),
		hostPort:     getenvWithDefault("HOST_PORT", "localhost:8090"),
		serverMock:   serverMock,
		secure:       secure,
	}, err
}

func getenvWithDefault(name string, defaultValue string) string {
	stringValue, ok := os.LookupEnv(name)
	if ok {
		return stringValue
	} else {
		return defaultValue
	}
}

func getenvBoolWithDefault(name string, defaultValue bool) (bool, error) {
	stringValue, ok := os.LookupEnv(name)
	if !ok {
		return defaultValue, nil
	}

	value, err := strconv.ParseBool(stringValue)
	if err != nil {
		return true, err
	}

	return value, nil
}
