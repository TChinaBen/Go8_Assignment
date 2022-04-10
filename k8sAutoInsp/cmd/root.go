/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"flag"
	"fmt"
	"github.com/golang/glog"
	"github.com/spf13/cobra"
	"k8sAutoInsp/check"
	"os"

	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   os.Args[0],
	Short: "Run CIS Benchmarks checks against a Kubernetes deployment",
	Long: `This tool runs the CIS Kubernetes Benchmark(https://www.cisecurity.org/benchmark/kubernetes)`,
	Run: func(cmd *cobra.Command,args []string){
		bv := "cis-1.20"
		if isMaster(){
			runChecks(check.MASTER,loadConfig(check.MASTER,bv))
			valid,err := ValidTargets(bv,[]string{string(check.CONTROLPLANE)}, viper.GetViper())
			if err != nil{
				err.Error()
				return
			}
			if valid{
				runChecks(check.CONTROLPLANE,loadConfig(check.CONTROLPLANE,bv))
			}
		}else{
			glog.V(1).Info("== Skipping master checks ==")
		}
		//valid,err := ValidTargets(bv,[]string{string(check.ETCD)},viper.GetViper())
		//if err != nil{
		//	err.Error()
		//	return
		//}
		//if valid && isEtcd(){
		//	glog.V(1).Info("== Running etcd checks ==")
		//	runChecks(check.ETCD,loadConfig(check.ETCD,bv))
		//}else{
		//	glog.V(1).Info("== Running etcd checks ==")
		//}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	flag.CommandLine.Parse([]string{})
	if err := rootCmd.Execute();err != nil{
		fmt.Println(err)
		glog.Flush()
		os.Exit(1)
	}
	glog.Flush()

}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.k8sAutoInsp.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName("config")
		viper.AddConfigPath(cfgDir)
	}
	viper.SetEnvPrefix(envVarsPrefix)
	viper.AutomaticEnv() // read in environment variables that match
	if err := viper.ReadInConfig();err != nil{
		if _,ok := err.(viper.ConfigFileNotFoundError);ok{
			configFileError = err
		}else{
			os.Exit(1)
		}
	}
}
