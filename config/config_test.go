package config

import (
	"os"
	"reflect"
	"testing"
)

func TestNewConfig(t *testing.T) {
	type args struct {
		configStruct any
		cfgOptions   *ConfigOptions
	}
	tests := []struct {
		name    string
		args    args
		want    any
		wantErr bool
	}{
		{
			name: "successful_config_with_env_vars",
			args: args{
				configStruct: &struct {
					JiraPassword string `env:"TEST_1234_JIRA_PASSWORD"`
					JiraUsername string `env:"TEST_1234_JIRA_USERNAME"`
				}{},
			},
			want: &struct {
				JiraPassword string `env:"TEST_1234_JIRA_PASSWORD"`
				JiraUsername string `env:"TEST_1234_JIRA_USERNAME"`
			}{
				JiraPassword: "password",
				JiraUsername: "username",
			},
			wantErr: false,
		},
	}

	os.Setenv("TEST_1234_JIRA_PASSWORD", "password")
	os.Setenv("TEST_1234_JIRA_USERNAME", "username")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewConfig(tt.args.configStruct, tt.args.cfgOptions)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
