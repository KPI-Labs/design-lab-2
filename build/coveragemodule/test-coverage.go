package coveragemodule

import (
	"fmt"
	"path"

	"github.com/google/blueprint"
	"github.com/roman-mazur/bood"
)

var (
	pctx = blueprint.NewPackageContext("github.com/KPI-Labs/design-lab-2/build/coveragemodule")

	goTestCoverage = pctx.StaticRule("testCoverage", blueprint.RuleParams{
		Command:     "cd $workDir && mkdir -p $outputPath && go test -v $pkg -coverprofile=$outputCoverage && go tool cover -html=$outputCoverage -o $outputHtml",
		Description: "test coverage for $pkg",
	}, "workDir", "pkg", "outputPath", "outputCoverage", "outputHtml")
)

type testCoverageModule struct {
	blueprint.SimpleName

	properties struct {
		Name string
		Pkg  string
		Srcs []string
	}
}

func convertPatternsIntoPaths(ctx blueprint.ModuleContext, patterns []string, excludePatterns []string) []string {
	var paths []string
	for _, src := range patterns {
		if matches, err := ctx.GlobWithDeps(src, excludePatterns); err == nil {
			paths = append(paths, matches...)
		} else {
			ctx.PropertyErrorf("srcs", "Cannot resolve files that match pattern %s", src)
			return nil
		}
	}
	return paths
}

func (gb *testCoverageModule) GenerateBuildActions(ctx blueprint.ModuleContext) {
	name := ctx.ModuleName()
	config := bood.ExtractConfig(ctx)

	pathToReports := path.Join(config.BaseOutputDir, "reports", name)
	pathToCoverageReports := path.Join(pathToReports, "coverage.out")
	pathToCoverageHtml := path.Join(pathToReports, "coverage.html")

	inputs := convertPatternsIntoPaths(ctx, gb.properties.Srcs, []string{})

	if inputs != nil {
		ctx.Build(pctx, blueprint.BuildParams{
			Description: fmt.Sprintf("Test coverage for %s", name),
			Rule:        goTestCoverage,
			Outputs:     []string{config.BaseOutputDir},
			Implicits:   inputs,
			Args: map[string]string{
				"outputPath":     pathToReports,
				"outputCoverage": pathToCoverageReports,
				"outputHtml":     pathToCoverageHtml,
				"workDir":        ctx.ModuleDir(),
				"pkg":            gb.properties.Pkg,
			},
		})
	}
}

func TestCoverageFactory() (blueprint.Module, []interface{}) {
	mType := &testCoverageModule{}
	return mType, []interface{}{&mType.SimpleName.Properties, &mType.properties}
}
