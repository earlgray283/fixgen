# fixgen

fixgen is a tool for generating fixture packages for tools like yo, ent.

## Install

```shell
go install github.com/earlgray283/fixgen@latest
```

## Usage

You can start generating by specifying the tool name after `fixgen`.

```shell
fixgen yo
```

## Options

| Option Name      | Flag                   | Default Value | Description                                                                                                       |
| ---------------- | ---------------------- | ------------- | ----------------------------------------------------------------------------------------------------------------- |
| Prefix           | `--prefix`             | `<null>`      | prefix for names of generated Go files                                                                            |
| Extension        | `--ext`                | `.gen.go`     | extension for names of generated Go files                                                                         |
| PackageName      | `--package`            | `fixture`     | package name for generated Go files                                                                               |
| Output           | `--output`             | `.`           | destination directory(if DestDir is `foo` and PackageName is `var`, then the directory `foo/var` will be created) |
| SkipConfirm      | `--skip-confirm`       | `false`       | skip confirm before generation                                                                                    |
| UseContext       | `--use-context`        | `false`       | if `true`, `context.Context` argument will be added to fixture function                                           |
| UseValueModifier | `--use-value-modifier` | `false`       | if `true`, type of modifier struct will be value                                                                  |
| Config           | `--config`             | `fixgen.yaml` | location of fixgen configration file                                                                              |

## Config

See [example](https://github.com/earlgray283/fixgen/tree/main/.examples/fixgen.yaml)

## Support Tools

| Tool Name | Repository URL                                |
| --------- | --------------------------------------------- |
| yo        | <https://github.com/cloudspannerecosystem/yo> |
| ent       | <https://github.com/ent/ent>                  |

## Examples

- [Yo](https://github.com/earlgray283/fixgen/tree/main/.examples/yo)
- [Ent](https://github.com/earlgray283/fixgen/tree/main/.examples/ent)
