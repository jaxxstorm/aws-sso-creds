package accounts

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
	results int32 // Adjusted to int32 as per v2 requirements
)

func Command() *cobra.Command {
	command := &cobra.Command{
		Use:   "accounts",
		Short: "List all accounts",
		Long:  "List all accounts",
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

			// Load default AWS config
			cfg, err := config.LoadDefaultConfig(context.TODO(),
				config.WithRegion(ssoConfig.Region),
				config.WithSharedConfigProfile(profile),
			)
			if err != nil {
				return fmt.Errorf("error loading AWS config: %v", err)
			}

			svc := sso.NewFromConfig(cfg)

			accounts, err := svc.ListAccounts(context.TODO(), &sso.ListAccountsInput{
				AccessToken: &token,
				MaxResults:  &results, // Note: MaxResults might need type adjustment
			})
			if err != nil {
				return fmt.Errorf("error listing accounts: %v", err)
			}

			writer := tabwriter.NewWriter(os.Stdout, tabwriterMinWidth, tabwriterWidth, tabwriterPadding, tabwriterPadChar, 0)
			fmt.Fprintln(writer, "ID\tNAME\tEMAIL ADDRESS")

			for _, account := range accounts.AccountList {
				fmt.Fprintf(writer, "%s\t%s\t%s\n", *account.AccountId, *account.AccountName, *account.EmailAddress)
			}

			writer.Flush()

			return nil
		},
	}

	command.Flags().Int32VarP(&results, "results", "r", 10, "Maximum number of accounts to return")

	return command
}
