package examples

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/mcsteele8/common-cli-utils/config"
)

type cliConfig struct {
	JiraUsername string `env:"CLI_JIRA_USERNAME" file:"jira_username" default:"empty"`
	JiraPassword string `env:"CLI_JIRA_PASSWORD" file:"jira_password" default:"empty" mask:"true"`
}

func exampleConfig() {

	cfg := &cliConfig{}

	cfgResult, err := config.NewConfig(cfg, nil)
	if err != nil {
		log.Fatalf("failed to set config values: %v", err)
	}
	cfg = cfgResult.(*cliConfig)

	fmt.Println("Jira Username:", cfg.JiraUsername)
	fmt.Println("Jira Password:", cfg.JiraPassword)

	os.Setenv("CLI_JIRA_USERNAME", "testuser")
	os.Setenv("CLI_JIRA_PASSWORD", "testpass")

	cfgResult, err = config.NewConfig(cfg, nil)
	if err != nil {
		log.Fatalf("failed to set config values: %v", err)
	}
	cfg = cfgResult.(*cliConfig)

	fmt.Println("Jira Username:", cfg.JiraUsername)
	fmt.Println("Jira Password:", cfg.JiraPassword)
}

func exampleConfigWithCtx() {
	ctx := context.Background()

	cfg := &cliConfig{}
	os.Setenv("CLI_JIRA_USERNAME", "testuser1")
	os.Setenv("CLI_JIRA_PASSWORD", "testpass1")

	ctx, _, err := config.NewCtxWithConfig(ctx, cfg, nil)
	if err != nil {
		log.Fatalf("failed to set config values: %v", err)
	}

	cfg = config.FromCtx(ctx).(*cliConfig)

	fmt.Println("Jira Username From Ctx:", cfg.JiraUsername)
	fmt.Println("Jira Password From Ctx:", cfg.JiraPassword)
}
