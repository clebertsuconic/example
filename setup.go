package example

import (
	"fmt"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/metrics"

	"github.com/caddyserver/caddy"
)

// init registers this plugin.
func init() { 
	plugin.Register("clebert", setup) 
        fmt.Println("Registered")
}

// setup is the function that gets called when the config parser see the token "example". Setup is responsible
// for parsing any extra options the example plugin may have. The first token this function sees is "example".
func setup(c *caddy.Controller) error {
	fmt.Println("setting up clebert")
	c.Next() // Ignore "example" and give us the next token.
	if c.NextArg() {
		// If there was another token, return an error, because we don't have any configuration.
		// Any errors returned from this setup function should be wrapped with plugin.Error, so we
		// can present a slightly nicer error message to the user.
		return plugin.Error("clebert", c.ArgErr())
	}

	// Add a startup function that will -- after all plugins have been loaded -- check if the
	// prometheus plugin has been used - if so we will export metrics. We can only register
	// this metric once, hence the "once.Do".
	c.OnStartup(func() error {
		once.Do(func() { metrics.MustRegister(c, requestCount) })
		return nil
	})

	// Add the Plugin to CoreDNS, so Servers can use it in their plugin chain.
	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		return Clebert{Next: next}
	})

	// All OK, return a nil error.
	return nil
}
