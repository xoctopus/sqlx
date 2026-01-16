package session

import (
	"net/url"

	"github.com/xoctopus/x/misc/must"
	"github.com/xoctopus/x/textx"
)

type EndpointOption struct {
	AutoMigration   bool `url:"-"`
	DryRun          bool `url:"-"`
	CreateTableOnly bool `url:"-"`
	MultiStatements bool
}

func (o *EndpointOption) SetDefault() {
	must.NoError(textx.UnmarshalURL(url.Values{}, o))
	o.MultiStatements = true
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
