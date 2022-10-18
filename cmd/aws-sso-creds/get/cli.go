package get

import (
	"fmt"
	"os"
	"time"

	"github.com/jaxxstorm/aws-sso-creds/pkg/credentials"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	. "github.com/logrusorgru/aurora"
)


func Command() *cobra.Command {
	command := &cobra.Command{
		Use:   "get",
		Short: "Get AWS temporary credentials to use on the command line",
		Long:  "Retrieve AWS temporary credentials",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {

			cmd.SilenceUsage = true

			profile := viper.GetString("profile")
			homeDir := viper.GetString("home-directory")

			if profile == "" {
				return fmt.Errorf("no profile specified")
			}

			creds, accountID, err := credentials.GetSSOCredentials(profile, homeDir)

			if err != nil {
				return err
			}

			fmt.Println(Sprintf("Your temporary credentials for account %s are:", White(accountID)))
			fmt.Println("")

			fmt.Fprintln(os.Stdout, "AWS_ACCESS_KEY_ID\t", *creds.RoleCredentials.AccessKeyId)
			fmt.Fprintln(os.Stdout, "AWS_SECRET_ACCESS_KEY\t", *creds.RoleCredentials.SecretAccessKey)
			fmt.Fprintln(os.Stdout, "AWS_SESSION_TOKEN\t", *creds.RoleCredentials.SessionToken)

			fmt.Println("")

			fmt.Println("These credentials will expire at:", Red(time.UnixMilli(*creds.RoleCredentials.Expiration).UTC()))

			return nil
		},
	}

	return command
}
