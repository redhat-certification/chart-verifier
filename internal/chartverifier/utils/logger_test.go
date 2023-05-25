package utils

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"strings"
	"testing"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func NewTestCmd(config *viper.Viper) *cobra.Command {
	return &cobra.Command{
		Use:   "test",
		Short: "A cobra command for testing logging",
		Long:  "A cobra command for testing logging",
		Run: func(cmd *cobra.Command, args []string) {
			switch args[0] {
			case "LogInfoOnly":
				InitLog(cmd, "", true)
				LogInfo(args[1])
			case "StdOutOnly":
				InitLog(cmd, "", true)
				WriteStdOut(args[1])
			case "LogInfoFile":
				InitLog(cmd, "", false)
				LogInfo(args[1])
			case "ErrorAndLog":
				InitLog(cmd, "", false)
				LogError(args[1])
			case "WarningAndLog":
				InitLog(cmd, "", false)
				LogWarning(args[1])
			case "StdOutFile":
				InitLog(cmd, "report.yaml", false)
				WriteStdOut(args[1])
			case "testPrune":
				InitLog(cmd, "", false)
				LogError("just for pruning")
			default:
				fmt.Println("No test", args[0])
			}
			WriteLogs("yaml")
		},
	}
}

func TestLogging(t *testing.T) {
	t.Run("LogInfo, no logfile, should be no output", func(t *testing.T) {
		tCmd := NewTestCmd(viper.New())
		outBuf := bytes.NewBufferString("")
		CmdStdout = outBuf
		errBuf := bytes.NewBufferString("")
		CmdStderr = errBuf
		tCmd.SetArgs([]string{"LogInfoOnly", "This should not be output"})
		require.NoError(t, tCmd.Execute())

		require.Empty(t, outBuf.String())
		require.Empty(t, errBuf.String())
	})

	t.Run("LogInfo, set logfile, should be a logfile created", func(t *testing.T) {
		tCmd := NewTestCmd(viper.New())
		outBuf := bytes.NewBufferString("")
		CmdStdout = outBuf
		errBuf := bytes.NewBufferString("")
		CmdStderr = errBuf

		output := "This should be output in a log file"
		tCmd.SetArgs([]string{"LogInfoFile", output})
		require.NoError(t, tCmd.Execute())

		require.Empty(t, outBuf.String())
		require.Empty(t, errBuf.String())

		require.True(t, checkAndOrDeleteFiles("log", output), fmt.Sprintf("Expected string not found in logs : %s", output))

		t.Cleanup(func() {
			checkAndOrDeleteFiles("delete", "don't care just be sure log is deleted")
		})
	})

	t.Run("Write to stdOut to a file", func(t *testing.T) {
		tCmd := NewTestCmd(viper.New())
		outBuf := bytes.NewBufferString("")
		CmdStdout = outBuf
		errBuf := bytes.NewBufferString("")
		CmdStderr = errBuf
		output := "This should be output to file only"
		tCmd.SetArgs([]string{"StdOutFile", output})
		require.NoError(t, tCmd.Execute())

		require.Empty(t, outBuf.String())
		require.Empty(t, errBuf.String())
		require.True(t, checkAndOrDeleteFiles("report", output), fmt.Sprintf("Expected string not found in report : %s", output))

		t.Cleanup(func() {
			checkAndOrDeleteFiles("delete", "don't care just be sure report is deleted")
		})
	})

	t.Run("Write Error to stdError and file", func(t *testing.T) {
		tCmd := NewTestCmd(viper.New())
		outBuf := bytes.NewBufferString("")
		CmdStdout = outBuf
		errBuf := bytes.NewBufferString("")
		CmdStderr = errBuf
		output := "This error should be output to stderr and log file"
		tCmd.SetArgs([]string{"ErrorAndLog", output})
		require.NoError(t, tCmd.Execute())

		require.Empty(t, outBuf.String())
		require.True(t, strings.Contains(errBuf.String(), output))
		require.True(t, checkAndOrDeleteFiles("log", output), fmt.Sprintf("Expected string not found in logs : %s", output))

		t.Cleanup(func() {
			checkAndOrDeleteFiles("delete", "don't care just be sure log is deleted")
		})
	})

	t.Run("Write Warning to stdError and file", func(t *testing.T) {
		tCmd := NewTestCmd(viper.New())
		outBuf := bytes.NewBufferString("")
		CmdStdout = outBuf
		errBuf := bytes.NewBufferString("")
		CmdStderr = errBuf
		output := "This warning should be output to stderr and log file."
		tCmd.SetArgs([]string{"WarningAndLog", output})
		require.NoError(t, tCmd.Execute())

		require.Empty(t, outBuf.String())
		require.True(t, strings.Contains(errBuf.String(), output))
		require.True(t, checkAndOrDeleteFiles("log", output), fmt.Sprintf("Expected string not found in logs : %s", output))

		t.Cleanup(func() {
			checkAndOrDeleteFiles("delete", "don't care just be sure log is deleted")
		})
	})

	t.Run("Test log file pruning", func(t *testing.T) {
		tCmd := NewTestCmd(viper.New())
		outBuf := bytes.NewBufferString("")
		CmdStdout = outBuf
		errBuf := bytes.NewBufferString("")
		CmdStderr = errBuf

		tCmd.SetArgs([]string{"testPrune"})
		for i := 1; i <= 15; i++ {
			require.NoError(t, tCmd.Execute())
			numlogFiles := howManyLogFiles()
			if i < 10 {
				require.True(t, numlogFiles == i, fmt.Sprintf("expected %d logfile but found %d", i, numlogFiles))
			} else {
				require.True(t, numlogFiles == 10, fmt.Sprintf("expected 10 logfile but found %d", numlogFiles))
			}
			time.Sleep(2 * time.Second)
		}

		t.Cleanup(func() {
			checkAndOrDeleteFiles("delete", "don't care just be sure log file are deleted")
		})
	})
}

