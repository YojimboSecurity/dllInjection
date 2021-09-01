package cmd

import (
	"fmt"
	"os"
	"github.com/YojimboSecurity/dllInjection/src"
	"github.com/spf13/cobra"
)

var (
	cfgFile string
	pid int16
	dll string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "DLLInjection",
	Short: "Inject a DLL into a process",
	Long: `Inject a DLL into a process`,
	Run: func(cmd *cobra.Command, args []string) { 
		src.DLLInjection(pid, dll)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize()

	rootCmd.Flags().Int16VarP(&pid, "pid", "p", 0, "PID to inject DLL into")
	rootCmd.Flags().StringVarP(&dll, "dll", "d", "", "DLL to inject")
}
