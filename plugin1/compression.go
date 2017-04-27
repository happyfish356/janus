package plugin1

import (
	"api"
	"github.com/pressly/chi/middleware"
	"router"
)

// Compression represents the compression plugin
type Compression struct{}

// NewCompression creates a new instance of Compression
func NewCompression() *Compression {
	return &Compression{}
}

// GetName retrieves the plugin's name
func (h *Compression) GetName() string {
	return "compression"
}

// GetMiddlewares retrieves the plugin's middlewares
func (h *Compression) GetMiddlewares(config api.Config, referenceSpec *api.Spec) ([]router.Constructor, error) {
	return []router.Constructor{
		middleware.DefaultCompress,
	}, nil
}
