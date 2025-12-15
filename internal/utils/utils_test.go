package utils

import (
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func createTestFS(t *testing.T) string {
	workspace := t.TempDir()
	testDirs := []string{
		"solutions/go/hello-world/1",
		"solutions/go/two-fer/1",
		"solutions/go/two-fer/2",
		"solutions/python/hello-world/1",
		"solutions/rust/leap/1",
		"solutions/rust/leap/2",
		"solutions/rust/pangram/1",
		"solutions/rust/pangram/2",
		"solutions/rust/pangram/3",
	}
	for _, dir := range testDirs {
		fullPath := filepath.Join(workspace, dir)
		err := os.MkdirAll(fullPath, 0755)
		if err != nil {
			t.Fatalf("Failed to create test directory %s: %v", fullPath, err)
		}
	}
	return workspace
}

func TestGetLanguageCounts(t *testing.T) {
	t.Run("with a valid filesystem", func(t *testing.T) {
		workspace := createTestFS(t)
		expected := []Count{
			{Slug: "rust", Count: 2},
			{Slug: "go", Count: 1},
			{Slug: "python", Count: 0},
		}
		counts, err := GetLanguageCounts(workspace)
		if err != nil {
			fmt := "GetLanguageCounts() returned an unexpected error: %v"
			t.Fatalf(fmt, err)
		}
		if !reflect.DeepEqual(counts, expected) {
			fmt := "GetLanguageCounts() got = %v, want = %v"
			t.Errorf(fmt, counts, expected)
		}
	})

	t.Run("with an error from GetLanguages", func(t *testing.T) {
		workspace := t.TempDir()
		_, err := GetLanguageCounts(workspace)
		if err == nil {
			t.Errorf("GetLanguageCounts() expected an error, but got nil")
		}
	})
}

func TestGetLanguages(t *testing.T) {
	t.Run("with a valid filesystem", func(t *testing.T) {
		workspace := createTestFS(t)
		expected := []string{"go", "python", "rust"}
		languages, err := GetLanguages(workspace)
		if err != nil {
			t.Fatalf("GetLanguages() returned an unexpected error: %v", err)
		}
		if !reflect.DeepEqual(languages, expected) {
			t.Errorf("GetLanguages() got = %v, want = %v", languages, expected)
		}
	})

	t.Run("with a missing solutions directory", func(t *testing.T) {
		workspace := t.TempDir()
		_, err := GetLanguages(workspace)
		if err == nil {
			t.Errorf("GetLanguages() expected an error, but got nil")
		}
	})

	t.Run("with a solutions file", func(t *testing.T) {
		workspace := t.TempDir()
		_, err := os.Create(filepath.Join(workspace, "solutions"))
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
		_, err = GetLanguages(workspace)
		if err == nil {
			t.Errorf("GetLanguages() expected an error, but got nil")
		}
	})
}

func TestGetExerciseCounts(t *testing.T) {
	workspace := createTestFS(t)

	t.Run("for rust", func(t *testing.T) {
		expected := []Count{
			{Slug: "pangram", Count: 3},
			{Slug: "leap", Count: 2},
		}
		counts, err := GetExerciseCounts(workspace, "rust")
		if err != nil {
			fmt := "GetExerciseCounts() returned an unexpected error: %v\n"
			t.Fatalf(fmt, err)
		}
		if !reflect.DeepEqual(counts, expected) {
			fmt := "GetExerciseCounts() got = %v, want = %v"
			t.Errorf(fmt, counts, expected)
		}
	})

	t.Run("for go", func(t *testing.T) {
		expected := []Count{
			{Slug: "two-fer", Count: 2},
			{Slug: "hello-world", Count: 1},
		}
		counts, err := GetExerciseCounts(workspace, "go")
		if err != nil {
			fmt := "GetExerciseCounts() returned an unexpected error: %v\n"
			t.Fatalf(fmt, err)
		}
		if !reflect.DeepEqual(counts, expected) {
			fmt := "GetExerciseCounts() got = %v, want = %v"
			t.Errorf(fmt, counts, expected)
		}
	})

	t.Run("for non-existent language", func(t *testing.T) {
		_, err := GetExerciseCounts(workspace, "elixir")
		if err == nil {
			t.Error("GetExerciseCounts() expected an error, but got nil")
		}
	})

	t.Run("with a language file", func(t *testing.T) {
		workspace := createTestFS(t)
		err := os.RemoveAll(filepath.Join(workspace, "solutions", "rust"))
		if err != nil {
			t.Fatalf("Failed to remove test directory: %v", err)
		}
		_, err = os.Create(filepath.Join(workspace, "solutions", "rust"))
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
		_, err = GetExerciseCounts(workspace, "rust")
		if err == nil {
			t.Errorf("GetExerciseCounts() expected an error, but got nil")
		}
	})

	t.Run("with a non-readable exercise directory", func(t *testing.T) {
		workspace := createTestFS(t)
		exercisePath := filepath.Join(workspace, "solutions", "rust", "leap")
		err := os.Chmod(exercisePath, 0000)
		if err != nil {
			t.Fatalf("Failed to change directory permissions: %v", err)
		}

		rescueStderr := os.Stderr
		r, w, _ := os.Pipe()
		os.Stderr = w

		_, err = GetExerciseCounts(workspace, "rust")
		if err != nil {
			fmt := "GetExerciseCounts() returned an unexpected error: %v"
			t.Errorf(fmt, err)
		}

		w.Close()
		out, _ := io.ReadAll(r)
		os.Stderr = rescueStderr

		warn := "Warning: failed to read rust/leap: "
		if !strings.Contains(string(out), warn) {
			t.Errorf("Expected stderr to contain %q, got %q", warn, string(out))
		}

		err = os.Chmod(exercisePath, 0755)
		if err != nil {
			t.Fatalf("Failed to restore directory permissions: %v", err)
		}
	})
}
