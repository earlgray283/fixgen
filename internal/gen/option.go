package gen

type option struct {
	packageName      string // default: "fixture"
	useContext       bool
	useValueModifier bool
	targetDir        string
}

func defaultOption() *option {
	return &option{
		packageName:      "fixture",
		useContext:       false,
		useValueModifier: false,
		targetDir:        "",
	}
}

type OptionFunc func(*option)

func (o *option) applyOptionFuncs(opts ...OptionFunc) {
	for _, opt := range opts {
		opt(o)
	}
}

func WithPackageName(packageName string) OptionFunc {
	return func(o *option) {
		o.packageName = packageName
	}
}

func UseContext() OptionFunc {
	return func(o *option) {
		o.useContext = true
	}
}

func UseValueModifier() OptionFunc {
	return func(o *option) {
		o.useValueModifier = true
	}
}

func TargetDir(targetDir string) OptionFunc {
	return func(o *option) {
		o.targetDir = targetDir
	}
}
