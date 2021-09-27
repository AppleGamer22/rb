package tasks

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidate(t *testing.T) {
	testDir, err := os.MkdirTemp("", "test_*")
	require.NoError(t, err)
	writer := new(strings.Builder)
	var testCases = []struct {
		title           string
		srcRootDir      string
		dirsWriter      io.Writer
		filesWriter     io.Writer
		errorsWriter    io.Writer
		isErrorExpected bool
		expectedErrText string
	}{
		{
			title:           "inaccessible srcRootDir=>expected error",
			srcRootDir:      "invalid directory path",
			dirsWriter:      writer,
			filesWriter:     writer,
			errorsWriter:    writer,
			isErrorExpected: true,
			expectedErrText: "SrcRootDir: must be accessible.",
		}, {
			title:           "empty srcRootDir=>expected error",
			srcRootDir:      "",
			dirsWriter:      writer,
			filesWriter:     writer,
			errorsWriter:    writer,
			isErrorExpected: true,
			expectedErrText: "SrcRootDir: cannot be blank.",
		}, {
			title:           "nil dirs writer=>expect error",
			srcRootDir:      testDir,
			dirsWriter:      nil,
			filesWriter:     writer,
			errorsWriter:    writer,
			isErrorExpected: true,
			expectedErrText: "DirsWriter: cannot be blank.",
		}, {
			title:           "nil files writer=>expect error",
			srcRootDir:      testDir,
			dirsWriter:      writer,
			filesWriter:     nil,
			errorsWriter:    writer,
			isErrorExpected: true,
			expectedErrText: "FilesWriter: cannot be blank.",
		}, {
			title:           "nil errors writer=>expect error",
			srcRootDir:      testDir,
			dirsWriter:      writer,
			filesWriter:     writer,
			errorsWriter:    nil,
			isErrorExpected: true,
			expectedErrText: "ErrorsWriter: cannot be blank.",
		}, {
			title:           "valid input=>success expected",
			srcRootDir:      testDir,
			dirsWriter:      writer,
			filesWriter:     writer,
			errorsWriter:    writer,
			isErrorExpected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			in := &NewSrcListerInput{
				SrcRootDir:   tc.srcRootDir,
				DirsWriter:   tc.dirsWriter,
				FilesWriter:  tc.filesWriter,
				ErrorsWriter: tc.errorsWriter,
			}
			sourceLister, err := NewSourceLister(in)

			if tc.isErrorExpected {
				assert.EqualError(t, err, tc.expectedErrText)
				assert.Nil(t, sourceLister)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, sourceLister)
			}
		})
	}

}

func TestListSources(t *testing.T) {
	// given
	srcRootDir, err := os.MkdirTemp("", "msgListSrcDir_*")
	t.Log("source root dir: ", srcRootDir)
	require.NoError(t, err)

	testCases := []struct {
		title             string
		testDirName       string
		subDirs           []string
		files             []string
		symSource         string
		symLinks          []string
		expectedErrorsLog string
	}{
		{
			title:       "with an empty source",
			testDirName: "empty-source",
		}, {
			title:       "with a single dir",
			testDirName: "single-dir",
			subDirs:     []string{"one"},
		}, {
			title:       "with nested dirs",
			testDirName: "nested-dirs",
			files:       []string{"first.txt", "one/second.txt"},
			subDirs:     []string{"one", "one/two", "one/two-two", "one/two/three"},
		}, {
			title:       "with irregular files",
			testDirName: "nested-dirs",
			files:       []string{"first.txt", "one/second.txt"},
			subDirs:     []string{"", "one", "one/two", "one/two-two", "one/two/three"},
			symSource:   "sym_src.txt",
			symLinks:    []string{"sym_one", "one/sym_two"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			testPath := filepath.Join(srcRootDir, tc.testDirName)
			t.Log("test-path: ", testPath)
			setupTestDirs(t, err, testPath, tc.subDirs)
			setupTestFiles(t, testPath, tc.files)
			setupSymlink(t, testPath, tc.symSource, tc.symLinks)
			var dirPaths = strings.Builder{}
			var filePaths = strings.Builder{}
			var errorsLog = strings.Builder{}
			srcListerInput := &NewSrcListerInput{
				SrcRootDir:   testPath,
				DirsWriter:   &dirPaths,
				FilesWriter:  &filePaths,
				ErrorsWriter: &errorsLog,
			}
			lister, testErr := NewSourceLister(srcListerInput)
			assert.NoError(t, testErr)
			expectedDirs := append(tc.subDirs, testPath)
			expectedDirsSlice := getExpectations(srcListerInput.SrcRootDir, listExpectedPaths, expectedDirs)
			expectedFiles := append(tc.files, tc.symSource)
			expectedFilesSlice := getExpectations(srcListerInput.SrcRootDir, listExpectedPaths, expectedFiles)
			errorLogsSlice := getExpectations(srcListerInput.SrcRootDir, getExpectedErrorLogs, tc.symLinks)

			// when
			testErr = lister.Do()

			// then
			assert.NoError(t, testErr)
			actualDirsSlice := cleanOutput(dirPaths)
			actualFilesSlice := cleanOutput(filePaths)
			actualLogSlice := cleanOutput(errorsLog)

			assert.Equal(t, expectedDirsSlice, actualDirsSlice)
			assert.Equal(t, expectedFilesSlice, actualFilesSlice)
			assert.Equal(t, errorLogsSlice, actualLogSlice)
		})
	}
}

