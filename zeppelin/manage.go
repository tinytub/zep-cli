package zeppelin

/*
#cgo CFLAGS: -I ${SRCDIR}/include
#cgo LDFLAGS: -L ${SRCDIR}/lib -lchash -lstdc++

#include "chash.h"
*/
import "C"
import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"os"
	"runtime"
	"strconv"

	"github.com/tinytub/zep-cli/proto/ZPMeta"
)

func ListNode(addrs []string) ([]*ZPMeta.NodeStatus, error) {
	conn, err := NewConn(addrs)
	if err != nil {
		return nil, err
	}
	//conn.mu.Lock()
	data, _ := conn.ListNode()
	conn.RecvDone <- true
	if data.Code.String() != "OK" {
		return nil, errors.New(*data.Msg)
	}
	nodes := data.ListNode.Nodes.Nodes

	return nodes, nil
}

type MetaNodes struct {
	Followers []*Node
	Leader    *Node
}
type Node struct {
	Ip   *string
	Port *int32
}

func ListMeta(addrs []string) (*ZPMeta.MetaNodes, error) {
	conn, err := NewConn(addrs)

	if err != nil {
		return &ZPMeta.MetaNodes{}, err
	}
	//conn.mu.Lock()
	data, _ := conn.ListMeta()
	if data.Code.String() != "OK" {
		return &ZPMeta.MetaNodes{}, errors.New(*data.Msg)
	}
	metas := data.ListMeta.Nodes

	conn.RecvDone <- true
	return metas, nil
}

func ListTable(addrs []string) (*ZPMeta.TableName, error) {
	conn, err := NewConn(addrs)
	if err != nil {
		return &ZPMeta.TableName{}, err
	}

	//conn.mu.Lock()
	data, _ := conn.ListTable()
	if data.Code.String() != "OK" {
		return &ZPMeta.TableName{}, errors.New(*data.Msg)
	}
	tables := data.ListTable.Tables
	conn.RecvDone <- true
	return tables, nil
}

func CreateTable(name string, num int32, addrs []string) {
	conn, err := NewConn(addrs)
	if err != nil {
		//return nil, err
	}

	//conn.mu.Lock()
	data, _ := conn.CreateTable(name, num)
	if data.Code.String() != "OK" {
		fmt.Println(*data.Msg)
		os.Exit(0)
	}
	fmt.Println(data)
	conn.RecvDone <- true
	return
}

func Ping(addrs []string) {
	conn, err := NewConn(addrs)
	if err != nil {
		//return nil, err
	}

	//conn.mu.Lock()
	data, _ := conn.Ping()
	fmt.Println(data)
	if data.Code.String() != "OK" {
		fmt.Println(*data.Msg)
		os.Exit(0)
	}
	fmt.Println(data)
	conn.RecvDone <- true
	return
}

func Set(tablename string, key string, value string, addrs []string) {
	partlocale := locationPartition(tablename, key, addrs)
	fmt.Println(partlocale)
	Nconn, err := NewConn(partlocale)
	if err != nil {
		//return nil, err
	}

	/*
		fmt.Println(Nconn)
		infostats, _ := Nconn.InfoStats(tablename)
		fmt.Println(infostats)
	*/
	v := []byte(value)
	setresp, err := Nconn.Set(tablename, key, v)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(setresp)
}

func Get(tablename string, key string, value string, addrs []string) {
	partlocale := locationPartition(tablename, key, addrs)
	Nconn, err := NewConn(partlocale)
	if err != nil {
		//		return nil, err
	}

	getresp, err := Nconn.Get(tablename, key)
	Nconn.RecvDone <- true
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(getresp)

}

func Space(tablename string, addrs []string) (int64, int64, error) {
	// pull table ---> get partition master info state ---> calculate
	Mconn, err := NewConn(addrs)
	if err != nil {
		return 0, 0, err
	}

	pullResp, err := Mconn.PullTable(tablename)
	if err != nil {
		return 0, 0, err
	}
	Mconn.RecvDone <- true
	pull := pullResp.Pull.GetInfo()[0]
	//	var masters []string
	var used int64 = 0
	var remain int64 = 0

	worker := runtime.NumCPU()
	//fmt.Println(worker)
	jobs := make(chan Jobber, worker)
	results := make(chan Result, len(pull.GetPartitions()))
	dones := make(chan struct{}, worker)

	//go addJob(pull.GetPartitions(), jobs, results, &JobInfoCap{})
	//通过 addjob 不太好搞...只能闭包构造 jobs 管道
	go func() {
		//var pmaddrs map[string]string
		pmaddrs := make(map[string]string)

		for _, partition := range pull.GetPartitions() {
			ip := partition.Master.GetIp()
			port := strconv.Itoa(int(partition.Master.GetPort()))
			pmaddrs[ip] = port
			//jobs <- job{addr, results}
		}
		for ip, port := range pmaddrs {
			jobs <- &JobInfoCap{ip + ":" + port, results}
		}
		close(jobs)
	}()

	for i := 0; i < worker; i++ {
		go doJob(jobs, dones, tablename)
	}
	data := awaitForCloseResult(dones, results, worker)
	for _, d := range data {
		used += d.Used
		remain += d.Remain
	}
	return used, remain, nil

}

