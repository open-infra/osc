package cmd

import (
	"fmt"

	"github.com/open-infra/osc/internal/color"
	"github.com/open-infra/osc/internal/config"
	"github.com/open-infra/osc/internal/ui"
	"github.com/spf13/cobra"
)

func infoCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "info",
		Short: "Print configuration info",
		Long:  "Print configuration information",
		Run: func(cmd *cobra.Command, args []string) {
			printInfo()
		},
	}
}

func printInfo() {
	const fmat = "%-25s %s\n"

	printLogo(color.Cyan)
	printTuple(fmat, "Configuration", config.OscConfigFile, color.Cyan)
	printTuple(fmat, "Logs", config.OscLogs, color.Cyan)
	printTuple(fmat, "Screen Dumps", config.OscDumpDir, color.Cyan)
}

func printLogo(c color.Paint) {
	for _, l := range ui.LogoSmall {
		fmt.Println(color.Colorize(l, c))
	}
	fmt.Println()
}
