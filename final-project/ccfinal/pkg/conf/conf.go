package conf

import (
	"fmt"
	"log"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/structs"
)

// Config
// struct type of app configs.
type Config struct {
	Monitor MonitorConfig `koanf:"monitor"`
	Mysql   mysqlConfig   `koanf:"mysql"`
}

// Load
// loading app configs.
func Load() Config {
	var instance Config

	k := koanf.New(".")

	// load default
	if err := k.Load(structs.Provider(Default(), "koanf"), nil); err != nil {
		_ = fmt.Errorf("error loading deafult: %v\n", err)
	}

	// load configs file
	if err := k.Load(file.Provider("config.yaml"), yaml.Parser()); err != nil {
		_ = fmt.Errorf("error loading config.yaml file: %v\n", err)
	}

	// unmarshalling
	if err := k.Unmarshal("", &instance); err != nil {
		log.Fatalf("error unmarshalling config: %v\n", err)
	}

	return instance
}

func Default() Config {
	return Config{
		Monitor: MonitorConfig{
			Port:      0,
			IntervalM: 0,
		},

		Mysql: mysqlConfig{
			WAddr: "",
			RAddr: "",
		},
	}
}

type MonitorConfig struct {
	Port      int `koanf:"port"`
	IntervalM int `koanf:"intervalM"`
}

type mysqlConfig struct {
	WAddr string `koanf:"waddr"`
	RAddr string `koanf:"raddr"`
}
