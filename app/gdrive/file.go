package gdrive

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/avast/retry-go"
	"google.golang.org/api/drive/v3"
)

type DriveFile struct {
	Remote     *drive.File
	Descriptor *os.File
	Path       string
	Size       int64
}

func (ds *DriveService) DownloadFile(driveFile *DriveFile) error {
	return retry.Do(func() error {
		ds.logger.Infof("Starting to download file: %s", driveFile.Remote.Name)

		if err := ds.getFileMetadata(driveFile); err != nil {
			ds.logger.Errorf("failed to get file metadata. %v", err)
			return err
		}

		ds.logger.Infof("Retrieved metadata of file '%s'", driveFile.Remote.Name)

		defer driveFile.Descriptor.Close()

		if driveFile.Size > 0 {
			if ds.checkWhetherFileIsCompleted(driveFile) {
				ds.logger.Infof("Download of file '%s' is already finished. Skipping", driveFile.Remote.Name)
				return nil
			}

			if ds.checkWhetherFileIsCorrupted(driveFile) {
				ds.logger.Infof("Local file '%s' is probably corrupted. Removing it..", driveFile.Descriptor.Name())

				if err := truncate(driveFile.Descriptor); err != nil {
					ds.logger.Errorf("Failed to truncate local file '%s' file. %v", driveFile.Descriptor.Name(), err)
					return err
				}
			}
		}

		content, err := ds.requestFileContent(driveFile)

		if err != nil {
			ds.logger.Errorf("Failed to fetch content of file '%s' from Google Drive. %v", driveFile.Remote.Name, err)
			return err
		}

		ds.logger.Infof("Fetched content of file '%s' from Google Drive.", driveFile.Remote.Name)

		w, err := io.Copy(driveFile.Descriptor, *content)
		driveFile.Size = w

		if err != nil {
			ds.logger.Errorf("Failed to write content to file. %v", err)
			return err
		}

		md5Checksum, err := getMd5Checksum(driveFile)

		if err != nil {
			ds.logger.Errorf("failed to calculate md5 checksum. %v", err)
			return err
		}

		if md5Checksum != driveFile.Remote.Md5Checksum {
			err := errors.New("checksum of local file != checksum of remote file. file is probably corrupted")
			ds.logger.Error(err)
			return err
		}

		ds.logger.Info("finished downloading file")

		return nil
	}, retry.Attempts(ds.conf.Download.RetryThreeshold))
}

func (ds *DriveService) checkWhetherFileIsCompleted(driveFile *DriveFile) bool {
	if driveFile.Size == driveFile.Remote.Size {
		md5Checksum, err := getMd5Checksum(driveFile)

		if err != nil {
			ds.logger.Errorf("failed to calculate md5 checksum. %v", err)
			return false
		}

		if md5Checksum != driveFile.Remote.Md5Checksum {
			return false
		}
	}

	return true
}

func (ds *DriveService) checkWhetherFileIsCorrupted(driveFile *DriveFile) bool {
	return driveFile.Size > driveFile.Remote.Size
}

func (ds *DriveService) requestFileContent(driveFile *DriveFile) (*io.ReadCloser, error) {
	request := ds.drive.Files.Get(driveFile.Remote.Id).
		SupportsAllDrives(true).
		SupportsTeamDrives(true)

	request.Header().Add("Range", fmt.Sprintf("bytes=%d-", driveFile.Size))

	response, err := request.Download()

	if err != nil {
		return nil, err
	}

	return &response.Body, err
}

func (ds *DriveService) getFileMetadata(driveFile *DriveFile) error {
	if err := getFileDescriptor(driveFile); err != nil {
		return err
	}

	if err := getFileSize(driveFile); err != nil {
		return err
	}

	return nil
}

func getFileDescriptor(driveFile *DriveFile) error {
	if err := os.MkdirAll(filepath.Dir(driveFile.Path), 0644); err != nil {
		return err
	}

	descriptor, err := os.OpenFile(driveFile.Path, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)

	if err != nil {
		return err
	}

	driveFile.Descriptor = descriptor

	return nil
}

func getFileSize(driveFile *DriveFile) error {
	stat, err := driveFile.Descriptor.Stat()

	if err != nil {
		return err
	}

	driveFile.Size = stat.Size()

	return nil
}

func getMd5Checksum(driveFile *DriveFile) (string, error) {
	hash := md5.New()

	_, err := driveFile.Descriptor.Seek(0, io.SeekStart)

	if err != nil {
		return "", err
	}

	if _, err := io.Copy(hash, driveFile.Descriptor); err != nil {
		return "", err
	}

	md5Checksum := fmt.Sprintf("%x", hash.Sum(nil))

	return md5Checksum, nil
}

func truncate(file *os.File) error {
	if err := file.Truncate(0); err != nil {
		return err
	}

	_, err := file.Seek(0, io.SeekStart)

	if err != nil {
		return err
	}

	return nil
}
