package plugin1

import (
	log "github.com/Sirupsen/logrus"
	"api"
	"middleware"
	"oauth"
	"router"
	"store"
)

// OAuth2 checks the integrity of the provided OAuth headers
type OAuth2 struct {
	authRepo oauth.Repository
	storage  store.Store
}

// NewOAuth2 creates a new instance of KeyExistsMiddleware
func NewOAuth2(authRepo oauth.Repository, storage store.Store) *OAuth2 {
	return &OAuth2{authRepo, storage}
}

// GetName retrieves the plugin's name
func (h *OAuth2) GetName() string {
	return "oauth2"
}

// GetMiddlewares retrieves the plugin's middlewares
func (h *OAuth2) GetMiddlewares(config api.Config, referenceSpec *api.Spec) ([]router.Constructor, error) {
	manager, err := h.getManager(config["server_name"].(string))
	if nil != err {
		log.WithError(err).Error("OAuth Configuration for this API is incorrect, skipping...")
		return nil, err
	}

	mw := middleware.NewKeyExistsMiddleware(manager)
	return []router.Constructor{
		mw.Handler,
	}, nil
}

func (h *OAuth2) getManager(oAuthServerName string) (oauth.Manager, error) {
	oauthServer, err := h.authRepo.FindByName(oAuthServerName)
	if nil != err {
		return nil, err
	}

	managerType, err := oauth.ParseType(oauthServer.TokenStrategy.Name)
	if nil != err {
		return nil, err
	}

	return oauth.NewManagerFactory(h.storage, oauthServer.TokenStrategy.Settings).Build(managerType)
}
