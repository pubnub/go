package contract

import (
	"flag"
	"github.com/cucumber/godog"
	"log"
	"os"
	"testing"
)

var path string
var tagsFilter string
var format string

func TestMain(m *testing.M) {
	flag.StringVar(&path, "path", "", "Path to feature files")
	flag.StringVar(&tagsFilter, "tagsFilter", "~@skip && ~@na=go && ~@beta", "Tags filter")
	flag.StringVar(&format, "format", "pretty", "Output formatter")
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
			Format: format,
			Paths:  []string{path},
			Tags:   tagsFilter,
			TestingT: t, // Testing instance that will run subtests.
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}
