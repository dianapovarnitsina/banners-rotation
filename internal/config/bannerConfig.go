package config

import (
	"github.com/pkg/errors"
)

var _ Configure = (*BannerConfig)(nil)

type BannerConfig struct {
	Logger   LoggerConf   `json:"logger"`
	FilePath string       `json:"file_path"` //nolint:tagliatelle
	Database DataBaseConf `json:"database"`
	GRPC     GRPC         `json:"grpc"`
	Storage  StorageConf  `json:"storage"`
	RMQ      RMQ          `json:"rmq"`
	Queues   struct {
		Events Queue
	}
}

func (b *BannerConfig) Init(file string) error {
	cfg, err := Init(file, b)

	_, ok := cfg.(*BannerConfig)
	if !ok {
		return errors.Wrap(err, "init config failed")
	}

	return nil
}
