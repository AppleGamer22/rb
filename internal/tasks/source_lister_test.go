package tasks

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"os"
	"strings"
	"testing"
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
