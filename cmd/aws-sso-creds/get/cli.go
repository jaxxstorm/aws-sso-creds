package get

import (
	"fmt"

	"github.com/jaxxstorm/aws-sso-creds/pkg/credentials"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Command() *cobra.Command {
	command := &cobra.Command{
		Use:   "get",
		Short: "Get AWS temporary credentials to use on the command line",
		Long:  "Retrieve AWS temporary credentials",
		RunE: func(cmd *cobra.Command, args []string) error {

			profile := viper.GetString("profile")
			homeDir := viper.GetString("home-directory")

			if profile == "" {
				return fmt.Errorf("no profile specified")
			}

			creds, err := credentials.GetSSOCredentials(profile, homeDir)

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
