package cmd

import (
	"log/slog"

	"github.com/ragarwalll/mta-forge.git/pkg/constants"
	"github.com/ragarwalll/mta-forge.git/pkg/forger"
	"github.com/spf13/cobra"
)

var extensionCmd = &cobra.Command{
	Use:   "extension",
	Short: "Generate *.mtaext files",
	Long:  `Use this command to generate *.mtaext files for your application.`,
	Run: func(_ *cobra.Command, _ []string) {
		if err := preRun(); err != nil {
			slog.Error(constants.ErrPreRunning, "error", err)
			return
		}

		forger := forger.NewForger(forgerArgs.BaseDir, forgerArgs.OutputDir)
		if err := forger.Generate("extension"); err != nil {
			slog.Error("Error generating descriptor", "error", err)
		}
	},
}

func init() {
	generateCmd.AddCommand(extensionCmd)
}
