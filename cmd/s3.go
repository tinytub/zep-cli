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
	"runtime"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tinytub/zep-cli/s3core"
)

// s3Cmd represents the s3 command
var s3Cmd = &cobra.Command{
	Use:   "s3",
	Short: "zep s3 gateway tools",
	Long: `A tool for zep s3 gateway
for normal test and bench`,
}

var s3ListBucket = &cobra.Command{
	Use:   "lb",
	Short: "list all bucket",
	Run: func(cmd *cobra.Command, args []string) {
		edp, acckey, sec := checkRegion(region)
		svc := s3core.NewClient(edp, acckey, sec)
		s3core.ListBucket(svc)
	},
}

var s3BenchBucket = &cobra.Command{
	Use:   "benchbk",
	Short: "zep s3 bench bucket",
	Long: `A tool for zep s3 gateway
		for normal test and bench`,
	Run: func(cmd *cobra.Command, args []string) {
		edp, acckey, sec := checkRegion(region)
		svc := s3core.NewClient(edp, acckey, sec)
		for i := 0; i < runs; i++ {
			s3core.CreateBucket(svc, fmt.Sprintf("test%d", i))
		}
	},
}

var s3CreateBucket = &cobra.Command{
	Use:   "createbk",
	Short: "zep s3 create bucket",
	Long: `A tool for zep s3 gateway
		for normal test create bucket`,
	Run: func(cmd *cobra.Command, args []string) {
		edp, acckey, sec := checkRegion(region)
		svc := s3core.NewClient(edp, acckey, sec)
		s3core.CreateBucket(svc, bucket)
	},
}

var s3DeleteBucket = &cobra.Command{
	Use:   "delbk",
	Short: "zep s3 delete bucket",
	Long: `A tool for zep s3 gateway
		for normal test create bucket`,
	Run: func(cmd *cobra.Command, args []string) {
		edp, acckey, sec := checkRegion(region)
		svc := s3core.NewClient(edp, acckey, sec)
		s3core.DeleteBucket(svc, bucket)
	},
}

var s3SetOBJ = &cobra.Command{
	Use:   "set",
	Short: "zep s3 set key",
	Long: `A tool for zep s3 gateway
		for normal set test`,
	Run: func(cmd *cobra.Command, args []string) {
		edp, acckey, sec := checkRegion(region)
		svc := s3core.NewClient(edp, acckey, sec)
		_, err := s3core.SetOBJ(svc, bucket, key, value, filename)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("set key ok")
	},
}

var s3GetOBJ = &cobra.Command{
	Use:   "get",
	Short: "zep s3 get key",
	Long: `A tool for zep s3 gateway
		for normal get test`,
	Run: func(cmd *cobra.Command, args []string) {
		edp, acckey, sec := checkRegion(region)
		svc := s3core.NewClient(edp, acckey, sec)
		res, err := s3core.GetOBJ(svc, bucket, key, output)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(res)

	},
}

var s3ListOBJ = &cobra.Command{
	Use:   "listobj",
	Short: "zep s3 list obj",
	Long: `A tool for zep s3 gateway
		for normal list obj`,
	Run: func(cmd *cobra.Command, args []string) {
		edp, acckey, sec := checkRegion(region)
		svc := s3core.NewClient(edp, acckey, sec)
		s3core.ListOBJ(svc, bucket)
	},
}

var s3DelOBJ = &cobra.Command{
	Use:   "delobj",
	Short: "zep s3 del obj",
	Long: `A tool for zep s3 gateway
		for normal list obj`,
	Run: func(cmd *cobra.Command, args []string) {
		edp, acckey, sec := checkRegion(region)
		svc := s3core.NewClient(edp, acckey, sec)
		s3core.DelOBJ(svc, bucket, key)
	},
}

var s3HeadOBJ = &cobra.Command{
	Use:   "headobj",
	Short: "zep s3 head obj",
	Long: `A tool for zep s3 gateway
		for normal list obj`,
	Run: func(cmd *cobra.Command, args []string) {
		edp, acckey, sec := checkRegion(region)
		svc := s3core.NewClient(edp, acckey, sec)
		s3core.HeadOBJ(svc, bucket, key)
	},
}

