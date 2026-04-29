package operations

import (
	"context"
	"fmt"

	"github.com/qiniu/go-sdk/v7/sandbox"

	sbClient "github.com/qiniu/qshell/v2/iqshell/sandbox"
)

// PublishInfo holds parameters for publishing/unpublishing templates.
type PublishInfo struct {
	TemplateIDs []string // One or more template IDs
	Yes         bool     // Skip confirmation
	Public      bool     // true = publish, false = unpublish
}

// Publish publishes or unpublishes one or more templates.
func Publish(info PublishInfo) {
	if len(info.TemplateIDs) == 0 {
		id, ok := templateIDFromCwdConfig()
		if !ok {
			return
		}
		if id != "" {
			info.TemplateIDs = []string{id}
		}
	}
	if len(info.TemplateIDs) == 0 {
		sbClient.PrintError("at least one template ID is required (positional args or qshell.sandbox.toml)")
		return
	}

	client, err := sbClient.NewSandboxClient()
	if err != nil {
		sbClient.PrintError("%v", err)
		return
	}

	ctx := context.Background()

	action := "publish"
	if !info.Public {
		action = "unpublish"
	}

	if !info.Yes {
		if !sbClient.IsInteractive() {
			sbClient.PrintError("confirmation required but stdin is not a terminal; pass --yes to confirm in non-interactive mode")
			return
		}
		fmt.Printf("Are you sure you want to %s %d template(s)? [y/N] ", action, len(info.TemplateIDs))
		var confirm string
		fmt.Scanln(&confirm)
		if confirm != "y" && confirm != "Y" {
			fmt.Println("Aborted")
			return
		}
	}

	for _, id := range info.TemplateIDs {
		if uErr := client.UpdateTemplate(ctx, id, sandbox.UpdateTemplateParams{
			Public: &info.Public,
		}); uErr != nil {
			sbClient.PrintError("%s template %s failed: %v", action, id, uErr)
			continue
		}
		if info.Public {
			sbClient.PrintSuccess("Template %s published", id)
		} else {
			sbClient.PrintSuccess("Template %s unpublished", id)
		}
	}
}
