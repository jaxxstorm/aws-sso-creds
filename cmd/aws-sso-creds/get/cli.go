package get

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sso"
	"github.com/jaxxstorm/aws-sso-creds/pkg/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
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

			ssoConfig, err := config.GetSSOConfig(profile)
			if err != nil {
				fmt.Errorf("error retrieving SSO config: %w", err)
			}

			cacheFiles, err := ioutil.ReadDir(fmt.Sprintf("%s/.aws/sso/cache", homeDir))
			if err != nil {
				fmt.Errorf("error retrieving cache files - perhaps you need to login?: %w", err)
			}

			token, err := config.GetSSOToken(cacheFiles, *ssoConfig)
			if err != nil {
				panic(err)
			}

			sess := session.Must(session.NewSession())
			svc := sso.New(sess, aws.NewConfig().WithRegion(ssoConfig.Region))

			creds, err := svc.GetRoleCredentials(&sso.GetRoleCredentialsInput{
				AccessToken: &token,
				AccountId: &ssoConfig.AccountID,
				RoleName: &ssoConfig.RoleName,
			})

			if err != nil {
				panic(err)
			}

			fmt.Println(creds)

			return nil
		},
	}


	return command
}
