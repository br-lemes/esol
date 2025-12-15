package cmd

import (
	"fmt"

	"github.com/br-lemes/esol/internal/utils"
	"github.com/spf13/cobra"
)

var singleCmd = &cobra.Command{
	Use:   "single [language]",
	Short: "Lists exercises with a single solution",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var (
			counts  []utils.Count
			err     error
			minimum int
		)
		if len(args) == 1 {
			minimum = 1
			counts, err = utils.GetExerciseCounts(workspace, args[0])
			if err != nil {
				return err
			}
		} else {
			minimum = 0
			counts, err = utils.GetLanguageCounts(workspace)
			if err != nil {
				return err
			}
		}
		for _, count := range counts {
			if count.Count == minimum {
				fmt.Println(count.Slug)
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(singleCmd)
}
