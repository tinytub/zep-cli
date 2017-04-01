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
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tinytub/zep-cli/s3core"
)

// monitorCmd represents the monitor command
var monitorCmd = &cobra.Command{
	Use:   "monitor",
	Short: "zep s3 gateway tools",
	Long: `A tool for zep s3 gateway
for normal test and bench`,
}

var BoolCode_value = map[bool]int32{
	true:  0,
	false: 1,
}

var StringCode_value = map[string]int32{
	"OK": 0,
	"":   1,
}

var s3MoSetOBJ = &cobra.Command{
	Use:   "set",
	Short: "zep s3 set key monit",
	Long: `A tool for zep s3 gateway
		for normal set test`,
	Run: func(cmd *cobra.Command, args []string) {
		edp, acckey, sec := checkRegion(region)
		svc := s3core.NewClient(edp, acckey, sec)
		//fmt.Println(mokey, value)
		res, _ := s3core.SetOBJ(svc, bucket, mokey, value, filename)
		fmt.Printf("S3Set=%v", BoolCode_value[res])
	},
}

var s3MoGetOBJ = &cobra.Command{
	Use:   "get",
	Short: "zep s3 get key",
	Long: `A tool for zep s3 gateway
		for normal set test`,
	Run: func(cmd *cobra.Command, args []string) {
		edp, acckey, sec := checkRegion(region)
		svc := s3core.NewClient(edp, acckey, sec)
		res, _ := s3core.GetOBJ(svc, bucket, mokey, "stdout")
		//s3core.DelOBJ(svc, bucket, mokey)
		fmt.Printf("S3Get=%v", StringCode_value[res])
	},
}

var s3MoSetGetDelOBJ = &cobra.Command{
	Use:   "sgd",
	Short: "zep s3 set get del key",
	Long: `A tool for zep s3 gateway
		for normal set test`,
	Run: func(cmd *cobra.Command, args []string) {
		edp, acckey, sec := checkRegion(region)
		svc := s3core.NewClient(edp, acckey, sec)
		setres, _ := s3core.SetOBJ(svc, bucket, mokey, value, filename)
		//fmt.Printf("S3Set=%v", BoolCode_value[setres])
		getres, _ := s3core.GetOBJ(svc, bucket, mokey, "stdout")
		//s3core.DelOBJ(svc, bucket, mokey)
		//fmt.Printf("S3Get=%v", StringCode_value[getres])
		delres, _ := s3core.DelOBJ(svc, bucket, mokey)
		fmt.Printf("S3Set=%v&S3Get=%v&S3Del=%v", BoolCode_value[setres], StringCode_value[getres], BoolCode_value[delres])

	},
}

var mokey string

func init() {
	t := time.Now()
	moKey := fmt.Sprintf("monit-%02d-%02d", t.Hour(), t.Minute())

	s3MoSetOBJ.Flags().StringVarP(&bucket, "bucket", "b", "monitor", "bucket name")
	s3MoSetOBJ.Flags().StringVarP(&mokey, "key", "k", moKey, "which key")
	s3MoSetOBJ.Flags().StringVarP(&value, "value", "v", "OK", "which value")
	s3MoSetOBJ.Flags().StringVarP(&filename, "f", "f", "", "filename which you want upload")
	//	s3MoSetOBJ.Flags().StringVar(&region, "region", "", "s3 region")

	s3MoGetOBJ.Flags().StringVarP(&bucket, "bucket", "b", "monitor", "bucket name")
	s3MoGetOBJ.Flags().StringVarP(&mokey, "key", "k", moKey, "which key")
	//	s3MoGetOBJ.Flags().StringVar(&region, "region", "", "s3 region")

	s3MoSetGetDelOBJ.Flags().StringVarP(&bucket, "bucket", "b", "monitor", "bucket name")
	s3MoSetGetDelOBJ.Flags().StringVarP(&mokey, "key", "k", moKey, "which key")
	s3MoSetGetDelOBJ.Flags().StringVarP(&value, "value", "v", "OK", "which value")
	s3MoSetGetDelOBJ.Flags().StringVarP(&filename, "f", "f", "", "filename which you want upload")
	//	s3MoSetGetDelOBJ.Flags().StringVar(&region, "region", "", "s3 region")

	monitorCmd.PersistentFlags().StringVar(&region, "region", "", "s3 region")

	monitorCmd.AddCommand(s3MoSetOBJ)
	monitorCmd.AddCommand(s3MoGetOBJ)
	monitorCmd.AddCommand(s3MoSetGetDelOBJ)
	RootCmd.AddCommand(monitorCmd)

	viper.Get("s3")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// monitorCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// monitorCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