func Stats(tablename string, addrs []string) (int64, int32, error) {
	Mconn, err := NewConn(addrs)
	if err != nil {
		return 0, 0, err
	}

	pullResp, err := Mconn.PullTable(tablename)
	if err != nil {
		return 0, 0, err
	}
	Mconn.RecvDone <- true
	pull := pullResp.Pull.GetInfo()[0]

	//	var masters []string

	var query int64
	var qps int32

	worker := runtime.NumCPU()
	//fmt.Println(worker)
	jobs := make(chan Jobber, worker)
	results := make(chan Result, len(pull.GetPartitions()))
	dones := make(chan struct{}, worker)

	go func() {
		pmaddrs := make(map[string]string)
		for _, partition := range pull.GetPartitions() {
			ip := partition.Master.GetIp()
			port := strconv.Itoa(int(partition.Master.GetPort()))
			fmt.Println(ip)
			fmt.Println(port)

			pmaddrs[ip] = port
			//jobs <- job{addr, results}
		}
		for ip, port := range pmaddrs {
			jobs <- &JobInfoStats{ip + ":" + port, results}
		}
		close(jobs)
	}()

	for i := 0; i < worker; i++ {
		go doJob(jobs, dones, tablename)
	}
	data := awaitForCloseResult(dones, results, worker)
	for _, d := range data {
		query += d.Query
		qps += d.QPS
	}
	return query, qps, nil

}

// 多 slave 时通过 goroutine 搞
func addJob(partitions []*ZPMeta.Partitions, jobs chan<- Jobber, results chan<- Result, jobtype interface{}) {
	for _, partition := range partitions {
		addr := partition.Master.GetIp() + ":" + strconv.Itoa(int(partition.Master.GetPort()))
		jobs <- &JobInfoCap{addr, results}
		//jobs <- job{addr, results}
	}
	close(jobs)
}

func doJob(jobs <-chan Jobber, dones chan<- struct{}, tablename string) {
	for job := range jobs {
		job.Do(tablename)
	}
	dones <- struct{}{}
}

type Jobber interface {
	Do(tablename string)
}

type JobInfoCap struct {
	addr   string
	result chan<- Result
}
type JobInfoStats struct {
	addr   string
	result chan<- Result
}
type Result struct {
	Used   int64
	Remain int64
	QPS    int32
	Query  int64
}

func (job *JobInfoCap) Do(tablename string) {
	//addr := partition.Master.GetIp() + ":" + strconv.Itoa(int(partition.Master.GetPort()))
	var used int64
	var remain int64
	Nconn, errN := NewConn([]string{job.addr})
	if errN != nil {
		return
	}

	inforesp, err := Nconn.InfoCapacity(tablename)
	Nconn.RecvDone <- true
	if err != nil {
		return
	}
	infoCap := inforesp.GetInfoCapacity()
	for _, i := range infoCap {
		used += i.GetUsed()
		remain += i.GetRemain()
	}
	job.result <- Result{Used: used, Remain: remain}
}

func (job *JobInfoStats) Do(tablename string) {
	//addr := partition.Master.GetIp() + ":" + strconv.Itoa(int(partition.Master.GetPort()))
	var query int64
	var qps int32
	Nconn, errN := NewConn([]string{job.addr})
	if errN != nil {
		return
	}

	inforesp, err := Nconn.InfoStats(tablename)
	Nconn.RecvDone <- true
	if err != nil {
		return
	}
	infoStats := inforesp.GetInfoStats()
	for _, i := range infoStats {
		fmt.Println("in manges:", i, job.addr)
		query += i.GetTotalQuerys()
		qps += i.GetQps()
	}
	job.result <- Result{Query: query, QPS: qps}
}

func awaitForCloseResult(dones <-chan struct{}, results chan Result, worker int) []Result {
	working := worker
	done := false
	//var totalused int64 = 0
	//var totalremain int64 = 0
	var data []Result
	for {
		select {
		case result := <-results:
			data = append(data, result)
			//totalused += result.Used
			//totalremain += result.Remain
		case <-dones:
			working -= 1
			if working <= 0 {
				done = true
			}
		default:
			if done {
				//fmt.Println("goroutine totalused", totalused)
				//fmt.Println("goroutine totalremain", totalremain)
				return data
			}
		}
	}
}

//
func locationPartition(tablename string, key string, addrs []string) []string {
	Mconn, err := NewConn(addrs)
	if err != nil {
		//return nil, err
	}

	tableinfo, _ := Mconn.PullTable(tablename)
	//./src/zp_table.cc:  int par_num = std::hash<std::string>()(key) % partitions_.size();

	/* 动态链接库的编译方法
	gcc -c chash.cc -std=c++11
	ar rv libchash.a chash.o
	mv libchash.a ../lib
	测试
	g++ -o chash chash.cc -std=c++11
	*/

	partcount := len(tableinfo.Pull.Info[0].Partitions)
	parNum := uint(C.chash(C.CString(key))) % uint(partcount)
	nodemaster := tableinfo.Pull.Info[0].Partitions[parNum].GetMaster()

	Mconn.RecvDone <- true
	var nodeaddrs []string
	nodeaddrs = append(nodeaddrs, nodemaster.GetIp()+":"+strconv.Itoa(int(nodemaster.GetPort())))
	return nodeaddrs

}

func getPartition(conn *Connection, tablename string) int {

	data, _ := conn.PullTable(tablename)

	if data.Code.String() != "OK" {
		fmt.Println(*data.Msg)
		os.Exit(0)
	}

	var pNum int
	for _, part := range data.Pull.Info {
		pNum += len(part.Partitions)
	}
	fmt.Println(pNum)
	return pNum
}

//func getBytes(key interface{}) ([]byte, error) {
func getBytes(key string) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(key)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
