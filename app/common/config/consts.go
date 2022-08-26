package config

import "errors"

var (
	errEmptyConfigClient = errors.New("configClient is null(try to initialize a new one)")
	errGetViper          = errors.New("get viper failed")
	errViperEmptyKey     = errors.New("get viper key failed")
)
