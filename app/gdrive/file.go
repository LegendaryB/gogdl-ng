package gdrive

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/LegendaryB/gogdl-ng/app/logging"
	"google.golang.org/api/drive/v3"
)

var logger = logging.NewLogger()

func DownloadFile(folderPath string, driveFile *drive.File) error {
	logger.Infof("starting to download file (name: %s, id: %s)", driveFile.Name, driveFile.Id)

	fp := filepath.Join(folderPath, driveFile.Name)
	file, err := getLocalFile(fp, driveFile.Size)

	if err != nil {
		logger.Errorf("failed to acquire local file: %w", err)
		return err
	}

	defer file.Close()

	request := service.Files.Get(driveFile.Id).
		SupportsAllDrives(true).
		SupportsTeamDrives(true)

	fi, err := file.Stat()

	if err != nil {
		logger.Errorf("failed to stat() file. %w", err)
		return err
	}

	if fi.Size() == driveFile.Size {
		logger.Infof("file is already finished. skipping..")
		return nil
	}

	request.Header().Add("Range", fmt.Sprintf("bytes=%d-", fi.Size()))

	response, err := request.Download()

	if err != nil {
		logger.Errorf("failed to fetch file. %w", err)
		return err
	}

	_, err = io.Copy(file, response.Body)

	if err != nil {
		logger.Errorf("failed to write buffer to file. %w", err)
		return err
	}

	md5checksum, err := getMd5Checksum(fp)

	if err != nil {
		logger.Errorf("")
		return err
	}

	if md5checksum != driveFile.Md5Checksum {
		err = errors.New("md5 checksum mismatch")
		logger.Errorf("the md5 checksum of the local file does not match checksum of the remote file. %w", err)
		return err
	}

	logger.Info("finished file download")
	return nil
}

func getMd5Checksum(path string) (string, error) {
	file, err := os.Open(path)

	if err != nil {
		logger.Errorf("failed to open file. %w", err)
		return "", err
	}

	hash := md5.New()

	if _, err = io.Copy(hash, file); err != nil {
		logger.Errorf("failed to write buffer. %w", err)
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

func getLocalFile(path string, maxSize int64) (*os.File, error) {
	fi, err := os.Stat(path)

	if errors.Is(err, os.ErrNotExist) {
		f, err := os.Create(path)

		if err != nil {
			logger.Errorf("failed to create file. %w", err)
			return nil, err
		}

		return f, nil
	}

	if fi.Size() > maxSize {
		logger.Warnf("local file size is greater than remote. file is probably corrupt. removing it..")
		os.Remove(path)
	}

	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)

	if err != nil {
		logger.Errorf("failed to open file. %w", err)
		return nil, err
	}

	return f, nil
}
