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
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tinytub/zep-cli/zeppelin"
)

var status = map[int32]string{
	0: "up",
	1: "down",
}

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
		meta := checkZepRegionNGetMeta(region)
		nodes, err := zeppelin.ListNode(meta)
		if err != nil {
			fmt.Println(err)
			return
		}
		for _, node := range nodes {
			fmt.Printf("IP: %s, Port: %d, Status: %s\n", *node.Node.Ip, *node.Node.Port, status[*node.Status])
		}
	},
}

var cmdListMeta = &cobra.Command{
	Use:   "listmeta",
	Short: "list Metas",
	Long:  "list all zeppelin Metas",
	Run: func(cmd *cobra.Command, args []string) {
		meta := checkZepRegionNGetMeta(region)

		metas, err := zeppelin.ListMeta(meta)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Leader:")
		fmt.Printf("IP: %s, Port: %d\n", *metas.Leader.Ip, *metas.Leader.Port)
		fmt.Println("Followers:")
		for _, follower := range metas.Followers {
			fmt.Printf("IP: %s, Port: %d\n", *follower.Ip, *follower.Port)
		}
	},
}

var cmdListTable = &cobra.Command{
	Use:   "listtable",
	Short: "list table",
	Long:  "list all zeppelin tables",
	Run: func(cmd *cobra.Command, args []string) {
		meta := checkZepRegionNGetMeta(region)

		tablelist, err := zeppelin.ListTable(meta)
		if err != nil {
			fmt.Println(err)
			return
		}
		for _, table := range tablelist.Name {
			fmt.Println(table)
		}

	},
}

var cmdCreateTable = &cobra.Command{
	Use:   "createtable",
	Short: "create table",
	Long:  "create a table",
	Run: func(cmd *cobra.Command, args []string) {
		meta := checkZepRegionNGetMeta(region)
		zeppelin.CreateTable(tname, tnum, meta)
	},
}

var cmdSet = &cobra.Command{
	Use:   "set",
	Short: "set key to table",
	Long:  "set key to table",
	Run: func(cmd *cobra.Command, args []string) {
		meta := checkZepRegionNGetMeta(region)

		zeppelin.Set(tname, ukey, uvalue, meta)
	},
}

var cmdGet = &cobra.Command{
	Use:   "get",
	Short: "get key from table",
	Long:  "get key from table",
	Run: func(cmd *cobra.Command, args []string) {
		meta := checkZepRegionNGetMeta(region)

		zeppelin.Get(tname, ukey, uvalue, meta)
	},
}

var cmdSpace = &cobra.Command{
	Use:   "space",
	Short: "space usage and remain for specified table",
	Run: func(cmd *cobra.Command, args []string) {
		meta := checkZepRegionNGetMeta(region)

		ptable, err := zeppelin.PullTable(tname, meta)
		if err != nil {
			fmt.Println(err)
			return
		}
		used, remain, err := ptable.Space(tname, meta)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Table: %s, Used: %d bytes, Ramain: %d bytes\n", tname, used, remain)

	},
}

var cmdStats = &cobra.Command{
	Use:   "stats",
	Short: "QPS and total Query for specified table",
	Run: func(cmd *cobra.Command, args []string) {
		meta := checkZepRegionNGetMeta(region)
		ptable, err := zeppelin.PullTable(tname, meta)
		if err != nil {
			fmt.Println(err)
			return
		}
		query, qps, err := ptable.Stats(tname, meta)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Table: %s, Query: %d, QPS: %d \n", tname, query, qps)

	},
}
var cmdOffset = &cobra.Command{
	Use:   "offset",
	Short: "offset specified table",
	Run: func(cmd *cobra.Command, args []string) {
		meta := checkZepRegionNGetMeta(region)

		ptable, err := zeppelin.PullTable(tname, meta)
		if err != nil {
			fmt.Println(err)
			return
		}
		unsyncoffset, err := ptable.Offset(tname, meta)
		if err != nil {
			fmt.Println(err)
			return
		}
		for pid, slave := range unsyncoffset {
			for _, offset := range slave {
				if offset != nil {
					fmt.Printf("unsynced partition id: %d, addr: %s, offset: %d\n", pid, offset.Addr, offset.Offset)
				}
			}
		}

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

	cmdCreateTable.Flags().StringVarP(&tname, "name", "n", "test", "table name")
	cmdCreateTable.Flags().Int32VarP(&tnum, "num", "N", 10, "table's partition num")
	//	cmdCreateTable.Flags().StringVar(&region, "region", "", "zep region")
	manageCmd.AddCommand(cmdCreateTable)

	cmdSet.Flags().StringVarP(&tname, "name", "t", "", "table name")
	cmdSet.Flags().StringVarP(&ukey, "key", "k", "", "key")
	cmdSet.Flags().StringVarP(&uvalue, "value", "v", "", "value")
	//	cmdSet.Flags().StringVar(&region, "region", "", "zep region")
	manageCmd.AddCommand(cmdSet)

	cmdGet.Flags().StringVarP(&tname, "name", "t", "", "table name")
	cmdGet.Flags().StringVarP(&ukey, "key", "k", "", "key")
	cmdGet.Flags().StringVarP(&uvalue, "value", "v", "", "value")
	//	cmdGet.Flags().StringVar(&region, "region", "", "zep region")
	manageCmd.AddCommand(cmdGet)

	cmdSpace.Flags().StringVarP(&tname, "name", "t", "", "table name")
	manageCmd.AddCommand(cmdSpace)

	cmdStats.Flags().StringVarP(&tname, "name", "t", "", "table name")
	manageCmd.AddCommand(cmdStats)

	cmdOffset.Flags().StringVarP(&tname, "name", "t", "", "table name")
	manageCmd.AddCommand(cmdOffset)

	manageCmd.PersistentFlags().StringVar(&region, "region", "", "zep region")

	//	cmdListNode.Flags().StringVar(&region, "region", "", "zep region")

	//	cmdListMeta.Flags().StringVar(&region, "region", "", "zep region")

	//	cmdListTable.Flags().StringVar(&region, "region", "", "zep region")

}

func getZepAllRegion() []string {
	conf := viper.Get("zep")
	var regionlist []string
	for region, _ := range conf.(map[string]interface{}) {
		regionlist = append(regionlist, region)
	}
	return regionlist
}

func checkZepRegionNGetMeta(region string) []string {
	//var meta []string
	path := fmt.Sprintf("zep.%s.meta", region)
	conf := viper.GetStringSlice(path)
	if conf == nil {
		ListZepRegion()
		os.Exit(0)
		//return nil, nil, nil
	}
	return conf
}
