package runner

import (
	"database/sql"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/joho/godotenv"

	"github.com/nofuture17/gonkey/checker"
	"github.com/nofuture17/gonkey/checker/response_body"
	"github.com/nofuture17/gonkey/checker/response_db"
	"github.com/nofuture17/gonkey/checker/response_header"
	"github.com/nofuture17/gonkey/fixtures"
	"github.com/nofuture17/gonkey/mocks"
	"github.com/nofuture17/gonkey/output"
	"github.com/nofuture17/gonkey/output/allure_report"
	testingOutput "github.com/nofuture17/gonkey/output/testing"
	"github.com/nofuture17/gonkey/testloader/yaml_file"
	"github.com/nofuture17/gonkey/variables"
)

type RunWithTestingParams struct {
	Server      *httptest.Server
	TestsDir    string
	Mocks       *mocks.Mocks
	FixturesDir string
	DB          *sql.DB
	DbType      fixtures.DbType
	EnvFilePath string
	OutputFunc  output.OutputInterface
	Checkers    []checker.CheckerInterface
}

// RunWithTesting is a helper function the wraps the common Run and provides simple way
// to configure Gonkey by filling the params structure.
func RunWithTesting(t *testing.T, params *RunWithTestingParams) {
	var mocksLoader *mocks.Loader
	if params.Mocks != nil {
		mocksLoader = mocks.NewLoader(params.Mocks)
	}

	if params.EnvFilePath != "" {
		if err := godotenv.Load(params.EnvFilePath); err != nil {
			t.Fatal(err)
		}
	}

	debug := os.Getenv("GONKEY_DEBUG") != ""

	var fixturesLoader fixtures.Loader
	if params.DB != nil {
		fixturesLoader = fixtures.NewLoader(&fixtures.Config{
			Location: params.FixturesDir,
			DB:       params.DB,
			Debug:    debug,
			DbType:   params.DbType,
		})
	}

	yamlLoader := yaml_file.NewLoader(params.TestsDir)
	yamlLoader.SetFileFilter(os.Getenv("GONKEY_FILE_FILTER"))

	r := New(
		&Config{
			Host:           params.Server.URL,
			Mocks:          params.Mocks,
			MocksLoader:    mocksLoader,
			FixturesLoader: fixturesLoader,
			Variables:      variables.New(),
		},
		yamlLoader,
	)

	if params.OutputFunc != nil {
		r.AddOutput(params.OutputFunc)
	} else {
		r.AddOutput(testingOutput.NewOutput(t))
	}

	if os.Getenv("GONKEY_ALLURE_DIR") != "" {
		allureOutput := allure_report.NewOutput("Gonkey", os.Getenv("GONKEY_ALLURE_DIR"))
		defer allureOutput.Finalize()
		r.AddOutput(allureOutput)
	}

	r.AddCheckers(response_body.NewChecker())
	r.AddCheckers(response_header.NewChecker())

	if params.DB != nil {
		r.AddCheckers(response_db.NewChecker(params.DB))
	}

	r.AddCheckers(params.Checkers...)

	_, err := r.Run()
	if err != nil {
		t.Fatal(err)
	}
}
