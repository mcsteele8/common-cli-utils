package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

/*
cliConfig struct is what is used to manage all configs for remi.
each struct value can be given multiple tags for reading config values from multiple areas.

env:      Is the tag used to pull environment variables during run time. Env tag value will hold priority over other tags
file:     Is the tag used to pull file values stored in $HOME/.remi-cli/cfg.yaml
default:  Is the tag that will be used if no env or file value can be found
mask:     Is the tag to mask the output of the value

Example:
type cliConfig struct {
	JiraUsername string `env:"REMI_CLI_JIRA_USERNAME" file:"jira_username" default:"empty"`
	JiraPassword string `env:"REMI_CLI_JIRA_PASSWORD" file:"jira_password" default:"empty" mask:"true"`
}
*/

type ctxKey string

const (
	cfgCtxKey ctxKey = "cfg-ctx-key"

	cfgTagEnv     = "env"
	cfgTagFile    = "file"
	cfgTagDefault = "default"
	cfgTagMask    = "mask"
)

type ConfigOptions struct {
	CfgDirectory             string
	CfgFilePath              string
	CfgFileName              string
	CfgFileType              string
	CreateEmptyCfgIfNotFound bool
}

var defaultCfgOptions = ConfigOptions{
	CfgDirectory: "",
	CfgFilePath:  "",
	CfgFileName:  "config",
	CfgFileType:  "yaml",
}

func NewConfig(configStruct any, cfgOptions *ConfigOptions) (any, error) {
	if cfgOptions == nil {
		cfgOptions = &defaultCfgOptions
	}

	// Check if configStruct is a pointer to a struct
	val := reflect.ValueOf(configStruct)
	if val.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("configStruct must be a pointer to a struct, got %v", val.Kind())
	}

	if val.Elem().Kind() != reflect.Struct {
		return nil, fmt.Errorf("configStruct must be a pointer to a struct, got a pointer to %v", val.Elem().Kind())
	}

	if cfgOptions.CfgFilePath != "" {
		viper.SetConfigFile(cfgOptions.CfgFilePath)
	} else {
		viper.SetConfigName(cfgOptions.CfgFileName)
		viper.SetConfigType(cfgOptions.CfgFileType)
		viper.AddConfigPath(cfgOptions.CfgDirectory)
	}

	cfgFileFound := true
	if err := viper.ReadInConfig(); err != nil && strings.Contains(err.Error(), "Not Found") {
		cfgFileFound = false
	}

	readStruct(val.Elem(), false)
	if !cfgFileFound && cfgOptions.CreateEmptyCfgIfNotFound {
		if err := initEmptyCfg(cfgOptions); err != nil {
			return nil, fmt.Errorf("failed to init empty config: %w", err)
		}
	}

	return configStruct, nil
}

func NewCtxWithConfig(ctx context.Context, configStruct any, cfgOptions *ConfigOptions) (context.Context, any, error) {
	config, err := NewConfig(configStruct, cfgOptions)
	if err != nil {
		return ctx, config, err
	}

	// set the config in context to easily pass to functions
	ctx = context.WithValue(ctx, cfgCtxKey, configStruct)

	return ctx, config, nil
}

func FromCtx(ctx context.Context) any {
	return ctx.Value(cfgCtxKey)
}

// func init() {
// 	var err error
// 	home_dir, err = os.UserHomeDir()
// 	if err != nil {
// 		log.Fatalln("failed to get user home directory: ", err.Error())
// 	}

// 	cli_cfg_dir = path.Join(home_dir, ".remi-cli")
// 	cli_cfg_path = path.Join(cli_cfg_dir, cli_cfg_file_name)
// }

func initEmptyCfg(cfgOptions *ConfigOptions) error {
	// create an empty config file with -rwxrwxrwx	0777  read, write, & execute for owner, group and others permissions
	os.Mkdir(cfgOptions.CfgDirectory, 0777)
	err := os.WriteFile(cfgOptions.CfgFilePath, []byte(""), 0777)
	if err != nil {
		return fmt.Errorf("failed to create an empty cfg file %w", err)
	}

	err = viper.WriteConfig()
	if err != nil {
		return fmt.Errorf("failed to write values to new cfg %w", err)
	}

	return nil
}

// readStruct is used to read the struct and will be recursively called
// to read all child structs within cfg
func readStruct(input reflect.Value, verbose bool) {
	inputType := input.Type()

	for i := 0; i < input.NumField(); i++ {
		fieldValue := input.Field(i)
		fieldName := inputType.Field(i).Name

		switch fieldValue.Kind() {
		case reflect.Struct:
			readStruct(fieldValue, verbose)
		case reflect.String:
			setString(fieldValue, inputType.Field(i).Tag)
		case reflect.Bool:
			setBool(fieldValue, inputType.Field(i).Tag)
		case reflect.Int:
			setInt(fieldValue, inputType.Field(i).Tag)
		default:
			log.Fatalf("Config type not supported yet: %s\n", fieldValue.Kind().String())
		}

		if verbose && fieldValue.Kind() != reflect.Struct {
			fmt.Printf("%s: %v\n", fieldName, getOutputValue(fieldValue, inputType.Field(i).Tag))
		}

	}

}

func getOutputValue(fieldValue reflect.Value, tag reflect.StructTag) interface{} {
	if tag.Get(cfgTagMask) == "true" {
		return "*********"
	}
	return fieldValue
}

func getTagValue(tag reflect.StructTag) string {
	envTag := tag.Get(cfgTagEnv)
	value := os.Getenv(envTag)
	if value == "" {
		value = viper.GetString(tag.Get(cfgTagFile))
	}

	if value == "" {
		value = tag.Get(cfgTagDefault)
	}
	return value
}

func setString(fieldValue reflect.Value, tag reflect.StructTag) {
	value := getTagValue(tag)

	// Ensure the value is addressable
	if fieldValue.CanSet() {
		// Set the field value
		fieldValue.SetString(value)
		viper.Set(tag.Get(cfgTagFile), value)
	}
}

func setInt(fieldValue reflect.Value, tag reflect.StructTag) {
	value := getTagValue(tag)

	intValue, err := strconv.Atoi(value)
	if err != nil {
		log.Fatalln("failed to set int config value: ", err.Error())
	}

	// Ensure the value is addressable
	if fieldValue.CanSet() {
		// Set the field value
		fieldValue.SetInt(int64(intValue))
		viper.Set(tag.Get(cfgTagFile), value)
	}
}

func setBool(fieldValue reflect.Value, tag reflect.StructTag) {
	value := getTagValue(tag)

	boolValue := false
	if value == "true" {
		boolValue = true
	}

	// Ensure the value is addressable
	if fieldValue.CanSet() {
		// Set the field value
		fieldValue.SetBool(boolValue)
		viper.Set(tag.Get(cfgTagFile), value)
	}
}
