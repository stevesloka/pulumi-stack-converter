/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/pulumi/pulumi/sdk/v3/go/auto/optup"

	"github.com/pulumi/pulumi/sdk/v3/go/auto"

	"github.com/spf13/cobra"
)

// loadgenCmd represents the loadgen command
var loadgenCmd = &cobra.Command{
	Use:   "loadgen",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		loadgen(cmd, args)
	},
}

var loadgenStackName, loadgenStackDirectory string

func init() {
	rootCmd.AddCommand(loadgenCmd)

	// Here you will define your flags and configuration settings.
	loadgenCmd.Flags().StringVar(&loadgenStackName, "stack-name", "dev", "Name of the stack to run against.")
	loadgenCmd.Flags().StringVar(&loadgenStackDirectory, "stack-directory", "", "Location of the of the local stack.")

	loadgenCmd.MarkFlagRequired("stack-name")
	loadgenCmd.MarkFlagRequired("stack-directory")
}

func loadgen(cmd *cobra.Command, args []string) {
	ctx := context.Background()

	s, err := auto.UpsertStackLocalSource(ctx, loadgenStackName, loadgenStackDirectory)
	if err != nil {
		fmt.Printf("Failed to create or select stack: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Created/Selected stack %q\n", loadgenStackName)

	fmt.Println("Starting refresh")

	_, err = s.Refresh(ctx)
	if err != nil {
		fmt.Printf("Failed to refresh stack: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Refresh succeeded!")

	fmt.Println("Starting update")

	// wire up our update to stream progress to stdout
	stdoutStreamer := optup.ProgressStreams(os.Stdout)

	// run the update
	_, err = s.Up(ctx, stdoutStreamer)
	if err != nil {
		fmt.Printf("Failed to update stack: %v\n\n", err)
		os.Exit(1)
	}
}
