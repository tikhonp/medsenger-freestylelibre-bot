// Code generated from Pkl module `config`. DO NOT EDIT.
package config

import (
	"context"

	"github.com/apple/pkl-go/pkl"
)

type Config struct {
	Server *Server `pkl:"server"`

	Db *Database `pkl:"db"`

	// The duration of the sleep between the requests to the LibreView API.
	FetchSleepDuration *pkl.Duration `pkl:"fetchSleepDuration"`
}

// LoadFromPath loads the pkl module at the given path and evaluates it into a Config
func LoadFromPath(ctx context.Context, path string) (ret *Config, err error) {
	evaluator, err := pkl.NewEvaluator(ctx, pkl.PreconfiguredOptions)
	if err != nil {
		return nil, err
	}
	defer func() {
		cerr := evaluator.Close()
		if err == nil {
			err = cerr
		}
	}()
	ret, err = Load(ctx, evaluator, pkl.FileSource(path))
	return ret, err
}

// Load loads the pkl module at the given source and evaluates it with the given evaluator into a Config
func Load(ctx context.Context, evaluator pkl.Evaluator, source *pkl.ModuleSource) (*Config, error) {
	var ret Config
	if err := evaluator.EvaluateModule(ctx, source, &ret); err != nil {
		return nil, err
	}
	return &ret, nil
}
