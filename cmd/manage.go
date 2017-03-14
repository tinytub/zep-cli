// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tinytub/zep-cli/zeppelin"
)

// manageCmd represents the manage command
var manageCmd = &cobra.Command{
	Use:   "manage",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	/*
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: Work your own magic here
			fmt.Println("manage called")
		},
	*/
}

var cmdListNode = &cobra.Command{
	Use:   "listnode",
	Short: "list nodes",
	Long:  "list all zeppelin nodes",
	Run: func(cmd *cobra.Command, args []string) {
		zeppelin.ListNode()
	},
}

var cmdListMeta = &cobra.Command{
	Use:   "listmeta",
	Short: "list Metas",
	Long:  "list all zeppelin Metas",
	Run: func(cmd *cobra.Command, args []string) {
		zeppelin.ListMeta()
	},
}

var cmdListTable = &cobra.Command{
	Use:   "listtable",
	Short: "list table",
	Long:  "list all zeppelin tables",
	Run: func(cmd *cobra.Command, args []string) {
		zeppelin.ListTable()
	},
}

var cmdCreateTable = &cobra.Command{
	Use:   "createtable",
	Short: "create table",
	Long:  "create a table",
	Run: func(cmd *cobra.Command, args []string) {
		zeppelin.CreateTable(tname, tnum)
	},
}

var cmdSet = &cobra.Command{
	Use:   "set",
	Short: "set key to table",
	Long:  "set key to table",
	Run: func(cmd *cobra.Command, args []string) {
		zeppelin.Set(tname, ukey, uvalue)
	},
}

var cmdGet = &cobra.Command{
	Use:   "get",
	Short: "get key from table",
	Long:  "get key from table",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

var tname string
var tnum int32
var ukey string
var uvalue string

func init() {
	//cmdListNode.Flags()
	RootCmd.AddCommand(manageCmd)

	manageCmd.AddCommand(cmdListNode)
	manageCmd.AddCommand(cmdListMeta)
	manageCmd.AddCommand(cmdListTable)

	cmdCreateTable.Flags().StringVarP(&tname, "name", "n", "", "table name")
	cmdCreateTable.Flags().Int32VarP(&tnum, "num", "N", 10, "table's partition num")
	manageCmd.AddCommand(cmdCreateTable)

	cmdSet.Flags().StringVarP(&tname, "name", "t", "", "table name")
	cmdSet.Flags().StringVarP(&ukey, "key", "k", "", "key")
	cmdSet.Flags().StringVarP(&uvalue, "value", "v", "", "value")
	manageCmd.AddCommand(cmdSet)
}
