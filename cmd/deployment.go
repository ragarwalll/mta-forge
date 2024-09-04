package cmd

import (
	"log/slog"

	"github.com/ragarwalll/mta-forge.git/pkg/constants"
	"github.com/ragarwalll/mta-forge.git/pkg/forger"
	"github.com/spf13/cobra"
)

var deploymentCmd = &cobra.Command{
	Use:   "deployment",
	Short: "Generate mta.yml file",
	Long:  `Use this command to generate mta.yml files for your application.`,
	Run: func(_ *cobra.Command, _ []string) {
		if err := preRun(); err != nil {
			slog.Error(constants.ErrPreRunning, "error", err)
			return
		}

		forger := forger.NewForger(forgerArgs.BaseDir, forgerArgs.OutputDir)
		if err := forger.Generate("deployment"); err != nil {
			slog.Error("Error generating deployment", "error", err)
		}
	},
}

func init() {
	generateCmd.AddCommand(deploymentCmd)
}
