package main

import (
	"os"

	"github.com/spf13/viper"

	"github.com/jaxxstorm/aws-sso-creds/cmd/aws-sso-creds/export"
	"github.com/jaxxstorm/aws-sso-creds/cmd/aws-sso-creds/exportps"
	"github.com/jaxxstorm/aws-sso-creds/cmd/aws-sso-creds/get"
	"github.com/jaxxstorm/aws-sso-creds/cmd/aws-sso-creds/helper"
	"github.com/jaxxstorm/aws-sso-creds/cmd/aws-sso-creds/list"
	"github.com/jaxxstorm/aws-sso-creds/cmd/aws-sso-creds/set"
	"github.com/jaxxstorm/aws-sso-creds/cmd/aws-sso-creds/version"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

var (
	profile string
)

func configureCLI() *cobra.Command {
	rootCommand := &cobra.Command{
		Use:  "aws-sso-creds",
		Long: "A helper utility to interact with AWS SSO",
	}

	rootCommand.AddCommand(get.Command())
	rootCommand.AddCommand(set.Command())
	rootCommand.AddCommand(version.Command())
	rootCommand.AddCommand(export.Command())
	rootCommand.AddCommand(exportps.Command())
	rootCommand.AddCommand(list.Command())
	rootCommand.AddCommand(helper.Command())

	homeDir, err := homedir.Dir()

	if err != nil {
		panic("Cannot find home directory, fatal error")
	}

	rootCommand.PersistentFlags().StringVarP(&profile, "profile", "p", "", "the AWS profile to use")
	rootCommand.PersistentFlags().StringVarP(&homeDir, "home-directory", "H", homeDir, "specify a path to a home directory")
	if err := viper.BindEnv("profile", "AWS_PROFILE"); err != nil {
		panic(err)
	}
	if err := viper.BindPFlag("profile", rootCommand.PersistentFlags().Lookup("profile")); err != nil {
		panic(err)
	}
	if err := viper.BindPFlag("home-directory", rootCommand.PersistentFlags().Lookup("home-directory")); err != nil {
		panic(err)
	}

	return rootCommand
}

func main() {
	rootCommand := configureCLI()

	if err := rootCommand.Execute(); err != nil {
		os.Exit(1)
	}
}
