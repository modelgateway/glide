package routers

import (
	"errors"

	"glide/pkg/telemetry"
)

var ErrRouterNotFound = errors.New("no router found with given ID")

type RouterManager struct {
	Config        *Config
	telemetry     *telemetry.Telemetry
	langRouterMap *map[string]*LangRouter
	langRouters   []*LangRouter
}

// NewManager creates a new instance of Router Manager that creates, holds and returns all routers
func NewManager(cfg *Config, tel *telemetry.Telemetry) (*RouterManager, error) {
	langRouters, err := cfg.BuildLangRouters(tel)
	if err != nil {
		return nil, err
	}

	langRouterMap := make(map[string]*LangRouter, len(langRouters))

	for _, router := range langRouters {
		langRouterMap[router.ID()] = router
	}

	manager := RouterManager{
		Config:        cfg,
		telemetry:     tel,
		langRouters:   langRouters,
		langRouterMap: &langRouterMap,
	}

	return &manager, err
}

func (r *RouterManager) GetLangRouters() []*LangRouter {
	return r.langRouters
}

// GetLangRouter returns a router by type and ID
func (r *RouterManager) GetLangRouter(routerID string) (*LangRouter, error) {
	if router, found := (*r.langRouterMap)[routerID]; found {
		return router, nil
	}

	return nil, ErrRouterNotFound
}
