package cmd

import (
	"fmt"
	"os"

	"github.com/mpppk/hlb/git"
	"github.com/mpppk/hlb/hlblib"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "hlb",
	Short: "multi git hosting service manager",
	Long:  ``,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		bypassCmds := []string{"create", "version", "init", "add-service"}
		configFilePath, err := hlblib.GetConfigDirPath()
		if err != nil {
			hlblib.PanicIfErrorExist(err)
		}

		for _, bypassCmd := range bypassCmds {
			if bypassCmd == cmd.Name() {
				return
			}
		}

		var config hlblib.Config
		err = viper.Unmarshal(&config)
		hlblib.PanicIfErrorExist(err)
		remote, err := git.GetDefaultRemote(".")
		hlblib.PanicIfErrorExist(err)
		serviceConfig, ok := config.FindServiceConfig(remote.ServiceHost)
		if !ok {
			fmt.Println(remote.ServiceHost, "is unknown host. Please add the service configuration to config file("+configFilePath+")")
			os.Exit(1)
		}
		if serviceConfig.Token == "" {
			serviceUrl := serviceConfig.Protocol + "://" + serviceConfig.Host
			addServiceCmd.Run(cmd, []string{serviceConfig.Type, serviceUrl})
		}
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	RootCmd.SetOutput(os.Stdout)
	if err := RootCmd.Execute(); err != nil {
		RootCmd.SetOutput(os.Stderr)
		RootCmd.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	configFilePath, err := hlblib.GetConfigFilePath()
	hlblib.PanicIfErrorExist(err)

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is "+configFilePath+")")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName(".hlb") // name of config file (without extension)
	configFilePath, err := hlblib.GetConfigDirPath()
	hlblib.PanicIfErrorExist(err)

	viper.AddConfigPath(configFilePath) // adding home directory as first search path
	viper.AutomaticEnv()                // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		initCmd.Run(nil, nil)
		err := viper.ReadInConfig()
		hlblib.PanicIfErrorExist(err)
	}
}
