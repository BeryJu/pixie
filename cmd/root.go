package cmd

import (
	"fmt"
	"os"

	"git.beryju.org/BeryJu.org/pixie/pkg/server"

	"git.beryju.org/BeryJu.org/pixie/pkg/config"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use: "pixie",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		if config.CfgDebug {
			log.SetLevel(log.DebugLevel)
		} else {
			log.SetFormatter(&log.JSONFormatter{})
		}
		log.Info("pixie starting on port 8080")
		server.Run()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&config.CfgRootDir, "root-dir", "r", ".", "Root directory to serve.")
	rootCmd.PersistentFlags().BoolVar(&config.CfgDebug, "debug", false, "Enable debug-mode.")
	rootCmd.PersistentFlags().BoolVar(&config.CfgPurgeExifGPS, "purge-exif-gps", true, "Purge GPS-Relateed EXIF metadata.")
}
