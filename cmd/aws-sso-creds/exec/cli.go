package exec

import (
	"fmt"
	"os"
	osexec "os/exec"
	"syscall"

	"github.com/jaxxstorm/aws-sso-creds/pkg/credentials"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	getSSOCredentials = credentials.GetSSOCredentials
	lookPath          = osexec.LookPath
	execCommand       = syscall.Exec
)

func Command() *cobra.Command {
	command := &cobra.Command{
		Use:   "exec -- <command> [args...]",
		Short: "Execute a command with AWS temporary credentials in the environment",
		Long:  "Execute a command with AWS temporary credentials injected as environment variables",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			cmd.SilenceUsage = true

			profile := viper.GetString("profile")
			homeDir := viper.GetString("home-directory")

			creds, _, cfg, err := getSSOCredentials(profile, homeDir)
			if err != nil {
				return err
			}

			binary, err := lookPath(args[0])
			if err != nil {
				return fmt.Errorf("command not found: %s", args[0])
			}

			env := os.Environ()
			env = append(env,
				fmt.Sprintf("AWS_ACCESS_KEY_ID=%s", *creds.RoleCredentials.AccessKeyId),
				fmt.Sprintf("AWS_SECRET_ACCESS_KEY=%s", *creds.RoleCredentials.SecretAccessKey),
				fmt.Sprintf("AWS_SESSION_TOKEN=%s", *creds.RoleCredentials.SessionToken),
				fmt.Sprintf("AWS_DEFAULT_REGION=%s", cfg.Region),
			)

			return execCommand(binary, args, env)
		},
	}

	return command
}
