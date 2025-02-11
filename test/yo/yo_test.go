package yo

import (
	"testing"

	"github.com/earlgray283/fixgen/internal/gen"
	gen_yo "github.com/earlgray283/fixgen/internal/gen/yo"
	goldiev2 "github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/require"
)

func Test_GoldenTest_yo(t *testing.T) {
	g, err := gen_yo.NewGenerator(gen.WithWorkDir("./project"))
	require.NoError(t, err)

	files, err := gen.GenerateWithFormat(g)
	require.NoError(t, err)

	goldie := goldiev2.New(t, goldiev2.WithDiffEngine(goldiev2.ColoredDiff), goldiev2.WithNameSuffix(".go"), goldiev2.WithFixtureDir("testdata"))
	for _, f := range files {
		goldie.Assert(t, "goldie-"+f.Name, f.Content)
	}
}
