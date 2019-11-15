package cmd

import (
	"fmt"
	"os"

	"git.beryju.org/BeryJu.org/pixie/pkg/config"
	"git.beryju.org/BeryJu.org/pixie/pkg/server"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use: "pixie",
	Run: func(cmd *cobra.Command, args []string) {
		if config.Current.Debug {
			log.SetLevel(log.DebugLevel)
		} else {
			log.SetFormatter(&log.JSONFormatter{})
		}
		log.Infof("pixie starting on port %d", config.Current.Port)
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
	config.Current = config.Defaults
	rootCmd.PersistentFlags().StringVarP(&config.Current.RootDir, "root-dir", "r", config.Defaults.RootDir, "Root directory to serve.")
	rootCmd.PersistentFlags().BoolVar(&config.Current.Debug, "debug", config.Defaults.Debug, "Enable debug-mode.")
	rootCmd.PersistentFlags().BoolVar(&config.Current.EXIFPurgeGPS, "exif-purge-gps", config.Defaults.EXIFPurgeGPS, "Purge GPS-Related EXIF metadata.")
	rootCmd.PersistentFlags().BoolVar(&config.Current.CacheEnabled, "cache-enabled", config.Defaults.CacheEnabled, "Enable in-memory cache")
	rootCmd.PersistentFlags().IntVar(&config.Current.CacheMaxItems, "cache-max-items", config.Defaults.CacheMaxItems, "Maximum Items to cache")
	rootCmd.PersistentFlags().IntVar(&config.Current.CacheMaxItemSize, "cache-max-item-size", config.Defaults.CacheMaxItemSize, "Maximum Item size to cache (in bytes).")
}
