package roles

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sso"
	cfg "github.com/jaxxstorm/aws-sso-creds/pkg/config"
	"github.com/jaxxstorm/aws-sso-creds/pkg/contract"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	tabwriterMinWidth = 6
	tabwriterWidth    = 4
	tabwriterPadding  = 3
	tabwriterPadChar  = ' '
)

var (
	results   int32 // Adjusted to int32 as per v2 requirements
	accountID string
)

type listAccountRolesAPI interface {
	ListAccountRoles(context.Context, *sso.ListAccountRolesInput, ...func(*sso.Options)) (*sso.ListAccountRolesOutput, error)
}

var (
	loadDefaultConfig = config.LoadDefaultConfig
	newSSOClient      = func(awsCfg aws.Config) listAccountRolesAPI {
		return sso.NewFromConfig(awsCfg)
	}
)

func Command() *cobra.Command {
	command := &cobra.Command{
		Use:   "roles ACCOUNT_ID",
		Short: "List roles for an account",
		Long:  "List all the roles for an account",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			cmd.SilenceUsage = true

			profile := viper.GetString("profile")
			homeDir := viper.GetString("home-directory")

			ssoConfig, err := cfg.GetSSOConfig(profile, homeDir)
			if err != nil {
				return fmt.Errorf("error retrieving SSO config: %w", err)
			}

			token, err := cfg.GetSSOToken(*ssoConfig, homeDir)
			if err != nil {
				return fmt.Errorf("error retrieving SSO token from cache files: %v", err)
			}

			awsCfg, err := loadDefaultConfig(context.TODO(),
				config.WithRegion(ssoConfig.Region),
				config.WithSharedConfigProfile(profile),
			)
			if err != nil {
				return fmt.Errorf("error loading AWS config: %v", err)
			}

			svc := newSSOClient(awsCfg)

			accountID = args[0]

			roles, err := svc.ListAccountRoles(context.TODO(), &sso.ListAccountRolesInput{
				AccessToken: &token,
				MaxResults:  &results, // Note: MaxResults might need type adjustment
				AccountId:   &accountID,
			})
			if err != nil {
				return fmt.Errorf("error listing roles: %v", err)
			}

			writer := tabwriter.NewWriter(os.Stdout, tabwriterMinWidth, tabwriterWidth, tabwriterPadding, tabwriterPadChar, 0)
			contract.IgnoreIoError(fmt.Fprintln(writer, "ROLE NAME"))

			for _, role := range roles.RoleList {
				contract.IgnoreIoError(fmt.Fprintf(writer, "%s\n", *role.RoleName))
			}

			if err := writer.Flush(); err != nil {
				return err
			}

			return nil
		},
	}

	command.Flags().Int32VarP(&results, "results", "r", 10, "Maximum number of roles to return") // Adjusted to Int32VarP

	return command
}
