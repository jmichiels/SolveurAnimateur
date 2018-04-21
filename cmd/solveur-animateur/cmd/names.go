package cmd

import (
	"log"

	"github.com/jmichiels/SolveurAnimateur"
	"github.com/spf13/cobra"
)

// personsCmd represents the persons command
var personsCmd = &cobra.Command{
	Use: "persons",
	Run: func(cmd *cobra.Command, args []string) {
		if err := solveuranimateur.GeneratePersonsDefinitions(); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	generateCmd.AddCommand(personsCmd)
}
