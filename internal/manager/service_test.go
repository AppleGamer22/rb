package manager

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewManager(t *testing.T) {
	srcRootDir, err := os.MkdirTemp("", "managerSrcDir_*")
	require.NoError(t, err)
	t.Log("source root dir: ", srcRootDir)
	targetRootDir, err := os.MkdirTemp("", "managerTargetDir_*")
	require.NoError(t, err)
	t.Log("target root dir: ", targetRootDir)

	type testCase struct {
		title          string
		input          ServiceInitInput
		expectedErrStr string
	}
	testCases := []testCase{
		{
			title: "valid input=>success expected",
			input: ServiceInitInput{
				SourceRootDir: srcRootDir,
				TargetRootDir: targetRootDir,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			m := NewService(tc.input)
			assert.NoError(t, err)
			assert.NotNil(t, m)
		})
	}
}

func TestListSources(t *testing.T) {
	// given
	srcRootDir, err := os.MkdirTemp("", "testListSrcDir_*")
	t.Log("source root dir: ", srcRootDir)
	require.NoError(t, err)
	type setupFunc func(t *testing.T, dirPath string) string

	testCases := []struct {
		title             string
		testDirName       string
		setupDirFunc      setupFunc
		expectedDirPaths  func(string) string
		expectedFilePaths string
		expectedErrorsLog string
	}{
		{
			title:       "empty source=>expect empty output",
			testDirName: "emptySrc",
			setupDirFunc: func(t *testing.T, dirName string) string {
				fullPath := filepath.Join(srcRootDir, dirName)
				t.Log("srcPath: ", fullPath)
				err = os.MkdirAll(fullPath, 0755)
				t.Log("err: ", err)
				require.NoError(t, err)
				return fullPath
			},
			expectedDirPaths:  func(testRoot string) string { return fmt.Sprintf("%s\n", testRoot) },
			expectedFilePaths: "",
			expectedErrorsLog: "",
		}, {
			title:       "single dir",
			testDirName: "singleDir",
			setupDirFunc: func(t *testing.T, testDirName string) string {
				fullPath := filepath.Join(srcRootDir, testDirName)
				t.Log("srcPath: ", fullPath)
				err = os.MkdirAll(fullPath, 0755)
				require.NoError(t, err)
				pathOne := filepath.Join(fullPath, "one")
				err = os.MkdirAll(pathOne, 0755)
				require.NoError(t, err)
				return fullPath
			},
			expectedDirPaths: func(testRoot string) string {
				return fmt.Sprintf("%s\n%s"+string(filepath.Separator)+"%s\n", testRoot, testRoot, "one")
			},
			expectedFilePaths: "",
			expectedErrorsLog: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			dirsWriter := &strings.Builder{}
			filesWriter := &strings.Builder{}
			errorsWriter := &strings.Builder{}
			testRoot := tc.setupDirFunc(t, tc.testDirName)
			api := NewService(ServiceInitInput{
				SourceRootDir: testRoot,
			})

			// when
			err = api.ListSources(dirsWriter, filesWriter, errorsWriter)

			// then
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedDirPaths(testRoot), dirsWriter.String())
			assert.Equal(t, tc.expectedFilePaths, filesWriter.String())
			assert.Equal(t, tc.expectedErrorsLog, errorsWriter.String())
		})
	}
}

func TestCreateTargetDirSkeleton(t *testing.T) {
	// given
	srcRootDir, err := os.MkdirTemp("", "testCreateTargetDirSkeleton_*")
	t.Log("source root dir: ", srcRootDir)
	require.NoError(t, err)
	type setupFunc func(t *testing.T, dirPath string) string

	dirsFunc := func(basePath string, dirsSubPaths []string) io.Reader {
		buf := strings.Builder{}
		for _, d := range dirsSubPaths {
			path := filepath.Join(basePath, d)
			buf.WriteString(fmt.Sprintf("%s\n", path))
		}
		out := strings.TrimSuffix(buf.String(), "\n")
		return strings.NewReader(out)
	}

	testCases := []struct {
		title                 string
		testDirName           string
		subDirs               []string
		setupFunc             setupFunc
		isErrorExpected       bool
		expectedErrorString   string
		expectedDirPathsFunc  func(string) io.Reader
		expectedErrorsLogFunc func(string) string
	}{
		{
			title:       "missing source directory",
			testDirName: "missingSource",
			subDirs:     []string{"missing"},
			setupFunc: func(t *testing.T, testDirName string) string {
				return filepath.Join(srcRootDir, testDirName)
			},
			isErrorExpected:     true,
			expectedErrorString: "CreateTargetDirSkeleton completed with errors",
			expectedDirPathsFunc: func(testRootDir string) io.Reader {
				_ = testRootDir
				return nil
			},
			expectedErrorsLogFunc: func(testDirPath string) string {
				return fmt.Sprintf("dir-skeleton-error missed-path: %s/missingSource/src/missing\n", testDirPath)
			},
		}, {
			title:       "empty source directory",
			testDirName: "emptySource",
			subDirs:     []string{},
			setupFunc: func(t *testing.T, testDirName string) string {
				testPath := filepath.Join(srcRootDir, testDirName)
				os.MkdirAll(filepath.Join(testPath, "src"), 0755)
				os.MkdirAll(filepath.Join(testPath, "target"), 0755)
				return testPath
			},
			expectedDirPathsFunc: func(testRootDir string) io.Reader {
				var dirs string
				return strings.NewReader(dirs)
			},
			expectedErrorsLogFunc: func(testDirPath string) string {
				return ""
			},
		}, {
			title:       "empty source directory",
			testDirName: "emptySource",
			subDirs:     []string{"one", "two", filepath.Join("two", "three")},
			setupFunc: func(t *testing.T, testDirName string) string {
				testPath := filepath.Join(srcRootDir, testDirName)
				os.MkdirAll(filepath.Join(testPath, "src", "one"), 0755)
				os.MkdirAll(filepath.Join(testPath, "src", "two", "three"), 0755)
				os.MkdirAll(filepath.Join(testPath, "target"), 0755)
				return testPath
			},
			expectedDirPathsFunc: func(targetRootDir string) io.Reader {
				buf := strings.Builder{}
				buf.WriteString(fmt.Sprintf("%s\n", filepath.Join(targetRootDir, "one")))
				buf.WriteString(fmt.Sprintf("%s", filepath.Join(targetRootDir, "two", "three")))
				return strings.NewReader(buf.String())
			},
			expectedErrorsLogFunc: func(testDirPath string) string {
				return ""
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			testRoot := tc.setupFunc(t, tc.testDirName)
			testSrcDir := filepath.Join(testRoot, "src")
			testTargetDir := filepath.Join(testRoot, "target")
			dirsReader := dirsFunc(testSrcDir, tc.subDirs)
			errorsWriter := &strings.Builder{}
			api := NewService(ServiceInitInput{
				SourceRootDir: testSrcDir,
				TargetRootDir: testTargetDir,
			})

			// when
			createdDirsReader, err := api.CreateTargetDirSkeleton(dirsReader, errorsWriter)

			// then
			if tc.isErrorExpected {
				assert.EqualError(t, err, tc.expectedErrorString)
			} else {
				assert.NoError(t, err)
			}
			if tc.expectedDirPathsFunc(srcRootDir) == nil {
				assert.Nil(t, createdDirsReader)
			} else {
				assert.Equal(t, tc.expectedDirPathsFunc(testTargetDir), createdDirsReader)
			}
			assert.Equal(t, tc.expectedErrorsLogFunc(srcRootDir), errorsWriter.String())
		})
	}
}
