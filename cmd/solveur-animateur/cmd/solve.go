package cmd

import (
	"fmt"
	"sort"

	"github.com/fatih/color"
	"github.com/jmichiels/SolveurAnimateur"
	"github.com/spf13/cobra"
)

var csv bool

// solveCmd represents the solve command
var solveCmd = &cobra.Command{
	Use: "solve",
	Run: func(cmd *cobra.Command, args []string) {
		_, distribution := solveuranimateur.Solve()

		if csv {
			const sep = `,`

			fmt.Print(sep)
			for index := range solveuranimateur.Numbers.Indexes() {
				fmt.Printf("%d%s", solveuranimateur.Numbers[solveuranimateur.NumberIndex(index)], sep)
			}

			var lines []string
			for personIndex, _ := range distribution {
				line := fmt.Sprintf("\n\"%s\"%s", solveuranimateur.Persons[personIndex], sep)
				for numberIndex := range distribution[personIndex] {
					line += fmt.Sprintf("%.3f%s", distribution[personIndex][numberIndex], sep)
				}
				lines = append(lines, line)
			}
			sort.Strings(lines)
			for _, line := range lines {
				fmt.Print(line)
			}
			fmt.Println()
		} else {
			bold := color.New(color.FgWhite, color.Bold)

			fmt.Printf("%20s", "")
			for index := range solveuranimateur.Numbers.Indexes() {
				bold.Printf("%4d ", solveuranimateur.Numbers[solveuranimateur.NumberIndex(index)])
			}

			var lines []string
			for personIndex, _ := range distribution {
				line := bold.Sprintf("\n%-20s", solveuranimateur.Persons[personIndex])
				for numberIndex := range distribution[personIndex] {
					probability := distribution[personIndex][numberIndex]
					if probability > 0 {
						line += fmt.Sprintf("%.2f ", probability)
					} else {
						line += fmt.Sprint(" --  ")
					}
				}
				lines = append(lines, line)
			}
			sort.Strings(lines)
			for _, line := range lines {
				fmt.Print(line)
			}
			fmt.Println()
		}
	},
}

func init() {
	rootCmd.AddCommand(solveCmd)

	solveCmd.Flags().BoolVar(&csv, "csv", false, "format output as csv")
}
