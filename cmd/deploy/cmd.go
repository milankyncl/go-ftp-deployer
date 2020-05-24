package deploy

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/milankyncl/go-ftp-deployer/internal/client/ftp"
	"github.com/milankyncl/go-ftp-deployer/internal/deployer"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

type Command struct {
	start    time.Time
	progress *deployer.Progress
}

func New() *Command {
	return &Command{
		start:    time.Now(),
		progress: deployer.NewProgress(),
	}
}

func (c *Command) Execute(rootDirectory string, config deployer.Config) {
	c.progress.Set(color.FgHiWhite)
	c.progress.Message(fmt.Sprintf("Initializing deploy at [%s]", c.start.Local().Format(time.RFC3339)))

	c.progress.Message("Connecting to server")
	// TODO: Get host and credentials from deployment config
	client, err := ftp.NewClient("", "", "")
	if err != nil {
		log.Fatal("Could not create FTP connection", err)
	}
	defer client.Close()

	c.progress.Set(color.FgGreen)
	c.progress.Message("Successfully connected to FTP server")
	c.progress.Message("")

	// Running deploy check

	// TODO: Think about some temporary file to serve content
	// 		 for website (maintenance)

	// TODO: Update to collect files first, then upload them
	c.progress.Set(color.FgHiWhite)
	c.progress.Message("Uploading:")
	c.progress.Set(color.FgHiGreen)
	deployPath := path.Join(rootDirectory, config.LocalPath)
	err = filepath.Walk(
		deployPath,
		func(filepath string, file os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if filepath != deployPath {
				extPath := strings.ReplaceAll(filepath, deployPath, "")
				if file.IsDir() {
					if client.FileExists(path.Join(config.ExtPath, extPath)) == false {
						c.progress.Message("Creating directory ", extPath)
						err = client.CreateDir(
							path.Join(config.ExtPath, extPath),
						)
					}
				} else {
					c.progress.Message("[1/2]", extPath)
					err = client.Upload(
						filepath,
						path.Join(config.ExtPath, extPath),
					)
				}
				if err != nil {
					return err
				}
			}

			return nil
		})
	if err != nil {
		log.Fatal(err)
	}

	c.progress.Message("")

	end := time.Now()
	elapsed := end.Sub(c.start)
	c.progress.Set(color.FgGreen, color.Bold)
	c.progress.Message(fmt.Sprintf("Finished deploy at [%s] in %.2fs", end.Local().Format(time.RFC3339), elapsed.Seconds()))
}
