package session

type EndpointOption struct {
	AutoMigration   bool `url:"-"`
	DryRun          bool `url:"-"`
	CreateTableOnly bool `url:"-"`
	MultiStatements bool
}

type AdaptorOption struct {
	ReadOnly bool
}

type OptionFunc func(*AdaptorOption)

func ReadOnly() OptionFunc {
	return func(o *AdaptorOption) {
		o.ReadOnly = true
	}
}
