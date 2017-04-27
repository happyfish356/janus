package loader

import (
	"api"
	"oauth"
	plugin "plugin1"
	"proxy"
	"router"
	"store"
)

// Params initialization options.
type Params struct {
	Router      router.Router
	Storage     store.Store
	APIRepo     api.Repository
	OAuthRepo   oauth.Repository
	ProxyParams proxy.Params
}

// Load loads all the basic components and definitions into a router
func Load(params Params) {
	pluginLoader := plugin.NewLoader()
	pluginLoader.Add(
		plugin.NewRateLimit(params.Storage),
		plugin.NewCORS(),
		plugin.NewOAuth2(params.OAuthRepo, params.Storage),
		plugin.NewCompression(),
	)

	prx := proxy.WithParams(params.ProxyParams)

	// create proxy register
	register := proxy.NewRegister(params.Router, prx)

	apiLoader := NewAPILoader(register, pluginLoader)
	apiLoader.LoadDefinitions(params.APIRepo)

	oauthLoader := NewOAuthLoader(register, params.Storage)
	oauthLoader.LoadDefinitions(params.OAuthRepo)
}
