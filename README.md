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

| Option Name        | Flag                    | Default Value | Description                                                                                                       |
| ------------------ | ----------------------- | ------------- | ----------------------------------------------------------------------------------------------------------------- |
| Prefix             | `-prefix`               | `<null>`      | prefix for names of generated Go files                                                                            |
| Extension          | `-ext`                  | `.gen.go`     | extension for names of generated Go files                                                                         |
| PackageName        | `-pkgname`              | `fixture`     | package name for generated Go files                                                                               |
| DestDir            | `-dest-dir`             | `.`           | destination directory(if DestDir is `foo` and PackageName is `var`, then the directory `foo/var` will be created) |
| CleanIfFailed      | `-clean-if-failed`      | `false`       | clean generated files and directories if the generation failed                                                    |
| Confirm            | `-confirm`              | `true`        | confirm before generation                                                                                         |
| UseContext         | `-use-context`          | `false`       | if `true`, `context.Context` argument will be added to fixture function                                           |
| UsePointerModified | `-use-pointer-modified` | `true`        | if `true`, type of modifier struct will be pointer                                                                |
| Config             | `-config`               | `fixgen.yaml` | location of fixgen configration file                                                                              |

## Config

See [example](https://github.com/earlgray283/fixgen/tree/main/.examples/fixgen.yaml)

## Support Tools

| Tool Name | Repository URL                                | Support | Experimental |
| --------- | --------------------------------------------- | ------- | ------------ |
| yo        | <https://github.com/cloudspannerecosystem/yo> | Yes     | Yes          |
| ent       | <https://github.com/ent/ent>                  | Yes     | Yes          |

## Examples

- [Yo](https://github.com/earlgray283/fixgen/tree/main/.examples/yo)
- [Ent](https://github.com/earlgray283/fixgen/tree/main/.examples/ent)
