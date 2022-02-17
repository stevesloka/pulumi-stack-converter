/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"context"
	"encoding/json"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
	"github.com/pulumi/pulumi/sdk/v3/go/common/apitype"
	"github.com/spf13/cobra"
)

// convertCmd represents the convert command
var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		createStack()
	},
}

var sourceStackName, sourceStackDirectory string

func init() {
	rootCmd.AddCommand(convertCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// convertCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	convertCmd.Flags().StringVar(&sourceStackName, "source-stack-name", "dev", "Name of the local stack to convert")
	convertCmd.Flags().StringVar(&sourceStackDirectory, "source-stack-directory", "", "Location of the of the local stack to convert")

	convertCmd.MarkFlagRequired("source-stack-name")
	convertCmd.MarkFlagRequired("source-stack-directory")
}

func createStack() {
	ctx := context.Background()

	// create a workspace from a local project
	w, _ := auto.NewLocalWorkspace(ctx, auto.WorkDir(sourceStackDirectory))
	stackName := auto.FullyQualifiedStackName("org", "proj", "existing_stack")
	dep, _ := w.ExportStack(ctx, stackName)
	// import/export is backwards compatible, and we must write code specific to the verison we're dealing with.
	if dep.Version != 3 {
		panic("expected deployment version 3")
	}
	var state apitype.DeploymentV3
	_ = json.Unmarshal(dep.Deployment, &state)

	// ... perform edits on the state ...

	// marshal out updated deployment state
	bytes, _ := json.Marshal(state)
	dep.Deployment = bytes
	// import our edited deployment state back to our stack
	_ = w.ImportStack(ctx, stackName, dep)
}
