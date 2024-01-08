package routers

import (
	"glide/pkg/providers"
	"glide/pkg/routers/retry"
	"glide/pkg/routers/routing"
	"glide/pkg/telemetry"
	"go.uber.org/multierr"
	"go.uber.org/zap"
)

type Config struct {
	LanguageRouters []LangRouterConfig `yaml:"language"` // the list of language routers
}

func (c *Config) BuildLangRouters(tel *telemetry.Telemetry) ([]*LangRouter, error) {
	routers := make([]*LangRouter, 0, len(c.LanguageRouters))

	var errs error

	for idx, routerConfig := range c.LanguageRouters {
		if !routerConfig.Enabled {
			tel.Logger.Info("router is disabled, skipping", zap.String("routerID", routerConfig.ID))
			continue
		}

		tel.Logger.Debug("init router", zap.String("routerID", routerConfig.ID))

		router, err := NewLangRouter(&c.LanguageRouters[idx], tel)
		if err != nil {
			errs = multierr.Append(errs, err)
			continue
		}

		routers = append(routers, router)
	}

	if errs != nil {
		return nil, errs
	}

	return routers, nil
}

// TODO: how to specify other backoff strategies?
// LangRouterConfig
type LangRouterConfig struct {
	ID              string                      `yaml:"id" json:"routers" validate:"required"`    // Unique router ID
	Enabled         bool                        `yaml:"enabled" json:"enabled"`                   // Is router enabled?
	Retry           *retry.ExpRetryConfig       `yaml:"retry" json:"retry"`                       // retry when no healthy model is available to router
	RoutingStrategy routing.Strategy            `yaml:"strategy" json:"strategy"`                 // strategy on picking the next model to serve the request
	Models          []providers.LangModelConfig `yaml:"models" json:"models" validate:"required"` // the list of models that could handle requests
}

// BuildModels creates LanguageModel slice out of the given config
func (c *LangRouterConfig) BuildModels(tel *telemetry.Telemetry) ([]*providers.LangModel, error) {
	var errs error

	models := make([]*providers.LangModel, 0, len(c.Models))

	for _, modelConfig := range c.Models {
		if !modelConfig.Enabled {
			tel.Logger.Info(
				"model is disabled, skipping",
				zap.String("router", c.ID),
				zap.String("model", modelConfig.ID),
			)

			continue
		}

		tel.Logger.Debug(
			"init lang model",
			zap.String("router", c.ID),
			zap.String("model", modelConfig.ID),
		)

		model, err := modelConfig.ToModel(tel)
		if err != nil {
			errs = multierr.Append(errs, err)
			continue
		}

		models = append(models, model)
	}

	if errs != nil {
		return nil, errs
	}

	return models, nil
}

func (c *LangRouterConfig) BuildRetry() *retry.ExpRetry {
	retryConfig := c.Retry

	return retry.NewExpRetry(
		retryConfig.MaxRetries,
		retryConfig.BaseMultiplier,
		retryConfig.MinDelay,
		retryConfig.MaxDelay,
	)
}

func DefaultLangRouterConfig() LangRouterConfig {
	return LangRouterConfig{
		Enabled:         true,
		RoutingStrategy: routing.Priority,
		Retry:           retry.DefaultExpRetryConfig(),
	}
}

func (c *LangRouterConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	*c = DefaultLangRouterConfig()

	type plain LangRouterConfig // to avoid recursion

	return unmarshal((*plain)(c))
}
