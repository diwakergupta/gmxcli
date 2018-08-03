package cmd

import (
	"log"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

func listLabels(cmd *cobra.Command, args []string) error {
	user, _ := cmd.Flags().GetString("user")
	// Get list of labels
	r, err := svc.Users.Labels.List(user).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve labels: %v", err)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Id", "Name"})
	for _, l := range r.Labels {
		table.Append([]string{l.Id, l.Name})
	}
	table.Render()
	return nil
}

var listLabelsCmd = &cobra.Command{
	Use:   "list",
	Short: "List all labels",
	RunE:  listLabels,
}

var labelsCmd = &cobra.Command{
	Use:   "labels",
	Short: "Operations on Labels",
}

func init() {
	labelsCmd.AddCommand(listLabelsCmd)
	RootCmd.AddCommand(labelsCmd)
}
