package yo

import (
	"os"
	"testing"

	goldiev2 "github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/require"

	"github.com/earlgray283/fixgen/internal/gen"
	gen_yo "github.com/earlgray283/fixgen/internal/gen/yo"
)

func Test_GoldenTest_yo(t *testing.T) {
	require.NoError(t, os.Chdir("./test"))

	g, err := gen_yo.NewGenerator(".")
	require.NoError(t, err)

	files, err := gen.GenerateWithFormat(g)
	require.NoError(t, err)

	goldie := goldiev2.New(t, goldiev2.WithDiffEngine(goldiev2.ColoredDiff), goldiev2.WithNameSuffix(".go"), goldiev2.WithFixtureDir("testdata"))
	for _, f := range files {
		goldie.Assert(t, "goldie-"+f.Name, f.Content)
	}
}
