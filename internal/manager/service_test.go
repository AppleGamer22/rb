package manager

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/AppleGamer22/recursive-backup/internal/workers"

	"github.com/AppleGamer22/recursive-backup/internal/tasks"
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
				return fmt.Sprintf("dir-skeleton-error missed-path: %s%[2]cmissingSource%[2]csrc%[2]cmissing\n", testDirPath, filepath.Separator)
			},
		}, {
			title:       "empty source directory",
			testDirName: "emptySource",
			subDirs:     []string{},
			setupFunc: func(t *testing.T, testDirName string) string {
				testPath := filepath.Join(srcRootDir, testDirName)
				err = os.MkdirAll(filepath.Join(testPath, "src"), 0755)
				require.NoError(t, err)
				err = os.MkdirAll(filepath.Join(testPath, "target"), 0755)
				require.NoError(t, err)
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
				err = os.MkdirAll(filepath.Join(testPath, "src", "one"), 0755)
				require.NoError(t, err)
				err = os.MkdirAll(filepath.Join(testPath, "src", "two", "three"), 0755)
				require.NoError(t, err)
				err = os.MkdirAll(filepath.Join(testPath, "target"), 0755)
				require.NoError(t, err)
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

func TestFilesCopySuccess(t *testing.T) {
	// given
	setupTestFunc := func(t *testing.T, testDirPath string, fileSubPaths, missingPaths []string) (srcPath, targetPath string, srcFilesReader io.Reader) {
		targetDirPath := filepath.Join(testDirPath, "target")
		err := os.MkdirAll(targetDirPath, 0755)
		require.NoError(t, err)

		missingMap := make(map[string]bool)
		for _, p := range missingPaths {
			missingMap[p] = true
		}

		srcDirPath := filepath.Join(testDirPath, "src")
		err = os.MkdirAll(srcDirPath, 0755)
		require.NoError(t, err)
		var filesList strings.Builder
		allPaths := append(fileSubPaths, missingPaths...)
		for _, subPath := range allPaths {
			filePath := filepath.Join(srcDirPath, subPath)
			if _, ok := missingMap[subPath]; !ok {
				dirPath := filepath.Dir(filePath)
				err := os.MkdirAll(dirPath, 0755)
				require.NoError(t, err)
				err = os.WriteFile(filePath, []byte("Hello"), 0755)
				require.NoError(t, err)
			}
			_, err := filesList.WriteString(fmt.Sprintf("%s\n", filePath))
			require.NoError(t, err)
		}
		filesReader := strings.NewReader(filesList.String())

		return srcDirPath, targetDirPath, filesReader
	}

	expectedLogsFunc := func(t *testing.T, testRootDir string, filePaths, missingPaths []string) []string {
		var out []string
		for _, p := range filePaths {
			out = append(out, fmt.Sprintf("true,0,%s,%s,%s", filepath.Join(testRootDir, "target", p), filepath.Join(testRootDir, "src", p), "success"))
		}
		for _, p := range missingPaths {
			var errorMsg string
			if runtime.GOOS == "windows" {
				errorMsg = fmt.Sprintf("CreateFile %s: The system cannot find the path specified.", filepath.Join(testRootDir, "src", p))
			} else {
				errorMsg = fmt.Sprintf("stat %s: no such file or directory", filepath.Join(testRootDir, "src", p))
			}
			out = append(out, fmt.Sprintf("false,0,%s,%s,%s", filepath.Join(testRootDir, "target", p), filepath.Join(testRootDir, "src", p), errorMsg))
		}
		return out
	}

	testCases := []struct {
		title                string
		batchID              uint
		generalRequestChan   chan tasks.GeneralRequest
		responseChan         chan tasks.BackupFileResponse
		testDirName          string
		filesSubPaths        []string
		missingFilesSubPaths []string
	}{
		{
			title:              "with a single file",
			generalRequestChan: make(chan tasks.GeneralRequest, 1),
			responseChan:       make(chan tasks.BackupFileResponse, 1),
			batchID:            1,
			testDirName:        "singleFile",
			filesSubPaths:      []string{"file_one"},
		}, {
			title:              "with multiple files",
			generalRequestChan: make(chan tasks.GeneralRequest, 3),
			responseChan:       make(chan tasks.BackupFileResponse, 3),
			batchID:            2,
			testDirName:        "multipleFiles",
			filesSubPaths: []string{
				"one",
				"two",
				filepath.Join("three", "four", "five"),
			},
		}, {
			title:                "with missing files files",
			generalRequestChan:   make(chan tasks.GeneralRequest, 3),
			responseChan:         make(chan tasks.BackupFileResponse, 3),
			batchID:              3,
			testDirName:          "multipleFiles",
			missingFilesSubPaths: []string{filepath.Join("missing", "foo"), filepath.Join("missing", "bar")},
		},
	}

	testRootDir, err := os.MkdirTemp("", "testFilesCopy_*")
	t.Log("source root dir: ", testRootDir)
	require.NoError(t, err)

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			testDirPath := filepath.Join(testRootDir, tc.testDirName)
			srcTestPath, targetTestPath, srcFileReader := setupTestFunc(t, testDirPath, tc.filesSubPaths, tc.missingFilesSubPaths)
			api := NewService(ServiceInitInput{
				SourceRootDir: srcTestPath,
				TargetRootDir: targetTestPath,
			})
			expectedLogs := expectedLogsFunc(t, testDirPath, tc.filesSubPaths, tc.missingFilesSubPaths)
			sort.Strings(expectedLogs)

			var wgRequest sync.WaitGroup
			updateOnQuit := func() {
				wgRequest.Done()
			}

			for i := 0; i < cap(tc.generalRequestChan); i++ {
				wgRequest.Add(1)
				workers.NewCopyWorker(uint(i), srcTestPath, targetTestPath, tc.generalRequestChan, updateOnQuit)
			}
			var logWriter strings.Builder
			go api.HandleFilesCopyResponse(&logWriter, tc.responseChan)

			// when
			api.RequestFilesCopy(srcFileReader, tc.batchID, tc.generalRequestChan, tc.responseChan)
			for i := 0; i < cap(tc.generalRequestChan); i++ {
				tc.generalRequestChan <- tasks.QuitRequest{}
			}
			wgRequest.Wait()
			close(tc.generalRequestChan)
			api.WaitForAllResponses()
			close(tc.responseChan)
			// then
			// time.Sleep(time.Millisecond * 500)
			logString := strings.TrimSuffix(logWriter.String(), "\n")
			logSlices := strings.Split(logString, "\n")
			assert.Equal(t, len(tc.filesSubPaths)+len(tc.missingFilesSubPaths)+1, len(logSlices), logSlices)
			assert.Equal(t, "status,duration [milli-sec],target,source,error_message", logSlices[0])
			partialLogSlices := logSlices[1:]
			require.Len(t, partialLogSlices, len(tc.filesSubPaths)+len(tc.missingFilesSubPaths))
			sort.Strings(partialLogSlices)
			for i := 0; i < len(expectedLogs); i++ {
				expectedLine := expectedLogs[i]
				expectedLineItems := strings.Split(expectedLine, ",")
				actualLine := partialLogSlices[i]
				actualLineItems := strings.Split(actualLine, ",")

				assert.Equal(t, len(expectedLineItems), len(actualLineItems), "unexpected line length")
				assert.Equal(t, expectedLineItems[0], actualLineItems[0], "unexpected completion status")
				expectedDuration, err := strconv.Atoi(expectedLineItems[1])
				require.NoError(t, err, "unexpected duration")
				assert.GreaterOrEqual(t, expectedDuration, 0, "unexpected duration")
				actualDuration, err := strconv.Atoi(actualLineItems[1])
				require.NoError(t, err, "unexpected duration")
				assert.GreaterOrEqual(t, actualDuration, 0, "unexpected duration")

				expectedLogsMessage := fmt.Sprintf("expected logs: %v", expectedLogs)
				actualLogsMessage := fmt.Sprintf("actual lines: %v", partialLogSlices)
				assert.Equal(t, expectedLineItems[2], actualLineItems[2], "unexpected target value", expectedLogsMessage, actualLogsMessage)
				assert.Equal(t, expectedLineItems[3], actualLineItems[3], "unexpected source value", expectedLogsMessage, actualLogsMessage)
				assert.Equal(t, expectedLineItems[4], actualLineItems[4], "unexpected message", expectedLogsMessage, actualLogsMessage)
			}
			//assert.Equal(t, expectedLogs, logSlices[1:])

			for _, subPath := range tc.filesSubPaths {
				srcPath := filepath.Join(srcTestPath, subPath)
				srcFileData, err := os.ReadFile(srcPath)
				require.NoError(t, err)
				targetPath := filepath.Join(targetTestPath, subPath)
				targetFileData, err := os.ReadFile(targetPath)
				require.NoError(t, err)
				assert.Equal(t, srcFileData, targetFileData)
			}
		})
	}
}
