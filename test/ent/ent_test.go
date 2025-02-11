package ent

import (
	"os"
	"testing"

	"github.com/earlgray283/fixgen/internal/gen"
	gen_ent "github.com/earlgray283/fixgen/internal/gen/ent"
	goldiev2 "github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/require"
)

func Test_GoldenTest_ent(t *testing.T) {
	// avoid failed to ent load: entc/load: parse schema dir: -: main module (github.com/earlgray283/fixgen) does not contain package github.com/earlgray283/fixgen/test/ent/project/ent/schema
	require.NoError(t, os.Chdir("./test"))

	g, err := gen_ent.NewGenerator()
	require.NoError(t, err)

	files, err := gen.GenerateWithFormat(g)
	require.NoError(t, err)

	goldie := goldiev2.New(t, goldiev2.WithDiffEngine(goldiev2.ColoredDiff), goldiev2.WithNameSuffix(".go"), goldiev2.WithFixtureDir("testdata"))
	for _, f := range files {
		goldie.Assert(t, "goldie-"+f.Name, f.Content)
	}
}
