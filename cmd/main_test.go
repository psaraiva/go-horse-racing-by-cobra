package main

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/psaraiva/go-horse-racing-by-cobra/internal"
	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	t.Run("rootCmd is initialized", func(t *testing.T) {
		assert.NotNil(t, rootCmd, "rootCmd should be initialized")
		assert.Equal(t, "lab-go-horse-racing-by-cobra", rootCmd.Use)
		assert.Equal(t, "The GO race of the horses", rootCmd.Short)
	})

	t.Run("flags are configured", func(t *testing.T) {
		horseLabelFlag := rootCmd.Flags().Lookup("horse-label")
		assert.NotNil(t, horseLabelFlag, "horse-label flag should exist")
		assert.Equal(t, internal.HorseLabelDefault, horseLabelFlag.DefValue)

		horsesQuantityFlag := rootCmd.Flags().Lookup("horses-quantity")
		assert.NotNil(t, horsesQuantityFlag, "horses-quantity flag should exist")

		scoreTargetFlag := rootCmd.Flags().Lookup("score-target")
		assert.NotNil(t, scoreTargetFlag, "score-target flag should exist")

		gameTimeoutFlag := rootCmd.Flags().Lookup("game-timeout")
		assert.NotNil(t, gameTimeoutFlag, "game-timeout flag should exist")
		assert.Equal(t, internal.GameTimeoutDefault, gameTimeoutFlag.DefValue)
	})

	t.Run("rootCmd.Run function is set", func(t *testing.T) {
		assert.NotNil(t, rootCmd.Run, "rootCmd.Run should be set")
	})
}

func TestRootCmdExecution(t *testing.T) {
	t.Run("executes with default flags", func(t *testing.T) {
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		rootCmd.SetArgs([]string{"--game-timeout", "10s", "--score-target", "15", "--horses-quantity", "2"})

		done := make(chan bool)
		go func() {
			rootCmd.Execute()
			w.Close()
			done <- true
		}()

		<-done

		os.Stdout = oldStdout

		var buf bytes.Buffer
		io.Copy(&buf, r)
		output := buf.String()

		assert.Contains(t, output, "H0", "Output should contain horse labels")
	})

	t.Run("executes with custom flags", func(t *testing.T) {
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		rootCmd.SetArgs([]string{
			"--horse-label", "A",
			"--horses-quantity", "3",
			"--score-target", "20",
			"--game-timeout", "10s",
		})

		done := make(chan bool)
		go func() {
			rootCmd.Execute()
			w.Close()
			done <- true
		}()

		<-done
		os.Stdout = oldStdout

		var buf bytes.Buffer
		io.Copy(&buf, r)
		output := buf.String()

		// Should contain custom horse label
		assert.True(t, strings.Contains(output, "A0") || strings.Contains(output, "tired"), "Output should contain custom horse labels or timeout message")
	})
}