var s3Test = &cobra.Command{
	Use:   "test",
	Short: "zep s3 test",
	Long: `A tool for zep s3 gateway
		for normal test`,
	Run: func(cmd *cobra.Command, args []string) {
		regionlist := getAllRegion()
		if region != "" {
			regionlist = []string{region}
		}

		runtime.GOMAXPROCS(4)

		doRun := make(chan string, 1)
		for i := 0; i < con; i++ {
			go roundTest(regionlist, doRun, i)
			//fmt.Println("ttt")
		}
		<-doRun

	},
}

func roundTest(regionlist []string, doRun chan string, con int) {
	timeout := make(<-chan time.Time)
	tick := make(<-chan time.Time)
	one := make(chan int, 1)
	if utime > 0 {
		tick = time.Tick(1000 * time.Millisecond)
		timeout = time.After(time.Duration(utime) * time.Second)
	} else {
		one <- utime
	}
	fmt.Println(tekey)
	for {
		select {

		case <-tick:
			t := time.Now()
			ntekey := fmt.Sprintf("t-%02d-%02d-con%d", t.Hour(), t.Minute(), con)
			nbucket := fmt.Sprintf("tb-%02d-%02d-con%d", t.Hour(), t.Minute(), con)
			fmt.Println("BUCKET!!!!!!!", bucket)
			for _, r := range regionlist {
				fmt.Printf("-------\n")
				fmt.Printf("checking region: %s\n", r)
				edp, acckey, sec := checkRegion(r)
				svc := s3core.NewClient(edp, acckey, sec)
				s3core.ListBucket(svc)
				s3core.CreateBucket(svc, nbucket)
				s3core.SetOBJ(svc, nbucket, ntekey, value, filename)
				s3core.GetOBJ(svc, nbucket, ntekey, output)
				s3core.ListOBJ(svc, nbucket)
				s3core.DelOBJ(svc, nbucket, ntekey)
				//s3core.DeleteBucket(svc, nbucket)
			}
		case <-timeout:
			fmt.Println("time out 5 second")
			doRun <- "done"
			return
			os.Exit(1)
		case <-one:
			fmt.Println("default!!!!!!!!!!!!!!!!!")
			for _, r := range regionlist {
				fmt.Printf("-------\n")
				fmt.Printf("checking region: %s\n", r)
				edp, acckey, sec := checkRegion(r)
				svc := s3core.NewClient(edp, acckey, sec)

				s3core.ListBucket(svc)

				s3core.CreateBucket(svc, bucket)

				_, err := s3core.SetOBJ(svc, bucket, tekey, value, filename)
				if err != nil {
					fmt.Println(err)
				} else {
					fmt.Println("set ok")
				}

				res, errG := s3core.GetOBJ(svc, bucket, tekey, output)
				if errG != nil {
					fmt.Println(errG)
				} else {
					fmt.Println(res)
				}

				s3core.ListOBJ(svc, bucket)
				_, errD := s3core.DelOBJ(svc, bucket, tekey)
				if errD != nil {
					fmt.Println(errD)
				} else {
					fmt.Println("delete ok")
				}

				//s3core.DeleteBucket(svc, bucket)
			}
			doRun <- "done"
			return
		}
	}
}

var (
	endpoint string
	runs     int
	region   string
	bucket   string
	key      string
	value    string
	filename string
	output   string
	tekey    string
	utime    int
	con      int
)

