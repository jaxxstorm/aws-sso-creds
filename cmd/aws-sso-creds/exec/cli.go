package exec

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/jaxxstorm/aws-sso-creds/pkg/credentials"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

			creds, _, err := credentials.GetSSOCredentials(profile, homeDir)
			if err != nil {
				return err
			}

			binary, err := exec.LookPath(args[0])
			if err != nil {
				return fmt.Errorf("command not found: %s", args[0])
			}

			env := os.Environ()
			env = append(env,
				fmt.Sprintf("AWS_ACCESS_KEY_ID=%s", *creds.RoleCredentials.AccessKeyId),
				fmt.Sprintf("AWS_SECRET_ACCESS_KEY=%s", *creds.RoleCredentials.SecretAccessKey),
				fmt.Sprintf("AWS_SESSION_TOKEN=%s", *creds.RoleCredentials.SessionToken),
			)

			return syscall.Exec(binary, args, env)
		},
	}

	return command
}
