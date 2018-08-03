package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func deleteAllFilters(cmd *cobra.Command, args []string) error {
	user, _ := cmd.Flags().GetString("user")
	r, err := svc.Users.Settings.Filters.List(user).Do()
	if err != nil {
		log.Printf("Unable to retrieve filters: %v", err)
		return err
	}

	for _, f := range r.Filter {
		log.Printf("Deleting filter ID %s Criteria %v\n", f.Id, f.Criteria)
		svc.Users.Settings.Filters.Delete(user, f.Id).Do()
	}
	return nil
}

var deleteFiltersCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete all filters",
	RunE:  deleteAllFilters,
}

func uploadFilters(cmd *cobra.Command, args []string) error {
	user, _ := cmd.Flags().GetString("user")
	// Get list of labels
	r, err := svc.Users.Labels.List(user).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve labels: %v", err)
	}
	// TODO: this needn't be a return condition
	if len(r.Labels) == 0 {
		return errors.New("No labels found")
	}
	// Make a map from the "label name" to the "label id"
	labelMap := make(map[string]string)
	for _, l := range r.Labels {
		labelMap[l.Name] = l.Id
	}

	// Read config file
	config, _ := cmd.Flags().GetString("config")
	yamlFile, err := ioutil.ReadFile(config)
	if err != nil {
		log.Printf("%e\n", err)
		return err
	}
	var c GMXConfig
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		log.Printf("%e\n", err)
		return err
	}

	// Replace string label IDs with canonical IDs
	// TODO: assumes all labels already exist
	for _, f := range c.Filters {
		for i, id := range f.Action.AddLabelIds {
			if name, ok := labelMap[id]; ok {
				f.Action.AddLabelIds[i] = name
			}
		}
		for i, id := range f.Action.RemoveLabelIds {
			if name, ok := labelMap[id]; ok {
				f.Action.RemoveLabelIds[i] = name
			}
		}

		fmt.Printf("Creating filter %v Action %v\n", f.Criteria, f.Action)

		// Create the filter
		if _, err := svc.Users.Settings.Filters.Create(user, &f).Do(); err != nil {
			fmt.Printf("Failed with %v\n", err)
		}
	}
	return nil
}

var uploadFiltersCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload filters",
	RunE:  uploadFilters,
}

var filtersCmd = &cobra.Command{
	Use:   "filters",
	Short: "Operations on Filters",
}

func init() {
	filtersCmd.AddCommand(deleteFiltersCmd)
	filtersCmd.AddCommand(uploadFiltersCmd)
	RootCmd.AddCommand(filtersCmd)
}
