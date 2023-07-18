package download

import (
	"os"
	"path/filepath"
)

func (jm *JobManager) createDriveIdFile(path string, driveId string) error {
	path = filepath.Join(path, driveIdFileName)

	if err := os.WriteFile(path, []byte(driveId), 0644); err != nil {
		jm.logger.Errorf("failed to write drive id file. %v", err)
		return err
	}

	return nil
}

func (jm *JobManager) readDriveIdFile(path string) (string, error) {
	path = filepath.Join(path, driveIdFileName)

	buf, err := os.ReadFile(path)

	if err != nil {
		jm.logger.Errorf("failed to read drive id file. %v", err)
		return "", err
	}

	driveId := string(buf)

	return driveId, nil
}
