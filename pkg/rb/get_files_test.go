package rb_test

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/AppleGamer22/recursive-backup/pkg/rb"
	"github.com/AppleGamer22/recursive-backup/pkg/utils"
	"github.com/stretchr/testify/assert"
)

// clears temporary directory of test-generated data
func ClearTemp(t *testing.T) {
	files, err := filepath.Glob(filepath.Join(os.TempDir(), "prefix-*"))
	assert.Nil(t, err)
	for _, file := range files {
		err = os.RemoveAll(file)
		assert.Nil(t, err)
	}
}

func TestInaccessibleFolder(t *testing.T) {
	ClearTemp(t)
	dirName1, err := ioutil.TempDir("", "prefix-")
	assert.Nil(t, err)
	err = os.Chmod(dirName1, 0220)
	assert.Nil(t, err)
	dirName2, err := ioutil.TempDir("", "prefix-")
	assert.Nil(t, err)
	assert.Nil(t, err)
	rber := rb.NewRecursiveBackupper(dirName1, dirName2, nil, false)
	sourcesLogPath, err := rber.BackupFilesSinceDate()
	assert.Nil(t, err)
	data, err := os.ReadFile(sourcesLogPath)
	assert.Nil(t, err)
	assert.True(t, strings.HasPrefix(string(data), "ERROR: "))
	err = os.Remove(sourcesLogPath)
	assert.Nil(t, err)
	ClearTemp(t)
}

func TestAccessibleFolder(t *testing.T) {
	ClearTemp(t)
	dirName1, err := ioutil.TempDir("", "prefix-")
	assert.Nil(t, err)
	dirName2, err := ioutil.TempDir("", "prefix-")
	assert.Nil(t, err)
	tempFile1, err := ioutil.TempFile(dirName1, "prefix-")
	assert.Nil(t, err)
	rber := rb.NewRecursiveBackupper(dirName1, dirName2, nil, false)
	sourcesLogPath, err := rber.BackupFilesSinceDate()
	assert.Nil(t, err)
	data, err := os.ReadFile(sourcesLogPath)
	assert.Nil(t, err)
	logs := string(data)
	tempFile2Path := strings.ReplaceAll(tempFile1.Name(), dirName1, dirName2)
	assert.True(t, strings.Contains(logs, tempFile1.Name()))
	assert.True(t, strings.Contains(logs, tempFile2Path))
	err = os.Remove(sourcesLogPath)
	assert.Nil(t, err)
	ClearTemp(t)
}

func TestRecoveryFileInTarget(t *testing.T) {
	ClearTemp(t)
	dirName1, err := ioutil.TempDir("", "prefix-")
	assert.Nil(t, err)
	dirName2, err := ioutil.TempDir("", "prefix-")
	assert.Nil(t, err)
	tempFile1, err := ioutil.TempFile(dirName1, "prefix-")
	assert.Nil(t, err)
	now := time.Now()
	rber := rb.NewRecursiveBackupper(dirName1, dirName2, &now, true)
	targetFilePath1, err := utils.Source2TargetPath(tempFile1.Name() ,dirName1, dirName2)
	assert.Nil(t, err)
	dest, err := os.Create(targetFilePath1)
	assert.Nil(t, err)
	defer dest.Close()
	_, err = io.Copy(dest, tempFile1)
	assert.Nil(t, err)
	found1, _ := utils.DoesTargetFileExist(tempFile1.Name(), dirName1, dirName2)
	assert.True(t, found1)
	sourcesLogPath, err := rber.BackupFilesSinceDate()
	assert.Nil(t, err)
	data, err := os.ReadFile(sourcesLogPath)
	assert.Nil(t, err)
	assert.False(t, strings.Contains(string(data), targetFilePath1))
	err = os.Remove(sourcesLogPath)
	assert.Nil(t, err)
	ClearTemp(t)
}

func TestRecoveryFileNotInTarget(t *testing.T) {
	ClearTemp(t)
	sourceDirName, err := ioutil.TempDir("", "src-dir-")
	assert.Nil(t, err)
	targetDirName, err := ioutil.TempDir("", "target-dir-")
	assert.Nil(t, err)
	tempFile1, err := ioutil.TempFile(sourceDirName, "test-file-")
	assert.Nil(t, err)
	now := time.Now()
	rber := rb.NewRecursiveBackupper(sourceDirName, targetDirName, &now, true)
	targetFilePath1, err := utils.Source2TargetPath(tempFile1.Name() ,sourceDirName, targetDirName)
	assert.Nil(t, err)
	found, _ := utils.DoesTargetFileExist(tempFile1.Name(), sourceDirName, targetDirName)
	assert.False(t, found)
	found, _ = utils.DoesTargetFileExist(tempFile1.Name(), sourceDirName, targetDirName)
	assert.False(t, found)
	sourcesLogPath, err := rber.BackupFilesSinceDate()
	assert.Nil(t, err)
	data, err := os.ReadFile(sourcesLogPath)
	assert.Nil(t, err)
	assert.True(t, strings.Contains(string(data), targetFilePath1))
	// err = os.Remove(sourcesLogPath)
	// assert.Nil(t, err)
	// ClearTemp(t)
}