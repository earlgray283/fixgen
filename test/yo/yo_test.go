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
					"UserType": {
						Value:          1,
						IsModifiedCond: `m.UserType != 1`,
					},
				},
			},
		},
		Imports: []*config.Import{
			{Package: "fmt"},
		},
	}

	tcs := map[string]struct {
		opts       []gen.OptionFunc
		fixtureDir string
	}{
		"minimum_option": {
			fixtureDir: "testdata",
		},
		"use-context": {
			fixtureDir: "testdata-context",
			opts:       []gen.OptionFunc{gen.UseContext()},
		},
		"use-value-modifier": {
			fixtureDir: "testdata-value-modifier",
			opts:       []gen.OptionFunc{gen.UseValueModifier()},
		},
	}

	for name, tc := range tcs {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			g, err := gen_yo.NewGenerator(".")
			require.NoError(t, err)
			files, err := gen.GenerateWithFormat(g, c, tc.opts...)
			require.NoError(t, err)

			goldie := goldiev2.New(t, goldiev2.WithDiffEngine(goldiev2.ColoredDiff), goldiev2.WithNameSuffix(".go"), goldiev2.WithFixtureDir(tc.fixtureDir))
			for _, f := range files {
				goldie.Assert(t, "goldie-"+f.Name, f.Content)
			}
		})
	}
}
