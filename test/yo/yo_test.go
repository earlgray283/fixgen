package yo

import (
	"os"
	"testing"

	goldiev2 "github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/require"

	"github.com/earlgray283/fixgen/internal/config"
	"github.com/earlgray283/fixgen/internal/gen"
	gen_yo "github.com/earlgray283/fixgen/internal/gen/yo"
)

func Test_GoldenTest_yo(t *testing.T) {
	require.NoError(t, os.Chdir("./test"))

	c := &config.Config{
		Structs: config.Structs{
			"User": {
				Fields: map[string]*config.Field{
					"Name": {
						Value: "Taro Yamada",
					},
					"IconURL": {
						Expr: `fmt.Sprintf("http://example.com/%d", rand.Int64())`,
					},
				},
			},
		},
		Imports: []*config.Import{
			{Package: "fmt"},
		},
	}

	g, err := gen_yo.NewGenerator(".", false, true)
	require.NoError(t, err)
	files, err := gen.GenerateWithFormat(g, c)
	require.NoError(t, err)

	goldie := goldiev2.New(t, goldiev2.WithDiffEngine(goldiev2.ColoredDiff), goldiev2.WithNameSuffix(".go"), goldiev2.WithFixtureDir("testdata"))
	for _, f := range files {
		goldie.Assert(t, "goldie-"+f.Name, f.Content)
	}

	g, err = gen_yo.NewGenerator(".", true, false)
	require.NoError(t, err)
	files, err = gen.GenerateWithFormat(g, c)
	require.NoError(t, err)

	goldie = goldiev2.New(t, goldiev2.WithDiffEngine(goldiev2.ColoredDiff), goldiev2.WithNameSuffix(".go"), goldiev2.WithFixtureDir("testdata-context"))
	for _, f := range files {
		goldie.Assert(t, "goldie-"+f.Name, f.Content)
	}
}
