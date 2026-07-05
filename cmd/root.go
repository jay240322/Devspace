package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootcmd = &cobra.Command{
	Use:   "devspace",
	Short: "Devspace is a Internal Developer Platform CLI",
	Long: "A high-performance platform engineering tool CLI tool to spin up secure, multi-cloud micro-services in seconds",
}

func Execute(){
	if err := rootcmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}