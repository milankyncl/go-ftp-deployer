package main

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

var (
	logger = logrus.New()

	rootCmd = &cobra.Command{
		Use:   "deploy",
		Short: "Deploy stuff on FTP server",
		Run: func(cmd *cobra.Command, args []string) {
			//
		},
	}
)

func main() {
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.InfoLevel)
	logger.SetFormatter(&logrus.TextFormatter{
		ForceColors: true,
		DisableTimestamp: true,
	})
	logger.Info("Starting the deployer")

	if err := rootCmd.Execute(); err != nil {
		logger.Fatal("Something went wrong")
	}
}
