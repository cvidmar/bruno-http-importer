package main

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"example.com/bruno-importer/importer"
	"github.com/stretchr/testify/require"

	_ "embed"
)

//go:embed "test/expect/one.bru"
var OneExpect []byte

//go:embed "test/expect/two.bru"
var TwoExpect []byte

func TestImport(t *testing.T) {
	cwd, errExec := os.Getwd()
	require.NoError(t, errExec)
	testDir := filepath.Join(cwd, "test")
	testInputDir := filepath.Join(testDir, "input")
	testOutputDir := filepath.Join(testDir, "output")
	errImport := importer.WalkDir(testInputDir, testOutputDir)
	require.NotEqual(t, io.EOF, errImport.Error())
	inputDirHandle, err := os.Open(testOutputDir)
	require.NoError(t, err)
	defer inputDirHandle.Close()
	for {
		file, err := inputDirHandle.Readdir(1)
		if err != nil {
			require.Equal(t, io.EOF, err)
		}
		if len(file) == 0 {
			break
		}
		fileName := file[0].Name()
		switch fileName {
		case "one.bru":
			one, err := os.ReadFile(testOutputDir + "/" + fileName)
			require.NoError(t, err)
			require.Equal(t, OneExpect, one)
		case "two.bru":
			two, err := os.ReadFile(testOutputDir + "/" + fileName)
			require.NoError(t, err)
			require.Equal(t, TwoExpect, two)
		}
	}
}