func init() {
	t := time.Now()
	teKey := fmt.Sprintf("t-%02d-%02d", t.Hour(), t.Minute())

	RootCmd.AddCommand(s3Cmd)

	s3BenchBucket.Flags().StringVarP(&endpoint, "endpoint", "e", "", "s3 endpoint")
	s3BenchBucket.Flags().IntVarP(&runs, "runs", "r", 3, "the number of times to run each test")
	s3BenchBucket.Flags().StringVar(&region, "region", "", "s3 region")

	s3ListBucket.Flags().StringVar(&region, "region", "", "s3 region")

	s3CreateBucket.Flags().StringVarP(&bucket, "bucket", "b", "monitor", "bucket name")
	s3CreateBucket.Flags().StringVar(&region, "region", "", "s3 region")

	s3DeleteBucket.Flags().StringVarP(&bucket, "bucket", "b", "monitor", "bucket name")
	s3DeleteBucket.Flags().StringVar(&region, "region", "", "s3 region")

	s3SetOBJ.Flags().StringVarP(&bucket, "bucket", "b", "monitor", "bucket name")
	s3SetOBJ.Flags().StringVarP(&key, "key", "k", "monit", "which key")
	s3SetOBJ.Flags().StringVarP(&value, "value", "v", "OK", "which value")
	s3SetOBJ.Flags().StringVarP(&filename, "f", "f", "", "filename which you want upload")
	s3SetOBJ.Flags().StringVar(&region, "region", "", "s3 region")

	s3GetOBJ.Flags().StringVarP(&bucket, "bucket", "b", "monitor", "bucket name")
	s3GetOBJ.Flags().StringVarP(&key, "key", "k", "monit", "which key")
	s3GetOBJ.Flags().StringVarP(&output, "output", "o", "stdout", "filename which you want download file to save")
	s3GetOBJ.Flags().StringVar(&region, "region", "", "s3 region")

	s3ListOBJ.Flags().StringVarP(&bucket, "bucket", "b", "monitor", "bucket name")
	s3ListOBJ.Flags().StringVar(&region, "region", "", "s3 region")

	s3DelOBJ.Flags().StringVarP(&bucket, "bucket", "b", "monitor", "bucket name")
	s3DelOBJ.Flags().StringVarP(&key, "key", "k", "monit", "which key")
	s3DelOBJ.Flags().StringVar(&region, "region", "", "s3 region")

	s3HeadOBJ.Flags().StringVarP(&bucket, "bucket", "b", "monitor", "bucket name")
	s3HeadOBJ.Flags().StringVarP(&key, "key", "k", "monit", "which key")
	s3HeadOBJ.Flags().StringVar(&region, "region", "", "s3 region")

	s3Test.Flags().StringVarP(&bucket, "bucket", "b", "monitor", "bucket name")
	s3Test.Flags().StringVarP(&tekey, "key", "k", teKey, "which key")
	s3Test.Flags().StringVarP(&value, "value", "v", "OK", "which value")
	s3Test.Flags().StringVarP(&filename, "f", "f", "", "filename which you want upload")
	s3Test.Flags().StringVarP(&output, "output", "o", "stdout", "filename which you want download file to save")
	s3Test.Flags().StringVar(&region, "region", "", "s3 region")
	s3Test.Flags().IntVarP(&utime, "timeout", "t", 0, "test timeout")
	s3Test.Flags().IntVarP(&con, "concurrency", "c", 1, "test concurrency")

	s3Cmd.AddCommand(s3BenchBucket)
	s3Cmd.AddCommand(s3ListBucket)
	s3Cmd.AddCommand(s3SetOBJ)
	s3Cmd.AddCommand(s3GetOBJ)
	s3Cmd.AddCommand(s3ListOBJ)
	s3Cmd.AddCommand(s3DelOBJ)
	s3Cmd.AddCommand(s3HeadOBJ)
	s3Cmd.AddCommand(s3CreateBucket)
	s3Cmd.AddCommand(s3DeleteBucket)
	s3Cmd.AddCommand(s3Test)
	//s3Cmd.Flags().StringVar(&region, "region", "", "s3 region")

	viper.Get("s3")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// s3Cmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// s3Cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

func getAllRegion() []string {
	conf := viper.Get("s3")
	var regionlist []string
	for region, _ := range conf.(map[string]interface{}) {
		regionlist = append(regionlist, region)
	}
	return regionlist
}

func checkRegion(region string) (string, string, string) {
	path := fmt.Sprintf("s3.%s", region)
	conf := viper.Get(path)
	if conf == nil {
		ListS3Region()
		os.Exit(0)
		//return nil, nil, nil
	}
	edp := viper.Get(path + ".endpoint").(string)
	key := viper.Get(path + ".key").(string)
	sec := viper.Get(path + ".secret").(string)
	return edp, key, sec
}
