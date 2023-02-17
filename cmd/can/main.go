package main

import (
	"flag"
	"fmt"
	"github.com/sasswart/gin-in-a-can/config"
	"github.com/sasswart/gin-in-a-can/openapi/operation"
	"github.com/sasswart/gin-in-a-can/openapi/path"
	"github.com/sasswart/gin-in-a-can/openapi/root"
	"github.com/sasswart/gin-in-a-can/openapi/schema"
	"github.com/sasswart/gin-in-a-can/render"
	"github.com/sasswart/gin-in-a-can/tree"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

func main() {
	configData, err := loadConfig()
	if err != nil {
		fmt.Println(fmt.Errorf("loadConfig error: %w", err))
		os.Exit(1)
	}

	fmt.Printf("Reading API specification from \"%s\"\n", absoluteOpenAPIFile(configData))
	apiSpec, err := root.LoadAPISpec(
		absoluteOpenAPIFile(configData),
	)
	if err != nil {
		fmt.Println(fmt.Errorf("openapi.LoadAPISpec error: %w", err))
		os.Exit(1)
	}

	engine := render.Engine{}
	canRenderer := engine.New(render.GinRenderer{})

	err = canRenderer.SetRenderer(apiSpec)
	if err != nil {
		fmt.Println(fmt.Errorf("openapi.SetRenderer error: %w", err))
		os.Exit(1)
	}

	apiSpec.SetMetadata(map[string]string{
		"package": configData.Generator.BasePackageName,
	})

	renderNode := buildRenderNode(configData)
	_, err = tree.Traverse(apiSpec, renderNode)
	if err != nil {
		fmt.Println(fmt.Errorf("openapi.Traverse(apiSpec, renderNode) error: %w", err))
		os.Exit(1)
	}
}

func buildRenderNode(config config.Config) tree.TraversalFunc {
	return func(key string, parent, child tree.NodeTraverser) (tree.NodeTraverser, error) {
		var templateFile string
		switch child.(type) {
		case *root.Root:
			templateFile = "openapi.tmpl"
		case *path.Item:
			templateFile = "path_item.tmpl"
		case *schema.Schema:
			schemaType := child.(*schema.Schema).Type
			if schemaType != "object" && schemaType != "array" {
				return child, nil
			}
			templateFile = "schema.tmpl"
		case *operation.Operation:
			templateFile = "operation.tmpl"
		}

		if templateFile == "" {
			return child, nil
		}
		_, err := render.Render(config, child, templateFile)
		if err != nil {
			return child, err
		}

		return child, nil
	}
}

func loadConfig() (config.Config, error) {
	exe, err := os.Readlink("/proc/self/exe")
	if err != nil {
		return config.Config{}, fmt.Errorf("could not read /proc/self/exe: %w\n", err)
	}

	wd, err := os.Getwd()
	if err != nil {
		return config.Config{}, fmt.Errorf("could not determine working directory: %w\n", err)
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
		return config.Config{}, fmt.Errorf("could not read config file: %w\n", err)
	}

	configData := config.Config{
		WorkingDirectory: wd,
		ConfigFilePath:   viper.ConfigFileUsed(),
		Generator: config.GeneratorConfig{
			TemplateDirectory: filepath.Join(filepath.Dir(exe), "templates"),
		},
	}

	err = viper.Unmarshal(&configData)
	if err != nil {
		return config.Config{}, fmt.Errorf("could not parse config file: %w\n", err)
	}

	return configData, nil
}

// absoluteOpenAPIFile uses the current working directory, resolved config file and the openAPI file that was specified
// in the config file to determine the absolute path to and Root file. It takes into account that any of these,
// except the working directory could be relative.
func absoluteOpenAPIFile(config config.Config) string {
	var absoluteOpenAPIFile string
	if filepath.IsAbs(config.OpenAPI.OpenAPIFile) {
		absoluteOpenAPIFile = config.OpenAPI.OpenAPIFile
	} else {
		if filepath.IsAbs(config.ConfigFilePath) {
			absoluteOpenAPIFile = filepath.Join(
				filepath.Dir(config.ConfigFilePath),
				config.OpenAPI.OpenAPIFile,
			)
		} else {
			absoluteOpenAPIFile = filepath.Join(
				config.WorkingDirectory,
				filepath.Dir(config.ConfigFilePath),
				config.OpenAPI.OpenAPIFile,
			)
		}
	}

	return absoluteOpenAPIFile
}
