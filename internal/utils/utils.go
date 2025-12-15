package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strconv"
)

type Count struct {
	Count int
	Slug  string
}

func sortCounts(counts []Count) {
	slices.SortFunc(counts, func(a, b Count) int {
		if a.Count != b.Count {
			return b.Count - a.Count
		}
		if a.Slug < b.Slug {
			return -1
		}
		if a.Slug > b.Slug {
			return 1
		}
		return 0
	})
}

func GetExerciseCounts(workspace, language string) ([]Count, error) {
	path := filepath.Join(workspace, "solutions", language)
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("language directory does not exist: %s", path)
	}

	solutionCounts := make(map[string]int)
	exerciseEntries, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read language directory: %w", err)
	}

	for _, exerciseEntry := range exerciseEntries {
		if !exerciseEntry.IsDir() {
			continue
		}
		exerciseSlug := exerciseEntry.Name()
		exercisePath := filepath.Join(path, exerciseSlug)
		solutionEntries, err := os.ReadDir(exercisePath)
		if err != nil {
			format := "Warning: failed to read %s/%s: %v\n"
			fmt.Fprintf(os.Stderr, format, language, exerciseSlug, err)
			continue
		}
		for _, solutionEntry := range solutionEntries {
			if !solutionEntry.IsDir() {
				continue
			}
			_, err := strconv.Atoi(solutionEntry.Name())
			if err == nil {
				solutionCounts[exerciseSlug]++
			}
		}
	}

	var counts []Count
	for slug, count := range solutionCounts {
		counts = append(counts, Count{Count: count, Slug: slug})
	}

	sortCounts(counts)
	return counts, nil
}

func GetLanguageCounts(workspace string) ([]Count, error) {
	languages, err := GetLanguages(workspace)
	if err != nil {
		return nil, err
	}

	var counts []Count
	for _, language := range languages {
		exerciseCounts, err := GetExerciseCounts(workspace, language)
		if err != nil {
			return nil, err
		}

		total := 0
		for _, count := range exerciseCounts {
			if count.Count > 1 {
				total++
			}
		}

		counts = append(counts, Count{Count: total, Slug: language})
	}

	sortCounts(counts)
	return counts, nil
}

func GetLanguages(workspace string) ([]string, error) {
	path := filepath.Join(workspace, "solutions")
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("solutions directory does not exist: %s", path)
	}

	langEntries, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read solutions directory: %w", err)
	}

	var languages []string
	for _, langEntry := range langEntries {
		if langEntry.IsDir() {
			languages = append(languages, langEntry.Name())
		}
	}
	slices.Sort(languages)
	return languages, nil
}
