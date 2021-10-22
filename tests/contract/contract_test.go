package contract

import (
	"flag"
	"github.com/cucumber/godog"
	"log"
	"os"
	"testing"
)

var path string
var tags string

func TestMain(m *testing.M) {
	flag.StringVar(&path, "path", "", "Path to feature files")
	flag.StringVar(&tags, "tags", "@feature=access", "Tags filter")
	flag.Parse()
	if path == "" {
		flag.Usage()
		log.Fatal("Please set the feature files path")
	}
	os.Exit(m.Run())
}

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format: "junit:Cucumber.xml",
			Paths:  []string{path},
			Tags:   tags,
			TestingT: t, // Testing instance that will run subtests.
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}
