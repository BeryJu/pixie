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
	Use: "pixie [directory to serve]",
	Run: func(cmd *cobra.Command, args []string) {
		if config.Current.Debug {
			log.SetLevel(log.DebugLevel)
		} else {
			log.SetFormatter(&log.JSONFormatter{})
		}
		if len(args) < 1 {
			log.Infof("No root-directory given, defaulting to Working Directory of %s", config.Current.RootDir)
		} else {
			config.Current.RootDir = args[0]
		}
		log.Infof("pixie with config %+v", config.Current)
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
	rootCmd.PersistentFlags().BoolVar(&config.Current.Debug, "debug", config.Defaults.Debug, "Enable debug-mode")
	rootCmd.PersistentFlags().BoolVar(&config.Current.Silent, "silent", config.Defaults.Silent, "Enable silent mode (no access logs)")
	rootCmd.PersistentFlags().BoolVar(&config.Current.EXIFPurgeGPS, "exif-purge-gps", config.Defaults.EXIFPurgeGPS, "Purge GPS-Related EXIF metadata")
	rootCmd.PersistentFlags().BoolVar(&config.Current.CacheEnabled, "cache-enabled", config.Defaults.CacheEnabled, "Enable in-memory cache")
	rootCmd.PersistentFlags().IntVar(&config.Current.CacheEviction, "cache-eviction", config.Defaults.CacheEviction, "Time after which entry can be evicted (in minutes)")
	rootCmd.PersistentFlags().IntVar(&config.Current.CacheMaxSize, "cache-max-size", config.Defaults.CacheMaxSize, "Maximum Cache size in MB (0 disables the limit)")
	rootCmd.PersistentFlags().IntVar(&config.Current.CacheMaxItemSize, "cache-max-item-size", config.Defaults.CacheMaxItemSize, "Maximum Item size to cache (in bytes)")
}
