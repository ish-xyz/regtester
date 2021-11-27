package cmd

import (
	"fmt"
	"os"

	"github.com/ish-xyz/regtester/pkg/config"
	"github.com/ish-xyz/regtester/pkg/controller"
	"github.com/ish-xyz/regtester/pkg/docker"
	"github.com/ish-xyz/regtester/pkg/logger"
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

	if err := cfg.Validate(); err != nil {
		logger.DebugLogger.Println("Configuration loaded: ", cfg)
		logger.ErrorLogger.Printf("configuration not valid!\n")
		os.Exit(1)
	}
	logger.DebugLogger.Printf("Configuration valid!\n")

	_ = controller.NewController(cfg)
	logger.DebugLogger.Printf("Controller loaded!\n")

	docker := docker.NewDockerClient(cfg)
	logger.DebugLogger.Printf("Docker client load!\n")

	report, _ := docker.Pull("http://localhost:5000", "rundeck/rundeck:latest")
	fmt.Println(report)
	//cannot use cfg (variable of type *config.Config) as *invalid type value in argument to controller.NewControllerc
}
