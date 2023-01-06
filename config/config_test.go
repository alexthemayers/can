package config

import (
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

const RepositoryName = "gin-in-a-can"

func TestConfig_LoadConfig(t *testing.T) {
	tests := []struct {
		configFile     string
		expectedConfig Config
		expectedErr    bool
	}{
		{configFile: "", expectedConfig: Config{}, expectedErr: true},
		{configFile: "test_fixtures/example.yaml", expectedConfig: testConfig},
	}
	for i, testCase := range tests {
		os.Args = []string{"can"}
		if testCase.configFile != "" {
			os.Args = append(os.Args, "--configFile", testCase.configFile)
		}
		configData, err := LoadConfig()

		if !testCase.expectedErr && err != nil {
			t.Log("Test Case: ", i)
			t.Error("Unexpected error occurred while loading config file:", err)
		}
		if !reflect.DeepEqual(configData, testCase.expectedConfig) {
			t.Log("Test Case: ", i)
			t.Error("Loaded OpenApiConfig did not match expected config")
		}
	}
}

var (
	repoLoc    = findRepoRoot(".")
	testConfig = Config{
		Generator: GeneratorConfig{
			ModuleName:        "github.com/sasswart/gin-in-a-can",
			BasePackageName:   "api",
			TemplateDirectory: repoLoc + "/templates/go-gin",
		},
		OpenAPIFile: OpenApiConfig{
			OpenAPIFile: repoLoc + "/openapi/fixtures/testRefs/validation.yaml",
		},
		OutputPath:       repoLoc + "/test/output",
		WorkingDirectory: repoLoc + "/test/output",
		ConfigFilePath:   repoLoc + "/templates/go-gin",
	}
)

func TestConfig_FindRepoRoot(t *testing.T) {
	wd, _ := os.Getwd()
	tests := []struct {
		input    string
		expected string
	}{
		{input: "/Users/alex/code/gin-in-a-can/render/docs/openapi.yml", expected: "/Users/alex/code/gin-in-a-can"},
		{input: "/Users/alex/code/gin-in-a-can/render", expected: "/Users/alex/code/gin-in-a-can"},
		{input: "./", expected: wd},
		{input: "", expected: ""},
		{input: "github.com/sasswart/gin-in-a-can", expected: ""},
	}

	for i, testCase := range tests {
		got := findRepoRoot(testCase.input)
		if got != testCase.expected {
			t.Logf("Test case number %d\n", i+1)
			t.Logf("got %s, expected %s\n\n", got, testCase.expected)
			t.Fail()
		}
	}
}

func findRepoRoot(s string) string {
	switch true {
	case !strings.HasPrefix(s, "/"):
		return ""
	case strings.HasPrefix(s, "./"), strings.HasPrefix(s, "."): // making some assumptions here
		s = getFileBasePath()
	case strings.HasSuffix(s, "/"):
		strings.TrimSuffix(s, "/")
	}

	words := strings.Split(s, "/")

	retArr := make([]string, 0, len(words))
	for _, w := range words {
		if w == RepositoryName {
			retArr = append(retArr, w)
			return strings.Join(retArr, "/")
		}
		retArr = append(retArr, w)
	}
	return ""
}
func getFileBasePath() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Dir(filename)
}
