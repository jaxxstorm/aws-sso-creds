package roles

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sso"
	"github.com/jaxxstorm/aws-sso-creds/pkg/config"
	"github.com/liggitt/tabwriter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	tabwriterMinWidth = 6
	tabwriterWidth    = 4
	tabwriterPadding  = 3
	tabwriterPadChar  = ' '
	tabwriterFlags    = tabwriter.RememberWidths
)

var (
	results   int64
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

			if profile == "" {
				return fmt.Errorf("no profile specified")
			}

			ssoConfig, err := config.GetSSOConfig(profile, homeDir)
			if err != nil {
				return fmt.Errorf("error retrieving SSO config: %w", err)
			}

			cacheFiles, err := os.ReadDir(filepath.Join(homeDir, ".aws", "sso", "cache"))
			if err != nil {
				return fmt.Errorf("error retrieving cache files - perhaps you need to login?: %w", err)
			}

			files := make([]fs.FileInfo, 0, len(cacheFiles))

			token, err := config.GetSSOToken(files, *ssoConfig, homeDir)
			if err != nil {
				return fmt.Errorf("error retrieving SSO token from cache files: %v", err)
			}

			sess := session.Must(session.NewSession())
			svc := sso.New(sess, aws.NewConfig().WithRegion(ssoConfig.Region))

			accountID = args[0]

			roles, err := svc.ListAccountRoles(&sso.ListAccountRolesInput{
				AccessToken: &token,
				MaxResults:  &results,
				AccountId:   &accountID,
			})
			if err != nil {
				return fmt.Errorf("error listing roles: %v", err)
			}

			writer := tabwriter.NewWriter(os.Stdout, tabwriterMinWidth, tabwriterWidth, tabwriterPadding, tabwriterPadChar, tabwriterFlags)
			fmt.Fprintln(writer, "ID\tROLE NAME")

			for _, results := range roles.RoleList {
				fmt.Fprintf(writer, "%s\t%s\n", *results.AccountId, *results.RoleName)
			}

			writer.Flush()

			return nil
		},
	}

	command.Flags().Int64VarP(&results, "results", "r", 10, "Maximum number of accounts to return")

	return command
}
