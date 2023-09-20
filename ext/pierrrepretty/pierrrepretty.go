// Package pierrrepretty provides an integration with github.com/pierrre/pretty.
package pierrrepretty

import (
	"github.com/pierrre/errors/errval"
	"github.com/pierrre/pretty"
)

// Configure configures the integration.
//
// It sets errval.VerboseStringer with pretty.Config.String.
func Configure(config *pretty.Config) {
	errval.VerboseWriter = config.Write
}

// ConfigureDefault calls Configure() with pretty.DefaultConfig.
func ConfigureDefault() {
	Configure(pretty.DefaultConfig)
}
