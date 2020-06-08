package coveragemodule

import (
	"bytes"
	"strings"
	"testing"

	"github.com/google/blueprint"
	"github.com/roman-mazur/bood"
)

var fileSystemDescriptions = []map[string][]byte{
	{
		"Blueprints": []byte(`
			test_coverage {
			  name: "package-out",
			  pkg: ".",
			  srcs: [ "main_test.go", "main.go",],
			}
		`),
		"main.go":      nil,
		"main_test.go": nil,
	},
}

var expectedOutput = [][]string{
	{
		"out:",
		"g.coveragemodule.testCoverage | main_test.go main.go",
		"description = Test coverage for package-out",
		"outputCoverage = out/reports/package-out/coverage.out",
		"outputHtml = out/reports/package-out/coverage.html",
	},
}

func TestTestCoverageFactory(t *testing.T) {
	for i, fs := range fileSystemDescriptions {
		t.Run(string(i), func(t *testing.T) {
			ctx := blueprint.NewContext()

			ctx.MockFileSystem(fs)

			ctx.RegisterModuleType("test_coverage", TestCoverageFactory)

			cfg := bood.NewConfig()

			_, errs := ctx.ParseBlueprintsFiles(".", cfg)

			if len(errs) != 0 {
				t.Fatalf("Syntax errors in the test blueprint file: %s", errs)
			}

			_, errs = ctx.PrepareBuildActions(cfg)

			if len(errs) != 0 {
				t.Errorf("Unexpected errors while preparing build actions: %s", errs)
			}

			buffer := new(bytes.Buffer)

			if err := ctx.WriteBuildFile(buffer); err != nil {

				t.Errorf("Error writing ninja file: %s", err)

			} else {

				text := buffer.String()
				//t.Logf("Generated ninja build file:\n%s", text) //For debug purposes

				for _, expectedStr := range expectedOutput[i] {

					if strings.Contains(text, expectedStr) != true {
						t.Errorf("Generated ninja file does not have expected string `%s`", expectedStr)
					}

				}

			}

		})
	}
}
