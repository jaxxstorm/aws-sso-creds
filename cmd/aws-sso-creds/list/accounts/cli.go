package accounts

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sso"
	"github.com/liggitt/tabwriter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/thuannfq/aws-sso-creds/pkg/config"
)

const (
	tabwriterMinWidth = 6
	tabwriterWidth    = 4
	tabwriterPadding  = 3
	tabwriterPadChar  = ' '
	tabwriterFlags    = tabwriter.RememberWidths
)

var (
	results int64
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

			ssoConfig, err := config.GetSSOConfig(profile, homeDir)
			if err != nil {
				return fmt.Errorf("error retrieving SSO config: %w", err)
			}

			cacheFiles, err := os.ReadDir(filepath.Join(homeDir, ".aws", "sso", "cache"))
			if err != nil {
				return fmt.Errorf("error retrieving cache files - perhaps you need to login?: %w", err)
			}

			token, err := config.GetSSOToken(cacheFiles, *ssoConfig, homeDir)
			if err != nil {
				return fmt.Errorf("error retrieving SSO token from cache files: %v", err)
			}

			sess := session.Must(session.NewSession())
			svc := sso.New(sess, aws.NewConfig().WithRegion(ssoConfig.Region))

			accounts, err := svc.ListAccounts(&sso.ListAccountsInput{
				AccessToken: &token,
				MaxResults:  &results,
			})
			if err != nil {
				return fmt.Errorf("error listing accounts: %v", err)
			}

			writer := tabwriter.NewWriter(os.Stdout, tabwriterMinWidth, tabwriterWidth, tabwriterPadding, tabwriterPadChar, tabwriterFlags)
			fmt.Fprintln(writer, "ID\tNAME\tEMAIL ADDRESS")

			for _, results := range accounts.AccountList {
				fmt.Fprintf(writer, "%s\t%s\t%s\n", *results.AccountId, *results.AccountName, *results.EmailAddress)
			}

			writer.Flush()

			return nil
		},
	}

	command.Flags().Int64VarP(&results, "results", "r", 10, "Maximum number of accounts to return")

	return command
}
