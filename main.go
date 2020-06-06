package main

import (
	"fmt"
	"github.com/milankyncl/go-ftp-deployer/cmd/deploy"
	"github.com/milankyncl/go-ftp-deployer/internal/deployer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
)

const (
	appVersion = "v0.1"
	appName    = "deployer"

	defaultConfigFile = "deployment.yml"

	configServerHost     = "deployer.server.host"
	configServerUser     = "deployer.server.user"
	configServerPassword = "deployer.server.password"
	configLocalPath      = "deployer.localPath"
	configExternalPath   = "deployer.externalPath"
)

var (
	configFile = new(string)

	rootCmd = &cobra.Command{
		Use:   appName,
		Short: "FTP Deployer - Written in GO",
		Run: func(cmd *cobra.Command, args []string) {
			err := cmd.Help()
			if err != nil {
				log.Fatal("app help error")
			}
		},
	}
	deployCmd = &cobra.Command{
		Use:   "deploy",
		Short: "Deploy stuff on FTP server",
		Run: func(cmd *cobra.Command, args []string) {
			rootDir, err := os.Getwd()
			if err != nil {
				log.Fatal(err)
			}

			deploy.Execute(
				rootDir,
				deployer.Config{
					Host:      viper.GetString(configServerHost),
					User:      viper.GetString(configServerUser),
					Password:  viper.GetString(configServerPassword),
					LocalPath: viper.GetString(configLocalPath),
					ExtPath:   viper.GetString(configExternalPath),
				},
			)
		},
	}
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Show the Deployer version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(
				fmt.Sprintf("FTP Deployer - Written in GO, %s", appVersion),
			)
		},
	}
)

func main() {
	cobra.OnInitialize(initConfig)

	rootCmd.AddCommand(
		versionCmd,
		deployCmd,
	)

	rootCmd.PersistentFlags().StringVarP(configFile, "config", "c", "deployment.yml", "configuration file, default is `deployment.yml`")

	config()

	_ = rootCmd.Execute()
}

func initConfig() {
	viper.AddConfigPath(".")

	if *configFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(*configFile)
	} else {
		viper.SetConfigFile(defaultConfigFile)
	}

	if err := viper.ReadInConfig(); err != nil {
		log.Panic(err)
	}
}

func config() {
	viper.SetDefault(configServerHost, "")
	viper.SetDefault(configServerHost, "")
	viper.SetDefault(configServerUser, "")
	viper.SetDefault(configServerPassword, "")
	viper.SetDefault(configLocalPath, "")
	viper.SetDefault(configExternalPath, "")
}
