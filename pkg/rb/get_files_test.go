package rb_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/AppleGamer22/recursive-backup/pkg/rb"
	"github.com/stretchr/testify/assert"
)

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
	sourcesLogPath, _, err := rb.GetFilePaths(dirName1, dirName2)
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
	sourcesLogPath, _, err := rb.GetFilePaths(dirName1, dirName2)
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