func checkAndOrDeleteFiles(fileType string, expectedContent string) bool {
	result := false
	if len(expectedContent) > 0 {
		currentDir, err := os.Getwd()
		if err != nil {
			fmt.Println(fmt.Sprintf("er"+
				"ror getting current working directory : %s", err))
			return false
		}
		logFilesPath := path.Join(currentDir, outputDirectory)

		files, err := os.ReadDir(logFilesPath)
		if err != nil {
			if fileType != "delete" {
				fmt.Printf("error reading log directory : %s : %s\n", logFilesPath, err)
				return false
			}
			return true
		}

		for _, file := range files {
			foundFile := false
			if fileType == "log" && strings.HasPrefix(file.Name(), "verifier") && strings.HasSuffix(file.Name(), ".log") {
				foundFile = true
			} else if fileType == "report" && file.Name() == "report.yaml" {
				foundFile = true
			} else if fileType == "delete" {
				foundFile = true
			}
			if foundFile {
				filePath := path.Join(logFilesPath, file.Name())
				if fileType != "delete" {
					logfileContent, err := os.ReadFile(filePath)
					if err != nil {
						fmt.Printf("error reading file %s", err)
						return false
					}
					if strings.Contains(string(logfileContent), expectedContent) {
						result = true
					}
				}
				os.Remove(filePath)
			}
		}
		os.Remove(logFilesPath)
	}
	return result
}

func howManyLogFiles() int {
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("error getting current working directory : %s\n", err)
		return 0
	}
	logFilesPath := path.Join(currentDir, outputDirectory)

	files, err := os.ReadDir(logFilesPath)
	if err != nil {
		fmt.Printf("error reading log directory : %s : %s\n", logFilesPath, err)
		return 0
	}
	numLogFiles := 0

	for _, file := range files {
		if strings.HasPrefix(file.Name(), "verifier") && strings.HasSuffix(file.Name(), ".log") {
			numLogFiles++
		}
	}
	return numLogFiles
}
