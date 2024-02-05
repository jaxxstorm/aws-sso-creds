package roles

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"text/tabwriter"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sso"
	cfg "github.com/jaxxstorm/aws-sso-creds/pkg/config"
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

			cacheFiles, err := os.ReadDir(filepath.Join(homeDir, ".aws", "sso", "cache"))
			if err != nil {
				return fmt.Errorf("error retrieving cache files - perhaps you need to login?: %w", err)
			}

			token, err := cfg.GetSSOToken(cacheFiles, *ssoConfig, homeDir)
			if err != nil {
				return fmt.Errorf("error retrieving SSO token from cache files: %v", err)
			}

			cfg, err := config.LoadDefaultConfig(context.TODO(),
				config.WithRegion(ssoConfig.Region),
				config.WithSharedConfigProfile(profile),
			)
			if err != nil {
				return fmt.Errorf("error loading AWS config: %v", err)
			}

			svc := sso.NewFromConfig(cfg)

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
			fmt.Fprintln(writer, "ROLE NAME")

			for _, role := range roles.RoleList {
				fmt.Fprintf(writer, "%s\t%s\n", *role.RoleName, *role.RoleName)
			}

			writer.Flush()

			return nil
		},
	}

	command.Flags().Int32VarP(&results, "results", "r", 10, "Maximum number of roles to return") // Adjusted to Int32VarP

	return command
}
