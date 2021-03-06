package plugin

import (
	"encoding/json"

	"github.com/hellofresh/janus/pkg/api"
	"github.com/hellofresh/janus/pkg/router"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/cors"
)

type corsConfig struct {
	AllowedOrigins []string `json:"domains"`
	AllowedMethods []string `json:"methods"`
	AllowedHeaders []string `json:"request_headers"`
	ExposedHeaders []string `json:"exposed_headers"`
}

// CORS represents the cors plugin
type CORS struct{}

// NewCORS creates a new instance of CORS
func NewCORS() *CORS {
	return &CORS{}
}

// GetName retrieves the plugin's name
func (h *CORS) GetName() string {
	return "cors"
}

// GetMiddlewares retrieves the plugin's middlewares
func (h *CORS) GetMiddlewares(rawConfig map[string]interface{}, referenceSpec *api.Spec) ([]router.Constructor, error) {
	var corsConfig corsConfig
	err := mapstructure.Decode(rawConfig, &corsConfig)
	if err != nil {
		return nil, err
	}

	// for some reasons mapstructure.Decode() gives empty arrays for all resulting config fields
	// this is quick workaround hack t make it work
	// TODO: investigate and fix mapstructure.Decode() behaviour and remove this dirty hack
	valJSON, err := json.Marshal(rawConfig)
	if nil != err {
		return nil, err
	}
	err = json.Unmarshal(valJSON, &corsConfig)
	if nil != err {
		return nil, err
	}

	mw := cors.New(cors.Options{
		AllowedOrigins:   corsConfig.AllowedOrigins,
		AllowedMethods:   corsConfig.AllowedMethods,
		AllowedHeaders:   corsConfig.AllowedHeaders,
		ExposedHeaders:   corsConfig.ExposedHeaders,
		AllowCredentials: true,
	})

	return []router.Constructor{
		mw.Handler,
	}, nil
}
