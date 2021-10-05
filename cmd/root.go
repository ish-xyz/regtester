package cmd

import (
	"fmt"

	"github.com/ish-xyz/regtester/pkg/config"
	"github.com/spf13/cobra"
)

var (
	// Used for flags.
	cfgFile string
	verbose bool

	rootCmd = &cobra.Command{
		Use:   "regtester",
		Short: "A tool to test registries",
		Long: `Regtester is a tool to run load tests against registries 
triggering parallel docker pulls and returning a detailed result `,
		Run: run,
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "Config file absolute path")
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "Run regtester in debug mode")
	rootCmd.MarkPersistentFlagRequired("config")
	//override flags
}

func run(cmd *cobra.Command, args []string) {
	cfg, _ := config.Load(cfgFile)

	fmt.Println(cfg)
}
