package deploy

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/milankyncl/go-ftp-deployer/internal/deployer"
	"github.com/milankyncl/go-ftp-deployer/internal/deployer/worker"
	"path"
	"time"
)

func Execute(rootDirectory string, config deployer.Config) {
	progress := deployer.NewProgress()
	start := time.Now()

	progress.
		Color(color.FgHiWhite, color.Bold).
		Message(fmt.Sprintf("Starting deploy at [%s]", start.Local().Format(time.RFC3339)))

	wrkr := deployer.NewWorker(config, progress)
	defer func() {
		if r := recover(); r != nil {
			progress.Color(color.FgWhite).Message("Something went wrong, disconnecting from server. Error:")
			progress.Color(color.FgRed).Message(r)
		}
		progress.Color(color.Bold, color.FgHiWhite).Message("Disconnected from server.")
		wrkr.Disconnect()
	}()

	// Running deploy check
	lastSync, err := wrkr.LastSyncTime()
	if err != nil {
		panic(err)
	}
	if lastSync == nil {
		progress.Color(color.FgWhite).Message("Deploying all contents")
	}

	// TODO: Think about some temporary file to serve content
	// 		 for website (maintenance)

	progress.Color(color.FgWhite).Message("Collecting files.")

	localPath := path.Join(rootDirectory, config.LocalPath)
	paths := wrkr.CollectPaths(localPath)
	toDep := make([]worker.Path, 0)

	for _, p := range paths {
		if lastSync == nil || lastSync.Before(p.Info().ModTime()) {
			toDep = append(toDep, p)
		}
	}

	progress.Message("")

	if len(toDep) == 0 {
		progress.Color(color.Bold, color.FgGreen).Message("Everything already synchronized!")
		return
	}

	progress.Color(color.FgHiWhite).Message("Uploading:")
	for i, p := range toDep {
		progress.
			Color(color.FgHiGreen).
			Message(fmt.Sprintf("(%d/%d) %s", i+1, len(paths), p.Path()))
		err := wrkr.Upload(path.Join(rootDirectory, config.LocalPath), p)
		if err != nil {
			panic(err)
		}
	}

	// Save last sync
	err = wrkr.UpdateLastSync()
	if err != nil {
		panic(err)
	}

	end := time.Now()
	elapsed := end.Sub(start)
	progress.
		Message().
		Color(color.FgGreen, color.Bold).
		Message(fmt.Sprintf("Finished deploy at [%s] in %.2fs", end.Local().Format(time.RFC3339), elapsed.Seconds()))
}
