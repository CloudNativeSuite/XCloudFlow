package modules

import (
	"context"
	"fmt"
	"os"

	"craftweave/core/parser"
	"craftweave/internal/ssh"
	"github.com/vultr/govultr/v3"
	"golang.org/x/oauth2"
)

func vultrHandler(ctx Context, task parser.Task) ssh.CommandResult {
	if task.Vultr == nil {
		return ssh.CommandResult{Host: ctx.Host.Name, ReturnMsg: "FAILED", ReturnCode: 1, Output: "missing vultr parameters"}
	}

	apiKey := task.Vultr.APIKey
	if apiKey == "" {
		apiKey = os.Getenv("VULTR_API_KEY")
	}
	if apiKey == "" {
		return ssh.CommandResult{Host: ctx.Host.Name, ReturnMsg: "FAILED", ReturnCode: 1, Output: "missing api key"}
	}

	config := &oauth2.Config{}
	ctxAPI := context.Background()
	ts := config.TokenSource(ctxAPI, &oauth2.Token{AccessToken: apiKey})
	client := govultr.NewClient(oauth2.NewClient(ctxAPI, ts))

	req := &govultr.InstanceCreateReq{
		Region: task.Vultr.Region,
		Plan:   task.Vultr.Plan,
		OsID:   task.Vultr.OsID,
		Label:  task.Vultr.Label,
	}

	inst, _, err := client.Instance.Create(ctxAPI, req)
	if err != nil {
		return ssh.CommandResult{Host: ctx.Host.Name, ReturnMsg: "FAILED", ReturnCode: 1, Output: err.Error()}
	}

	return ssh.CommandResult{Host: ctx.Host.Name, ReturnMsg: "CHANGED", ReturnCode: 0, Output: fmt.Sprintf("ID:%s IP:%s", inst.ID, inst.MainIP)}
}

func init() { Register("vultr_instance", vultrHandler) }
