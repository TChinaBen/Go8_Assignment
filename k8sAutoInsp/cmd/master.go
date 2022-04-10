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
	"github.com/spf13/cobra"
	"k8sAutoInsp/check"
)

// masterCmd represents the master command
var masterCmd = &cobra.Command{
	Use:   "master",
	Short: "Run Kubernetes benchmark checks from the master.yaml file",
	Long: `Run Kubernetes benchmark checks from the master.yaml file in cfg/<version>`,
	Run: func(cmd *cobra.Command, args []string) {
		// 获取当前kubernetes版本号
		bv := "cis-1.20"
		// 获取文件路径
		filename := loadConfig(check.MASTER,bv)
		runChecks(check.MASTER,filename)
	},
}

func init() {

	masterCmd.PersistentFlags().StringVarP(&masterFile,
		"file",
		"f",
		"/master.yaml",
		"Alternative YAML file for master checks")
	rootCmd.AddCommand(masterCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// masterCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// masterCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
