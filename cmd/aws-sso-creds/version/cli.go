package version

import (
	"fmt"
	"os"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/pulumi/pulumictl/pkg/gitversion"
	"github.com/spf13/cobra"
	"github.com/thuannfq/aws-sso-creds/pkg/version"
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
				version, err := gitversion.GetLanguageVersions(workingDir, plumbing.Revision(commitish))
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
