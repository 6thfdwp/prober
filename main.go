package main

import (
	"fmt"

	"github.com/6thfdwp/prober/cmd"
	"github.com/spf13/cobra"
)

const (
// sub1 = kallangur-qld-4503,petrie-qld-4502",warner-qld-4500,griffin-qld-4503
// daisy-hill-qld-4127, marsden-qld-4132
// springfield-qld-4300
// upper-coomera-qld-4209,pacific-pines-qld-4211,arundel-qld-4214

// coorparoo-qld-4151,
// carina-qld-4152,carina-heights-qld-4152,cannon-hill-qld-4170,holland-park-qld-4121,holland-park-west-qld-4121
// carina-4152/florence-st/16, carina-4152/eleanor-st/46, carina-4152/lunga-st/67,4
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
