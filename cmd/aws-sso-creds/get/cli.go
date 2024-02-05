package get

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/jaxxstorm/aws-sso-creds/pkg/credentials"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/logrusorgru/aurora"
)

type JSON struct {
	AwsAccessKeyID     string    `json:"aws_access_key_id"`
	AwsSecretAccessKey string    `json:"aws_secret_access_key"`
	SessionToken       string    `json:"aws_session_token"`
	ExpireAt           time.Time `json:"expire_at"`
}

func Command() *cobra.Command {
	command := &cobra.Command{
		Use:          "get",
		Short:        "Get AWS temporary credentials to use on the command line",
		Long:         "Retrieve AWS temporary credentials",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {

			cmd.SilenceUsage = true

			profile := viper.GetString("profile")
			homeDir := viper.GetString("home-directory")
			exportJSON, _ := cmd.Flags().GetBool("json")

			creds, accountID, err := credentials.GetSSOCredentials(profile, homeDir)

			if err != nil {
				return err
			}

			if exportJSON {
				credJSON := JSON{
					AwsAccessKeyID:     *creds.RoleCredentials.AccessKeyId,
					AwsSecretAccessKey: *creds.RoleCredentials.SecretAccessKey,
					SessionToken:       *creds.RoleCredentials.SessionToken,
					ExpireAt:           time.UnixMilli(creds.RoleCredentials.Expiration),
				}
				output, err := json.Marshal(credJSON)
				if err != nil {
					return err
				}
				fmt.Println(string(output))
			} else {

				fmt.Println(aurora.Sprintf("Your temporary credentials for account %s are:", aurora.White(accountID)))
				fmt.Println("")

				fmt.Fprintln(os.Stdout, "AWS_ACCESS_KEY_ID\t", *creds.RoleCredentials.AccessKeyId)
				fmt.Fprintln(os.Stdout, "AWS_SECRET_ACCESS_KEY\t", *creds.RoleCredentials.SecretAccessKey)
				fmt.Fprintln(os.Stdout, "AWS_SESSION_TOKEN\t", *creds.RoleCredentials.SessionToken)

				fmt.Println("")

				fmt.Println("These credentials will expire at:", aurora.Red(time.UnixMilli(creds.RoleCredentials.Expiration)))
			}

			return nil
		},
	}
	command.PersistentFlags().BoolP("json", "j", false, "print output in json format")
	return command
}
