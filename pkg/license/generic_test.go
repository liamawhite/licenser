// Copyright (c) Tetrate, Inc 2021 All Rights Reserved.

package license

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReader(t *testing.T) {
	t.Run("Reader has correct bytes form string", func(t *testing.T) {
		a := FromTemplateString(license, mark, 2019, "Test")
		want, _ := ioutil.ReadFile("testdata/apache.golden")
		got, _ := ioutil.ReadAll(a.Reader())
		assert.Equal(t, want, got)
	})

	t.Run("Reader has correct bytes from file ", func(t *testing.T) {
		a := FromTemplateFile("testdata/apache.golden", mark, 2019, "Test")
		want, _ := ioutil.ReadFile("testdata/apache.golden")
		got, _ := ioutil.ReadAll(a.Reader())
		assert.Equal(t, want, got)
	})
}

func TestIsPresent(t *testing.T) {
	tests := []struct {
		name      string
		inputFile string
		want      bool
	}{
		{
			name:      "License is present",
			inputFile: "testdata/apache.golden",
			want:      true,
		},
		{
			name:      "License is not present",
			inputFile: "testdata/nolicense.golden",
			want:      false,
		},
	}
	for _, tt := range tests {
		tc := tt
		t.Run(tc.name, func(t *testing.T) {
			a := FromTemplateFile("testdata/apache.golden", mark, 0, "") // The presence check doesn't care about these values
			inputReader, _ := os.Open(tc.inputFile)
			assert.Equal(t, tc.want, a.IsPresent(inputReader))
		})
	}
}
