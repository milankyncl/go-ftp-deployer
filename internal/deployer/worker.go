package deployer

import (
	"errors"
	"github.com/fatih/color"
	"github.com/milankyncl/go-ftp-deployer/internal/deployer/worker"
	"github.com/milankyncl/go-ftp-deployer/internal/platform/ftp"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

const (
	tempFile = ".deployment"
)

type Worker struct {
	config   Config
	progress *Progress
	client   *ftp.Client
}

func NewWorker(config Config, progress *Progress) *Worker {
	progress.Color(color.Reset, color.FgWhite).Message("Connecting to server")

	client, err := ftp.NewClient(config.Host, config.User, config.Password, config.ExtPath)
	if err != nil {
		log.Fatal("Could not create FTP connection", err)
	}

	progress.
		Color(color.FgGreen).
		Message("Successfully connected to FTP server")

	return &Worker{
		config:   config,
		progress: progress,
		client:   client,
	}
}

func (w *Worker) LastSyncTime() (*time.Time, error) {
	if w.client.FileExists(tempFile) {
		// Reach content
		f, err := w.client.ReadFile(tempFile)
		if err != nil &&
			err != io.ErrShortWrite {
			return nil, err
		}
		if f == nil {
			return nil, errors.New("unable to reach temp file")
		}
		t, err := time.Parse(time.RFC3339, string(f.Content()))
		if err != nil {
			return nil, err
		}
		return &t, nil
	}
	return nil, nil
}

func (w *Worker) CollectPaths(localPath string) []worker.Path {
	var paths []worker.Path

	err := filepath.Walk(
		localPath,
		func(filepath string, file os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if filepath != localPath {
				extPath := strings.ReplaceAll(filepath, localPath, "")
				paths = append(paths, worker.NewPath(extPath, file))
			}
			return nil
		})
	if err != nil {
		log.Println(err)

	}

	return paths
}

func (w *Worker) Upload(localPath string, p worker.Path) error {
	var err error
	if p.Info().IsDir() {
		if !w.client.FileExists(p.Path()) {
			err = w.client.CreateDir(p.Path())
		}
	} else {
		err = w.client.Upload(path.Join(localPath, p.Path()), p.Path())
	}
	return err
}

func (w *Worker) UpdateLastSync() error {
	file, err := ioutil.TempFile(".", "*.deployment")
	if err != nil {
		log.Fatal(err)
	}
	_, err = file.Write([]byte(time.Now().Format(time.RFC3339)))
	if err != nil {
		log.Fatalln(err)
	}

	tfs, err := os.Open(file.Name())
	if err != nil {
		log.Fatal(err)
	}

	defer os.Remove(file.Name())
	return w.client.WriteFile(tempFile, tfs)
}

func (w *Worker) Disconnect() {
	w.client.Close()
}
