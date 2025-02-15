package gen

type option struct {
	packageName string // default: "fixture"
}

func defaultOption() *option {
	return &option{
		packageName: "fixture",
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
