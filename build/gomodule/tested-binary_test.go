package gomodule

import (
	"bytes"
	"github.com/google/blueprint"
	"github.com/roman-mazur/bood"
	"strings"
	"testing"
)

var fileSystemDescriptions = []map[string][]byte{
	{
		"Blueprints": []byte(`
			go_binary {
			  name: "package-out",
			  pkg: ".",
              testPkg: ".",
			  srcs: [ "main.go",],
			  testSrcs: [ "main_test.go",],
			}
		`),
		"main.go":      nil,
		"main_test.go": nil,
	},
	{
		"Blueprints": []byte(`
			go_binary {
			  name: "new-package",
			  pkg: ".",
              testPkg: ".",
			  srcs: [ "main.go",],
			  testSrcs: [ "main_test.go",],
			  vendorFirst: true
			}
		`),
		"main.go":      nil,
		"main_test.go": nil,
	},
}

var expectedOutput = [][]string{
	{
		"out/bin/package-out:",
		"g.gomodule.binaryBuild | main.go\n",
		"out/reports/package-out/test.txt",
		"g.gomodule.test | main_test.go main.go",
	},
	{
		"out/bin/new-package",
		"g.gomodule.binaryBuild | main.go\n",
		"build vendor: g.gomodule.vendor | go.mod\n",
		"out/reports/new-package/test.txt",
		"g.gomodule.test | main_test.go main.go",
	},
}



func TestTestCoverageFactory(t *testing.T) {
	for i, fs := range fileSystemDescriptions {
		t.Run(string(i), func(t *testing.T) {
			ctx := blueprint.NewContext()

			ctx.MockFileSystem(fs)

			ctx.RegisterModuleType("go_binary", SimpleBinFactory)

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