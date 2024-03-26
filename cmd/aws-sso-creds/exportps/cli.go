package exportps

import (
	"fmt"

	"github.com/jaxxstorm/aws-sso-creds/pkg/credentials"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Command() *cobra.Command {
	command := &cobra.Command{
		Use:          "export-ps",
		Short:        "Generates a set of powershell environment assignments to define the AWS temporary creds to your environment",
		Long:         "Generates a set of powershell environment assignments to define the AWS temporary creds to your environment",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {

			cmd.SilenceUsage = true

			profile := viper.GetString("profile")
			homeDir := viper.GetString("home-directory")

			creds, _, err := credentials.GetSSOCredentials(profile, homeDir)

			if err != nil {
				return err
			}

			fmt.Printf("$env:AWS_ACCESS_KEY_ID='%s'\n", *creds.RoleCredentials.AccessKeyId)
			fmt.Printf("$env:AWS_SECRET_ACCESS_KEY='%s'\n", *creds.RoleCredentials.SecretAccessKey)
			fmt.Printf("$env:AWS_SESSION_TOKEN='%s'\n", *creds.RoleCredentials.SessionToken)

			return nil
		},
	}

	return command
}