func cleanOutput(lines strings.Builder) []string {
	linesSlice := strings.Split(lines.String(), "\n")
	sort.Strings(linesSlice)
	if linesSlice[0] == "" {
		linesSlice = linesSlice[1:]
	}
	return linesSlice
}

type doFunc func(string, []string) []string

func listExpectedPaths(rootDir string, paths []string) []string {
	var out []string
	for _, p := range paths {
		if len(p) == 0 {
			continue
		}
		if len(strings.TrimPrefix(p, rootDir)) > 0 {
			fullPath := filepath.Join(rootDir, p)
			out = append(out, fullPath)
		} else {
			out = append(out, p)
		}
	}
	return out
}

func getExpectedErrorLogs(srcRootDir string, errorPathsSlice []string) []string {
	var out []string
	sort.Strings(errorPathsSlice)
	for _, p := range errorPathsSlice {
		if len(p) > 0 {
			path := filepath.Join(srcRootDir, p)
			line := fmt.Sprintf("path: %s, type: L--------- error_msg: unexpected_element", path)
			out = append(out, line)
		}
	}
	return out
}

func getExpectations(srcRootDir string, do doFunc, paths []string) []string {
	out := do(srcRootDir, paths)
	if len(out) > 0 {
		sort.Strings(out)
	} else {
		out = []string{}
	}
	return out
}

func setupTestDirs(t *testing.T, err error, testPath string, subDirs []string) {
	err = os.MkdirAll(testPath, 0755)
	require.NoError(t, err)
	for _, p := range subDirs {
		path := filepath.Join(testPath, p)
		err = os.MkdirAll(path, 0755)
		require.NoError(t, err)
	}
}

func setupTestFiles(t *testing.T, testPath string, files []string) {
	for _, f := range files {
		path := filepath.Join(testPath, f)
		file, err := os.Create(path)
		require.NoError(t, err)
		_, err = file.WriteString(fmt.Sprint(rand.Int()))
		require.NoError(t, err)
		err = file.Close()
		require.NoError(t, err)
	}
}

func setupSymlink(t *testing.T, testPath, symSource string, symLinks []string) {
	var symSrcFilePath string
	if len(symSource) > 0 {
		symSrcFilePath = filepath.Join(testPath, symSource)
		file, err := os.Create(symSrcFilePath)
		assert.NoError(t, err)
		err = file.Close()
		assert.NoError(t, err)
		_, err = file.WriteString(fmt.Sprint(rand.Int()))
		assert.Error(t, err)
		for _, sym := range symLinks {
			path := filepath.Join(testPath, sym)
			err := os.Symlink(symSrcFilePath, path)
			require.NoError(t, err)
		}
	}
}
