package main

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "zwallet",
	Short: "Use Zwallet to store, send and execute smart contract on 0Chain platform",
	Long: `Use Zwallet to store, send and execute smart contract on 0Chain platform.
			Complete documentation is available at https://docs.zus.network/guides/zwallet-cli`,
	// PersistentPreRun: ,
}

func main() {

}
