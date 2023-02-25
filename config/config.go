package config

import (
	"flag"
	"fmt"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

// SemVer should be updated on any new release!!
const SemVer = "0.0.6"

// r represents System level settings for flags and environmental context tracking
var r *System

type System struct {
	Output, Debug, Dryrun, VersionFlagSet *bool
	// ConfigPath is `.` if not set through the `-configFile` flag
	ConfigPath *string

	// ProcWorkingDir is set through calling os.Getwd()
	ExePath, ProcWorkingDir string
}

func Runtime() *System {
	return r
}

// Data represents the config data used in the day-to-day running of can
//
//	TODO this may be vague in definition for the sake of its legibility in use
//	TODO check for redundancy
type Data struct {
	Generator
	Template

	// left public due to it's need to be unmarshalled by Data.Load()
	TemplatesDir string `yaml:"templatesDir"`

	// OpenAPIFile represents the path to the yaml OpenAPI 3 file to render
	OpenAPIFile    string
	absOpenAPIPath string

	OutputPath    string
	absOutputPath string
}

type Generator struct {
	ModuleName string

	// BasePackageName represents the
	BasePackageName string
}

// Template is populated based on it's Name variable set as a CLI flag
type Template struct {
	Name *string

	// Directory should be ./templates/${Name} by default
	Directory    string
	absDirectory string
}

func (d *Data) Load() (err error) {
	r = &System{}
	// Setup config pre-unmarshalling. This assumes we don't change directory before this is executed
	r.ProcWorkingDir, err = os.Getwd()
	if err != nil {
		return fmt.Errorf("r.Config.load:: could not determine working directory: %w\n", err)
	}
	r.ExePath, err = os.Executable()
	if err != nil {
		return fmt.Errorf("Config.load:: could not determine executable directory: %w\n", err)
	}

	// flags
	//var VersionFlagSet *bool
	r.VersionFlagSet = flag.Bool("version", false, "Print Can version and exit")
	r.Debug = flag.Bool("debug", false, "Enable debug logging")
	r.ConfigPath = flag.String("configFile", ".", "Specify which config file to use")
	r.Dryrun = flag.Bool("dryrun", false,
		"Toggles whether to perform a render without writing to disk."+
			"This works particularly well in combination with -debug")
	d.Template.Name = flag.String("template", "", "Specify which template set to use")
	flag.Parse()

	absCfgPath, err := filepath.Abs(*r.ConfigPath)
	if err != nil {
		return fmt.Errorf("could not resolve relative config path: %w", err)
	}
	r.ConfigPath = &absCfgPath

	if *r.VersionFlagSet {
		fmt.Printf("Can Version: %s\n", SemVer)
		os.Exit(0)
	}

	// config load
	if *r.Debug {
		fmt.Printf("[v%s]::Using config file \"%s\".\n", SemVer, *r.ConfigPath)
	}
	viper.SetConfigFile(*r.ConfigPath)

	err = viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("loadConfig:: could not read config file: %w\n", err)
	}

	err = viper.Unmarshal(&d)
	if err != nil {
		return fmt.Errorf("loadConfig:: could not parse config file: %w\n", err)
	}

	// This should always happen at the end of this function
	// Handle Templates
	if *d.Template.Name == "" {
		fmt.Printf("template is a required flag\nexiting...")
		os.Exit(1)
	}
	if !d.validTemplateName() {
		fmt.Printf("%s does not exist in %s\nexiting...\n", *d.Template.Name, d.TemplatesDir)
		fmt.Println("Valid template names are:")
		names, err := d.validTemplates()
		if err != nil {
			return fmt.Errorf("could not read templates: %w", err)
		}
		for _, name := range names {
			fmt.Println(name)
		}
		os.Exit(1)
	}

	// resolve paths
	err = d.resolveTemplateConfig()
	if err != nil {
		return err
	}

	return nil
}

func (d *Data) GetTemplateDir() (path string) {
	if d.Template.absDirectory != "" {
		return d.Template.absDirectory
	}
	switch true {
	case filepath.IsAbs(d.TemplatesDir):
		d.Template.absDirectory = filepath.Join(d.TemplatesDir, *d.Template.Name)
		return d.Template.absDirectory
	case filepath.IsAbs(*r.ConfigPath):
		d.Template.absDirectory = filepath.Join(filepath.Dir(*r.ConfigPath), d.TemplatesDir, *d.Template.Name)
		return d.Template.absDirectory
	default:
		d.Template.absDirectory = filepath.Join(r.ProcWorkingDir, filepath.Dir(*r.ConfigPath), d.TemplatesDir, *d.Template.Name)
		return d.Template.absDirectory
	}
}

// GetOutputDir is used by the render engine to determine where rendered files will be written to
func (d *Data) GetOutputDir() (path string) {
	if d.absOutputPath != "" {
		return d.absOutputPath
	}
	switch true {
	case filepath.IsAbs(d.OutputPath):
		d.absOutputPath = d.OutputPath
		return d.absOutputPath
	case filepath.IsAbs(*r.ConfigPath):
		d.absOutputPath = filepath.Join(
			filepath.Dir(*r.ConfigPath),
			d.OutputPath,
		)
		return d.absOutputPath
	default:
		d.absOutputPath = filepath.Join(
			r.ProcWorkingDir,
			filepath.Dir(*r.ConfigPath),
			d.OutputPath,
		)
		return d.absOutputPath
	}
}

// GetOpenAPIFilepath uses the current working directory, resolved config file and the openAPI file that was specified
// in the config file to determine the absolute path to an OpenAPI file. It takes into account that any of these,
// except the working directory could be relative. It returns the absolute value on every call by caching the result of
// it's first run and returning that on successive calls
func (d *Data) GetOpenAPIFilepath() (path string) {
	if d.absOpenAPIPath != "" { // we shouldn't have to run below logic multiple times
		return d.absOpenAPIPath
	}
	if filepath.IsAbs(d.OpenAPIFile) {
		d.absOpenAPIPath = d.OpenAPIFile
		return d.absOpenAPIPath
	} else {
		if filepath.IsAbs(*r.ConfigPath) {
			d.absOpenAPIPath = filepath.Join(
				filepath.Dir(*r.ConfigPath),
				d.OpenAPIFile,
			)
			return d.absOpenAPIPath
		} else {
			d.absOpenAPIPath = filepath.Join(
				// TODO test this
				// not relative as per above comment
				r.ProcWorkingDir,
				filepath.Dir(*r.ConfigPath),
				d.OpenAPIFile,
			)
			return d.absOpenAPIPath
		}
	}
}

func (d *Data) validTemplateName() bool {
	dirs, err := d.validTemplates()
	if err != nil {
		fmt.Println(fmt.Errorf("could not list valid templates in %s :: %w", d.TemplatesDir, err))
		return false
	}
	for _, dir := range dirs {
		if dir == *d.Template.Name {
			return true
		}
	}
	return false
}

func (d *Data) validTemplates() (templates []string, err error) {
	dirs, err := os.ReadDir(d.TemplatesDir)
	if err != nil {
		fmt.Println(fmt.Errorf("could not list directories %w", err))
		return nil, err
	}
	for _, dir := range dirs {
		templates = append(templates, dir.Name())
	}
	return templates, nil
}
func (d *Data) resolveTemplateConfig() error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	if d.TemplatesDir == "" {
		d.TemplatesDir = filepath.Join(filepath.Dir(exe), "templates")
	}
	if d.Template.Directory == "" {
		d.Template.Directory = "./templates/" + *d.Template.Name
	}
	return nil
}
