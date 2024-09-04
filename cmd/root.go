package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/ragarwalll/mta-forge.git/pkg/cli"
	"github.com/ragarwalll/mta-forge.git/pkg/constants"
	"github.com/ragarwalll/mta-forge.git/pkg/logger"
	"github.com/sagikazarmark/slog-shim"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// forgerArgs is the arguments for the forger
var forgerArgs = &cli.ForgerArgs{}

// rootCmd is the root command for the application
var rootCmd = &cobra.Command{
	Use:   "mf",
	Short: "Mta Forge is a tool to generate MTA & MTA Extension resources",
	Long:  `You can use Mta Forge to generate MTA & MTA Extension resources with templates.`,
}

// Execute is the entry point for the application
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// init is called before the command is run
func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&forgerArgs.BaseDir, "base-dir", "b", "", "Specify the base directory where the templates are located (default: current directory)")
	rootCmd.PersistentFlags().StringVarP(&forgerArgs.OutputDir, "output-dir", "o", "", "Specify the output directory where the generated files will be saved (default: ./output)")
	rootCmd.PersistentFlags().BoolVarP(&forgerArgs.Verbose, "verbose", "v", false, "Enable verbose output")
	rootCmd.PersistentFlags().BoolVarP(&forgerArgs.Local, "local", "l", false, "Enable local output format")
	rootCmd.PersistentFlags().BoolVar(&forgerArgs.ExpandSource, "expand-source", false, "Expand source information in logs")

	_ = viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	_ = viper.BindPFlag("local", rootCmd.PersistentFlags().Lookup("local"))
	_ = viper.BindPFlag("expand-source", rootCmd.PersistentFlags().Lookup("expand-source"))
	_ = viper.BindPFlag("output-dir", rootCmd.PersistentFlags().Lookup("output-dir"))
}

// initConfig is called before the command is run
func initConfig() {
	viper.AutomaticEnv()
}

// preRun is called before the command is run
func preRun() error {
	if forgerArgs.BaseDir == "" {
		forgerArgs.BaseDir = "."
	}

	if forgerArgs.OutputDir == "" {
		forgerArgs.OutputDir = path.Join(forgerArgs.BaseDir, "output")
		if err := os.MkdirAll(forgerArgs.OutputDir, constants.DefaultDirPermissions); err != nil {
			return fmt.Errorf("error creating output directory: %w", err)
		}
	}

	cli.SetForgerArgs(forgerArgs)
	logger.InitLogger()

	slog.Info("running mta-forge", "expand-source", forgerArgs.ExpandSource, "local", forgerArgs.Local, "verbose", forgerArgs.Verbose, "output-dir", forgerArgs.OutputDir)

	return nil
}
