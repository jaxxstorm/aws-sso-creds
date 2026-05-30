package version

import (
	"fmt"
	"os"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/jaxxstorm/aws-sso-creds/pkg/version"
	"github.com/pulumi/pulumictl/pkg/gitversion"
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	command := &cobra.Command{
		Use:   "version",
		Short: "Get the current version",
		Long:  `Get the current version of pulumictl`,
		RunE: func(cmd *cobra.Command, args []string) error {

			v := version.Version
			// If we haven't set a version with linker flags, use this tool to get the version
			if v == "" {
				commitish := "HEAD"
				workingDir, err := os.Getwd()
				if err != nil {
					return fmt.Errorf("error obtaining working directory: %w", err)
				}
				repo, err := git.PlainOpen(workingDir)
				if err != nil {
					return fmt.Errorf("error opening git repository: %w", err)
				}
				version, err := gitversion.GetLanguageVersions(repo, plumbing.Revision(commitish), false, "", false)
				if err != nil {
					return fmt.Errorf("error calculating version: %w", err)
				}

				v = version.SemVer
			}
			fmt.Println(v)
			return nil
		},
	}
	return command
}
