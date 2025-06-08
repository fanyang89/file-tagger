package main

import (
	"context"
	"os"
	"path/filepath"

	"github.com/cockroachdb/errors"
	"github.com/goccy/go-yaml"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v3"
	gormzerolog "github.com/vitaliy-art/gorm-zerolog"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	_ "github.com/mattn/go-sqlite3"

	"github.com/fanyang89/file-tagger/ft/v1"
	"github.com/fioepq9/pzlog"
)

type Config struct {
	DSN string `yaml:"dsn"`
}

func defaultConfig() *Config {
	return &Config{}
}

func (c *Config) Validate() error {
	if c.DSN == "" {
		return errors.New("dsn is required")
	}
	return nil
}

func loadConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, "open config file")
	}
	defer func() { _ = f.Close() }()

	var config Config
	decoder := yaml.NewDecoder(f, yaml.Strict(), yaml.DisallowUnknownField())
	err = decoder.Decode(&config)
	if err != nil {
		return nil, errors.Wrap(err, "decode config file")
	}
	return &config, nil
}

func saveConfig(path string, config *Config) error {
	err := os.MkdirAll(filepath.Dir(path), 0755)
	if err != nil {
		return errors.Wrap(err, "save config")
	}

	buf, err := yaml.Marshal(config)
	if err != nil {
		return errors.Wrap(err, "save config")
	}

	return errors.Wrap(os.WriteFile(path, buf, 0644), "save config")
}

var cmd = &cli.Command{
	Name:  "ft",
	Usage: "File tagging tool",
	Commands: []*cli.Command{
		cmdTag,
	},
}

var cmdTag = &cli.Command{
	Name: "tag",
	Arguments: []cli.Argument{
		&cli.StringArg{Name: "path"},
		&cli.StringArgs{Name: "tags", Max: 10},
	},
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "config",
			Aliases: []string{"c"},
			Value:   filepath.Join(homeDir(), ".config", "file-tagger", "config.yaml"),
		},
	},
	Action: func(ctx context.Context, command *cli.Command) error {
		configPath := command.String("config")
		config, err := loadConfig(configPath)
		if err != nil {
			if !os.IsNotExist(err) {
				return err
			}
			config = defaultConfig()
			err = saveConfig(configPath, config)
			if err != nil {
				return err
			}
		}

		err = config.Validate()
		if err != nil {
			return err
		}

		filePath := command.StringArg("path")
		if filePath == "" {
			return errors.New("argument path is required")
		}
		_, err = os.Stat(filePath)
		if err != nil {
			if os.IsNotExist(err) {
				return errors.Errorf("file %s does not exist", filePath)
			}
			return err
		}

		filePath, err = filepath.Abs(filePath)
		if err != nil {
			return err
		}

		tagArgs := command.StringArgs("tags")
		if len(tagArgs) == 0 {
			return errors.New("argument tags is required")
		}

		db, err := gorm.Open(sqlite.Open(config.DSN), &gorm.Config{
			Logger: gormzerolog.NewGormLogger(),
		})
		if err != nil {
			return errors.Wrap(err, "open db")
		}

		tagger := ft.NewTagger(db)
		err = tagger.Migrate()
		if err != nil {
			return err
		}

		for _, s := range tagArgs {
			name, value := ft.ParseTagKeyValue(s)
			err = tagger.Tag(filePath, name, value)
			if err != nil {
				return errors.Wrap(err, "tag")
			}
		}

		return nil
	},
}

func homeDir() string {
	h, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return h
}

func main() {
	zerolog.SetGlobalLevel(zerolog.TraceLevel)
	log.Logger = zerolog.New(pzlog.NewPtermWriter()).With().Timestamp().Caller().Stack().Logger()
	err := cmd.Run(context.Background(), os.Args)
	if err != nil {
		log.Error().Err(err).Msg("Unexpected error")
	}
}
