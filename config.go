package ermeslog

import (
	"io"

	"github.com/sirupsen/logrus"
)

type Config struct {
	Business    string
	Service     string
	Version     string
	Environment string
	Level       logrus.Level
	Output      io.Writer
	FileOutput  *FileConfig
	Hooks       []logrus.Hook
}

type FileConfig struct {
	Filename   string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Compress   bool
}

// DefaultConfig returns a config with sensible defaults
func DefaultConfig() Config {
	return Config{
		Level:       logrus.InfoLevel,
		Environment: "development",
	}
}
