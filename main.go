package main

import (
	"fmt"

	"github.com/6thfdwp/prober/cmd"
	"github.com/spf13/cobra"
)

const (
// sub1 = "kallangur-qld-4503"
// sub2 = "petrie-qld-4502"
// warner-qld-4500,griffin-qld-4503
// daisy-hill-qld-4127,

// ripley-qld-4306,redbank-plains-qld-4301, heritage-park-qld-4118
// coorparoo-qld-4151,cannon-hill-qld-4170,carina-qld-4152
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "prober",
		Short: "Prober CLI tool",
		Long:  `Prober CLI tool to extract housing key info for suburbs and streets.`,
	}

	rootCmd.AddCommand(cmd.NewSuburbCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		return
	}
}
