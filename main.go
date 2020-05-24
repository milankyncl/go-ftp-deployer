package main

import (
	"github.com/milankyncl/go-ftp-deployer/cmd/deploy"
	"github.com/milankyncl/go-ftp-deployer/internal/deployer"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var (
	rootCmd = &cobra.Command{
		Use:   "deploy",
		Short: "Deploy stuff on FTP server",
		Run: func(cmd *cobra.Command, args []string) {
			rootDir, err := os.Getwd()
			if err != nil {
				log.Fatal(err)
			}

			dep := deploy.New()
			dep.Execute(
				rootDir,
				deployer.Config{
					LocalPath: "example/",
					ExtPath:   "/test/",
				},
			)
		},
	}
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal("Something went wrong")
	}
}
