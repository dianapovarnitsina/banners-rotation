package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/dianapovarnitsina/banners-rotation/internal/app/banner"
	"github.com/dianapovarnitsina/banners-rotation/internal/config"
	"github.com/pkg/errors"
	"log"
	"sync"
)

var (
	bannerConfigFile string
	wg               sync.WaitGroup
)

func init() {
	flag.StringVar(&bannerConfigFile, "config", "banner_config.yaml", "Path to configuration file")
}

func main() {
	if err := mainImpl(); err != nil {
		log.Fatal(err)
	}
}

func mainImpl() error {
	ctx := context.TODO()
	flag.Parse()

	if bannerConfigFile == "" {
		return fmt.Errorf("please set: '--config=<Path to configuration file>'")
	}

	conf := new(config.BannerConfig)
	if err := conf.Init(bannerConfigFile); err != nil {
		return errors.Wrap(err, "init config failed")
	}

	app, err := banner.NewApp(ctx, conf)
	if err != nil {
		return fmt.Errorf("failed to create bannerApp: %w", err)
	}

	wg.Add(1)

	go func() {
		defer wg.Done()
		<-app.GetGrpcServerShutdownSignal()
	}()

	wg.Wait()

	return nil
}
