package manager

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewManager(t *testing.T) {
	srcDir, err := os.MkdirTemp("", "managerSrcDir_*")
	require.NoError(t, err)
	targetDir, err := os.MkdirTemp("", "managerTargetDir_*")
	writer := new(strings.Builder)
	reader := new(strings.Reader)

	type testCase struct {
		title          string
		input          NewManagerInput
		expectedErrStr string
	}
	testCases := []testCase{
		{
			title: "empty input=>error expected",
			input: NewManagerInput{
				SourceRootDir:          "",
				TargetRootDir:          "",
				ListingDirPathsWriter:  nil,
				ListingFilePathsWriter: nil,
				ListingErrorsLogWriter: nil,
				FilePathsReader:        nil,
				FileCopyPipelineLength: 0,
				FileBackupLogWriter:    nil,
			},
			expectedErrStr: "FileBackupLogWriter: cannot be blank; FileCopyPipelineLength: cannot be blank; FilePathsReader: cannot be blank; ListingDirPathsWriter: cannot be blank; ListingErrorsLogWriter: cannot be blank; ListingFilePathsWriter: cannot be blank; SourceRootDir: cannot be blank; TargetRootDir: cannot be blank.",
		}, {
			title: "invalid pipeline=>error expected",
			input: NewManagerInput{
				SourceRootDir:          srcDir,
				TargetRootDir:          targetDir,
				ListingDirPathsWriter:  writer,
				ListingFilePathsWriter: writer,
				ListingErrorsLogWriter: writer,
				FilePathsReader:        reader,
				FileCopyPipelineLength: 0,
				FileBackupLogWriter:    writer,
			},
			expectedErrStr: "FileCopyPipelineLength: cannot be blank.",
		}, {
			title: "success",
			input: NewManagerInput{
				SourceRootDir:          srcDir,
				TargetRootDir:          targetDir,
				ListingDirPathsWriter:  writer,
				ListingFilePathsWriter: writer,
				ListingErrorsLogWriter: writer,
				FilePathsReader:        reader,
				FileCopyPipelineLength: 2,
				FileBackupLogWriter:    writer,
			},
			expectedErrStr: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			m, err := NewManager(tc.input)
			if len(tc.expectedErrStr) > 0 {
				assert.EqualError(t, err, tc.expectedErrStr)
				assert.Nil(t, m)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, m)
				assertedManager, ok := m.(*manager)
				assert.True(t, ok)
				assert.Len(t, assertedManager.FileBackupWorkers, tc.input.FileCopyPipelineLength)
				assert.Equal(t, cap(assertedManager.TasksPipeline), tc.input.FileCopyPipelineLength)
			}
		})
	}
}
