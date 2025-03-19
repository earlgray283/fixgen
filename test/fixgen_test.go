package yo

import (
	"os"
	"path/filepath"
	"testing"

	goldiev2 "github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/require"

	"github.com/earlgray283/fixgen/internal/config"
	"github.com/earlgray283/fixgen/internal/gen"
	gen_ent "github.com/earlgray283/fixgen/internal/gen/ent"
	gen_yo "github.com/earlgray283/fixgen/internal/gen/yo"
)

func Test_GoldenTest(t *testing.T) {
	wd, err := os.Getwd()
	require.NoError(t, err)

	commonConfig := func() *config.Config {
		return &config.Config{
			Structs: config.Structs{
				"User": {
					Fields: map[string]*config.Field{
						"Name": {
							Value: "Taro Yamada",
						},
						"IconURL": {
							Expr: `fmt.Sprintf("http://example.com/%d", 123456)`,
						},
						"UserType": {
							Value:          1,
							IsModifiedCond: `m.UserType != 1`,
						},
					},
				},
				"Todo": {
					Fields: map[string]*config.Field{
						"Title": {
							MustOverwrite: true,
						},
					},
				},
			},
			Imports: []*config.Import{
				{Package: "fmt"},
			},
		}
	}

	generators := []string{"yo", "ent"}

	for _, typ := range generators {
		testDir := filepath.Join(wd, typ, "test")
		t.Run(typ, func(t *testing.T) {
			chdir(t, wd, testDir)
			defer chdir(t, testDir, wd)

			tcs := map[string]struct {
				opts       []gen.OptionFunc
				modify     func(*config.Config)
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
				"use-math-v1": {
					fixtureDir: "testdata-math-v1",
					modify: func(c *config.Config) {
						c.DefaultValuePolicy = &config.DefaultValuePolicy{
							Type: config.DefaultValuePolicyTypeRandv1,
						}
					},
				},
				"use-zero": {
					fixtureDir: "testdata-zero",
					modify: func(c *config.Config) {
						c.DefaultValuePolicy = &config.DefaultValuePolicy{
							Type: config.DefaultValuePolicyTypeZero,
						}
					},
				},
			}

			for name, tc := range tcs {
				tc := tc
				t.Run(name, func(t *testing.T) {
					c := commonConfig()
					if tc.modify != nil {
						tc.modify(c)
					}

					g := mustNewGenerator(t, typ)
					files, err := gen.GenerateWithFormat(g, c, tc.opts...)
					require.NoError(t, err)

					goldie := goldiev2.New(t, goldiev2.WithDiffEngine(goldiev2.ColoredDiff), goldiev2.WithNameSuffix(".go"), goldiev2.WithFixtureDir(filepath.Join(wd, typ, "test", tc.fixtureDir)))
					for _, f := range files {
						goldie.Assert(t, "goldie-"+f.Name, f.Content)
					}
				})
			}
		})
	}
}

func mustNewGenerator(t *testing.T, typ string) gen.Generator {
	t.Helper()

	var (
		g   gen.Generator
		err error
	)
	switch typ {
	case "yo":
		g, err = gen_yo.NewGenerator(".")
	case "ent":
		g, err = gen_ent.NewGenerator(".")
	default:
		t.Fatalf("unrecognized generator type `%s`", typ)
	}
	require.NoError(t, err)

	return g
}

func chdir(t *testing.T, from, to string) {
	t.Helper()

	t.Logf("chdir: `%s` → `%s`", from, to)
	require.NoError(t, os.Chdir(to))
}
