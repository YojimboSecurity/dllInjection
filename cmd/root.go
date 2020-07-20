// Copyright © 2020 David Johnson <david@yojimbosecurity.ninja>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"
	"log"
	"os"

	"git.yojimbosecurity.com/dllInjection/src"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var pid int16
var dllPath string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dllInjection",
	Short: "DLL Injection POC",
	Long: `I wrote this to test Sysmon Create Remote Thread EventID 8

To test Sysmon EventID 8 copy the example below to a file such as sysmon_config.xml 
and load it Sysmon.exe -c sysmon_config.xml

<Sysmon schemaversion="4.1">
	<EventFiltering>
	   <CreateRemoteThread onmatch="include">
		   <StartFunction name="technique_id=T1055,technique_name=Process Injection" condition="contains">LoadLibrary</StartFunction>
		   <TargetImage name="technique_id=T1055,technique_name=Process Injection" condition="is">C:\Windows\System32\rundll32.exe</TargetImage>
		   <TargetImage name="technique_id=T1055,technique_name=Process Injection" condition="is">C:\Windows\System32\svchost.exe</TargetImage>
		   <TargetImage name="technique_id=T1055,technique_name=Process Injection" condition="is">C:\Windows\System32\sysmon.exe</TargetImage>
	   </CreateRemoteThread>
   </EventFiltering>
 </Sysmon>
`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		if pid == 0 {
			log.Fatalln("[!] Missing PID")
		}
		if dllPath == "" {
			log.Fatalln("[!] Missing DLL Path")
		} 
		src.DLLInjection(pid, dllPath)
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
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.dllInjection.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.Flags().StringVarP(&dllPath, "dll-path", "d", "", "Path to DLL to be injected")
	rootCmd.Flags().Int16VarP(&pid, "pid", "p", 0, "PID of process to inject DLL into")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".dllInjection" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".dllInjection")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
