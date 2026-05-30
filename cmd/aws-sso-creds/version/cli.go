package version

import (
	"fmt"
	"os"

	appversion "github.com/jaxxstorm/aws-sso-creds/pkg/version"
	"github.com/jaxxstorm/vers"
	"github.com/spf13/cobra"
)

var calculateFallbackVersion = func() (string, error) {
	workingDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("error obtaining working directory: %w", err)
	}
	repo, err := vers.OpenRepository(workingDir)
	if err != nil {
		return "", fmt.Errorf("error opening git repository: %w", err)
	}
	version, err := vers.Calculate(vers.Options{
		Repository: repo,
	})
	if err != nil {
		return "", fmt.Errorf("error calculating version: %w", err)
	}

	return version.SemVer, nil
}

func Command() *cobra.Command {
	command := &cobra.Command{
		Use:   "version",
		Short: "Get the current version",
		Long:  `Get the current version of aws-sso-creds`,
		RunE: func(cmd *cobra.Command, args []string) error {

			v := appversion.Version
			// If we haven't set a version with linker flags, use this tool to get the version
			if v == "" {
				version, err := calculateFallbackVersion()
				if err != nil {
					return err
				}

				v = version
			}
			fmt.Println(v)
			return nil
		},
	}
	return command
}
