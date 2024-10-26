package main

import (
	"fmt"

	"github.com/6thfdwp/prober/cmd"
	"github.com/spf13/cobra"
)

const (
// sub1 = "kallangur-qld-4503"
// sub2 = "petrie-qld-4502"
// mango-hill-qld-4509, lawnton-qld-4501, warner-qld-4500,griffin-qld-4503
// daisy-hill-qld-4127, heritage-park-qld-4118
// DomainUrl = "https://www.domain.com.au/suburb-profile/regents-park-qld-4118"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "prober",
		Short: "Prober CLI tool",
		Long:  `Prober CLI tool to extract key info.`,
	}

	rootCmd.AddCommand(cmd.NewSuburbCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		return
	}
}
