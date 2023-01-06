package config

import (
	"flag"
	"fmt"
	"github.com/spf13/viper"
	"os"
)

type Config struct {
	Generator        GeneratorConfig
	OpenAPIFile      OpenApiConfig
	OutputPath       string
	WorkingDirectory string
	ConfigFilePath   string
}
type OpenApiConfig struct {
	OpenAPIFile string
}

type GeneratorConfig struct {
	ModuleName        string
	BasePackageName   string
	TemplateDirectory string
}

func LoadConfig() (Config, error) {
	wd, err := os.Getwd()
	if err != nil {
		return Config{}, fmt.Errorf("could not determine working directory: %w\n", err)
	}

	args := flag.NewFlagSet("can", flag.ExitOnError)

	var configFilePath = args.String("configFile", "", "Specify which config file to use")
	_ = args.Parse(os.Args[1:])

	if configFilePath == nil {
		fmt.Println("No config file specified.")
		viper.SetConfigName("config")
		viper.AddConfigPath(".")
	} else {
		fmt.Printf("Using config file \"%s\" as specified.\n", *configFilePath)
		viper.SetConfigFile(*configFilePath)
	}

	err = viper.ReadInConfig()
	if err != nil {
		return Config{}, fmt.Errorf("could not read config file: %w\n", err)
	}

	configData := Config{
		WorkingDirectory: wd,
		ConfigFilePath:   viper.ConfigFileUsed(),
	}

	_ = viper.Unmarshal(&configData)

	return configData, nil
}
