// Package cmd
package cmd

import (
	"fmt"
	"github.com/ahmedkhaeld/pScan/scan"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io"
	"os"
)

// scanCmd represents the scan command
var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Run a port scan on a list of hosts",
	RunE: func(cmd *cobra.Command, args []string) error {
		hostsFile := viper.GetString("hosts-file")

		ports, err := cmd.Flags().GetIntSlice("ports")
		if err != nil {
			return err
		}

		return scanAction(os.Stdout, hostsFile, ports)
	},
}

func scanAction(out io.Writer, hostsFile string, ports []int) error {
	hl := &scan.HostsList{}
	if err := hl.Load(hostsFile); err != nil {
		return err
	}
	results := scan.Run(hl, ports)
	return printResults(out, results)

}

func printResults(out io.Writer, results []scan.Results) error {
	message := ""
	for _, result := range results {
		message += fmt.Sprintf("%s:\n", result.Host)

		if result.NotFound {
			message += fmt.Sprintf("\tHost not found\n")
			continue
		}
		message += fmt.Sprintln()
		for _, portState := range result.PortStates {
			message += fmt.Sprintf("\tPort %d: %s\n", portState.Port, portState.Open)
		}
		message += fmt.Sprintln()
	}
	_, err := fmt.Fprint(out, message)
	return err
}

func init() {
	rootCmd.AddCommand(scanCmd)

	scanCmd.Flags().IntSliceP("ports", "p", []int{80, 443}, "ports to scan")

}
