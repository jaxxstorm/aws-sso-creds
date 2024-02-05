package set

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/jaxxstorm/aws-sso-creds/pkg/credentials"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/bigkevmcd/go-configparser"
)

var (
	credsFile *configparser.ConfigParser
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

			credsPath := filepath.Join(homeDir, ".aws", "credentials")
			cfgPath := filepath.Join(homeDir, ".aws", "config")

			fmt.Println(credsPath)
			fmt.Println(cfgPath)

			creds, _, err := credentials.GetSSOCredentials(profile, homeDir)
			if err != nil {
				return err
			}

			credsFile, err = configparser.NewConfigParserFromFile(credsPath)
			if err != nil {
				if os.IsNotExist(err) {
					// Ensure the new empty credentials file is not readable by others.
					if f, err := os.OpenFile(credsPath, os.O_CREATE, 0600); err != nil {
						f.Close()
						return err
					}

					credsFile = configparser.New()
				}
				return fmt.Errorf("error parsing config file: %v", err)
			}

			configFile, err := configparser.NewConfigParserFromFile(cfgPath)
			if err != nil {
				return err
			}

			// create a new credentials section
			if err := credsFile.AddSection(args[0]); err != nil {
				return fmt.Errorf("error creating credentials section in creds file: %v", err)
			}

			if err := configFile.AddSection(fmt.Sprintf("profile %s", args[0])); err != nil {
				return fmt.Errorf("error creating credentials section in config file: %v", err)
			}

			if err := credsFile.Set(args[0], "aws_access_key_id", *creds.RoleCredentials.AccessKeyId); err != nil {
				return fmt.Errorf("error setting access key id: %v", err)
			}
			if err := credsFile.Set(args[0], "aws_secret_access_key", *creds.RoleCredentials.SecretAccessKey); err != nil {
				return fmt.Errorf("error setting secret access key: %v", err)
			}
			if err := credsFile.Set(args[0], "aws_session_token", *creds.RoleCredentials.SessionToken); err != nil {
				return fmt.Errorf("error setting session token: %v", err)
			}

			if err := credsFile.SaveWithDelimiter(credsPath, "="); err != nil {
				return fmt.Errorf("error saving credentials file: %v", err)
			}

			if err := configFile.SaveWithDelimiter(cfgPath, "="); err != nil {
				return fmt.Errorf("error saving config file: %v", err)
			}

			fmt.Printf("credentials saved to profile: %s\n", args[0])
			fmt.Printf("these credentials will expire:  %s\n", time.Unix(creds.RoleCredentials.Expiration, 0).Format(time.UnixDate))

			return nil
		},
	}

	return command
}
