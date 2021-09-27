package tasks

import (
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBackupFile_Do_Success(t *testing.T) {
	// given
	srcRootPath, err := os.MkdirTemp("", "srcDir_*")
	fileName := "test_file.txt"
	srcFilePath := filepath.Join(srcRootPath, fileName)
	srcFile, err := os.Create(srcFilePath)
	require.NoError(t, err)
	defer srcFile.Close()
	testText := "testing123\n"
	n, err := srcFile.WriteString(testText)
	require.NoError(t, err)
	assert.Equal(t, n, len(testText))
	targetRootPath, err := ioutil.TempDir("", "testTarget_*")
	require.NoError(t, err)
	targetFilePath := filepath.Join(targetRootPath, fileName)
	testResponseChannel := make(chan BackupFileResponse)
	now := time.Now()
	testTask := BackupFile{
		CreationRequestTime: now,
		SourcePath:          srcFilePath,
		TargetPath:          targetFilePath,
		ResponseChannel:     testResponseChannel,
	}

	// when
	go testTask.Do()

	// then
	resp := <-testTask.ResponseChannel
	t.Log(resp)
	assert.Equal(t, true, resp.CompletionStatus)
	assert.Equal(t, srcFilePath, resp.SourcePath)
	assert.Equal(t, targetFilePath, resp.TargetPath)
	assert.Equal(t, now, resp.CreationRequestTime)
	assert.True(t, now.Before(resp.CompletionTime))
}

func TestBackupFile_Do_Fail(t *testing.T) {
	// given
	srcRootPath, err := os.MkdirTemp("", "srcDir_*")
	fileName := "test_file.txt"
	srcFilePath := filepath.Join(srcRootPath, fileName)
	srcFile, err := os.Create(srcFilePath)
	require.NoError(t, err)
	defer srcFile.Close()
	os.Chmod(srcFilePath, fs.ModeIrregular)
	targetRootPath, err := ioutil.TempDir("", "testTarget_*")
	require.NoError(t, err)
	targetFilePath := filepath.Join(targetRootPath, fileName)
	testResponseChannel := make(chan BackupFileResponse)
	now := time.Now()
	testTask := BackupFile{
		CreationRequestTime: now,
		SourcePath:          srcFilePath,
		TargetPath:          targetFilePath,
		ResponseChannel:     testResponseChannel,
	}

	// when
	go testTask.Do()

	// then
	resp := <-testTask.ResponseChannel
	t.Log(resp)
	assert.Equal(t, false, resp.CompletionStatus)
	assert.Equal(t, srcFilePath, resp.SourcePath)
	assert.Equal(t, targetFilePath, resp.TargetPath)
	assert.True(t, now.Before(resp.CompletionTime))
}
