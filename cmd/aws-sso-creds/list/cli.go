package list

import (
	"github.com/jaxxstorm/aws-sso-creds/cmd/aws-sso-creds/list/accounts"
	"github.com/jaxxstorm/aws-sso-creds/cmd/aws-sso-creds/list/roles"
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	command := &cobra.Command{
		Use:   "list",
		Short: "list commands",
		Long:  "Commands that list things",
	}

	command.AddCommand(accounts.Command())
	command.AddCommand(roles.Command())

	return command
}
