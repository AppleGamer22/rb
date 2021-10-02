package tasks

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBackupDirSkeleton_Do_missing_source_dir(t *testing.T) {
	// given
	srcRootPath, err := os.MkdirTemp("", "srcDir_*")
	t.Log("srcRootPath: ", srcRootPath)
	require.NoError(t, err)
	err = os.MkdirAll(srcRootPath+string(os.PathSeparator)+"one", 0755)
	require.NoError(t, err)
	err = os.MkdirAll(srcRootPath+string(os.PathSeparator)+"two", 0755)
	require.NoError(t, err)
	err = os.MkdirAll(srcRootPath+string(os.PathSeparator)+"three", 0755)
	require.NoError(t, err)
	paths := fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n",
		srcRootPath,
		filepath.Join(srcRootPath, "one"),
		filepath.Join(srcRootPath, "two"),
		filepath.Join(srcRootPath, "three"),
		filepath.Join(srcRootPath, "no-real"))
	testDirReader := strings.NewReader(paths)
	targetRootPath, err := os.MkdirTemp("", "testTarget_*")
	require.NoError(t, err)
	testTask := backupDirSkeleton{
		SrcRootPath:          srcRootPath,
		SrcDirectoriesReader: testDirReader,
		TargetRootPath:       targetRootPath,
	}

	// when
	_, errs := testTask.Do()

	assert.Len(t, errs, 1)
	assert.EqualError(t, errs[0], fmt.Sprintf("missed directories: [%s]", filepath.Join(srcRootPath, "no-real")))
}

func TestBackupDirSkeleton_Do_success(t *testing.T) {
	// given
	srcRootPath, err := os.MkdirTemp("", "srcDir_*")
	t.Log("srcRootPath: ", srcRootPath)
	require.NoError(t, err)
	err = os.MkdirAll(srcRootPath+string(os.PathSeparator)+"one", 0755)
	require.NoError(t, err)
	err = os.MkdirAll(srcRootPath+string(os.PathSeparator)+"two", 0755)
	require.NoError(t, err)
	err = os.MkdirAll(srcRootPath+string(os.PathSeparator)+"three", 0755)
	require.NoError(t, err)
	paths := fmt.Sprintf("%s\n%s\n%s\n%s\n",
		srcRootPath,
		filepath.Join(srcRootPath, "one"),
		filepath.Join(srcRootPath, "two"),
		filepath.Join(srcRootPath, "three"))
	testDirReader := strings.NewReader(paths)
	targetRootPath, err := os.MkdirTemp("", "testTarget_*")
	require.NoError(t, err)
	testTask := backupDirSkeleton{
		SrcRootPath:          srcRootPath,
		SrcDirectoriesReader: testDirReader,
		TargetRootPath:       targetRootPath,
	}

	// when
	_, errs := testTask.Do()

	assert.Len(t, errs, 0)
}
