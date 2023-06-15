package export

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/thuannfq/aws-sso-creds/pkg/credentials"
)

func Command() *cobra.Command {
	command := &cobra.Command{
		Use:          "export",
		Short:        "Generates a set of shell commands to export AWS temporary creds to your environment",
		Long:         "Generates a set of shell commands to export AWS temporary creds to your environment",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {

			cmd.SilenceUsage = true

			profile := viper.GetString("profile")
			homeDir := viper.GetString("home-directory")

			creds, _, err := credentials.GetSSOCredentials(profile, homeDir)

			if err != nil {
				return err
			}

			fmt.Printf("export AWS_ACCESS_KEY_ID=%s\n", *creds.RoleCredentials.AccessKeyId)
			fmt.Printf("export AWS_SECRET_ACCESS_KEY=%s\n", *creds.RoleCredentials.SecretAccessKey)
			fmt.Printf("export AWS_SESSION_TOKEN=%s\n", *creds.RoleCredentials.SessionToken)

			return nil
		},
	}

	return command
}
