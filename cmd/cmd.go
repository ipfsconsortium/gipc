package cmd

import (
	"fmt"
	"os"

	cmd "github.com/ipfsconsortium/gipc/commands"
	cfg "github.com/ipfsconsortium/gipc/config"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// cfgFile is the configuration file path.
	cfgFile string
	// verbose is the verbosity level used in logrus.
	verbose string
)

// RootCmd represents the base command when called without any subcommands.
var RootCmd = &cobra.Command{
	Use:   "gipc",
	Short: "IPFS pinning consortium",
	Long:  "IPFS pinning consortium",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync",
	Long:  "Sync",
	Run:   cmd.Sync,
}

var syncOnceCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync",
	Long:  "Sync",
	Run:   cmd.Sync,
}

var dbDumpCmd = &cobra.Command{
	Use:   "db-dump",
	Short: "Dumps the database",
	Long:  "Dumps the database",
	Run:   cmd.DumpDb,
}

var dbInitCmd = &cobra.Command{
	Use:   "db-init",
	Short: "Initializes the database",
	Long:  "Initialized the database",
	Run:   cmd.InitDb,
}

var ipfscInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize ipfsc",
	Long:  "Initialize ipfsc",
	Run:   cmd.IpfscInit,
}

var ipfscLsCmd = &cobra.Command{
	Use:   "ls",
	Short: "Info of local ens",
	Long:  "Info of local ens",
	Run:   cmd.IpfscLs,
}

var ipfscAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add hash to IPFS",
	Long:  "Add hash to IPFS",
	Run:   cmd.IpfscAdd,
}
var ipfscRmCmd = &cobra.Command{
	Use:   "rm",
	Short: "Remove hash to IPFS",
	Long:  "Remove hash to IPFS",
	Run:   cmd.IpfscRemove,
}

// ExecuteCmd adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func ExecuteCmd() {

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)

	}
}

// init is called when the package loads and initializes cobra.
func init() {

	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file")
	RootCmd.PersistentFlags().StringVar(&verbose, "verbose", "INFO", "verbose level")

	RootCmd.AddCommand(syncCmd)

	RootCmd.AddCommand(dbDumpCmd)
	RootCmd.AddCommand(dbInitCmd)

	RootCmd.AddCommand(ipfscInitCmd)
	RootCmd.AddCommand(ipfscLsCmd)
	RootCmd.AddCommand(ipfscAddCmd)
	RootCmd.AddCommand(ipfscRmCmd)

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	if logLevel, err := log.ParseLevel(verbose); err == nil {
		log.SetLevel(logLevel)
	} else {
		panic(err)
	}

	viper.SetConfigType("yaml")
	viper.SetConfigName("gipc")  // name ofconfig file (without extension)
	viper.AddConfigPath(".")     // adding current directory as first search path
	viper.AddConfigPath("$HOME") // adding home directory as first search path
	viper.SetEnvPrefix("GIPC")   // so viper.AutomaticEnv will get matching envvars starting with O2M_
	viper.AutomaticEnv()         // read in environment variables that match

	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	log.WithField("file", viper.ConfigFileUsed()).Debug("Using config file")

	if err := viper.Unmarshal(&cfg.C); err != nil {
		panic(err)
	}

}
