package set

import (
	"fmt"
	"os"
	"time"

	"github.com/jaxxstorm/aws-sso-creds/pkg/credentials"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/bigkevmcd/go-configparser"
)

func Command() *cobra.Command {
	command := &cobra.Command{
		Use:   "set PROFILE",
		Short: "Create a new AWS profile with temporary credentials",
		Long:  "Create a new AWS profile with temporary credentials",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			cmd.SilenceUsage = true
			profile := viper.GetString("profile")
			homeDir := viper.GetString("home-directory")
			credsPath := fmt.Sprintf("%s/.aws/credentials", homeDir)
			cfgPath := fmt.Sprintf("%s/.aws/config", homeDir)

			if profile == "" {
				return fmt.Errorf("no profile specified")
			}

			creds, _, err := credentials.GetSSOCredentials(profile, homeDir)
			if err != nil {
				return err
			}

			credsFile, err := configparser.NewConfigParserFromFile(credsPath)
			if os.IsNotExist(err) {
				// Ensure the new empty credentials file is not readable by others.
				if f, err := os.OpenFile(credsPath, os.O_CREATE, 0600); err != nil {
					return err
				} else {
					f.Close()
				}
				credsFile = configparser.New()
			} else if err != nil {
				return err
			}

			configFile, err := configparser.NewConfigParserFromFile(cfgPath)
			if err != nil {
				return err
			}

			// create a new credentials section
			credsFile.AddSection(args[0])
			configFile.AddSection(fmt.Sprintf("profile %s", args[0]))

			credsFile.Set(args[0], "aws_access_key_id", *creds.RoleCredentials.AccessKeyId)
			credsFile.Set(args[0], "aws_secret_access_key", *creds.RoleCredentials.SecretAccessKey)
			credsFile.Set(args[0], "aws_session_token", *creds.RoleCredentials.SessionToken)

			credsFile.SaveWithDelimiter(credsPath, "=")
			configFile.SaveWithDelimiter(cfgPath, "=")

			fmt.Printf("credentials saved to profile: %s\n", args[0])
			fmt.Printf("these credentials will expire:  %s\n", time.Unix(*creds.RoleCredentials.Expiration, 0).Format(time.UnixDate))

			return nil
		},
	}

	return command
}
