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
	"reflect"
	"runtime"
	"strconv"

	"github.com/tinytub/zep-cli/proto/ZPMeta"
	"github.com/tinytub/zep-cli/proto/client"
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

type PTable struct {
	pull *ZPMeta.Table
}

func PullTable(tablename string, addrs []string) (PTable, error) {
	var ptable PTable
	Mconn, err := NewConn(addrs)
	if err != nil {
		return ptable, err
	}

	pullResp, err := Mconn.PullTable(tablename)
	if err != nil {
		return ptable, err
	}
	Mconn.RecvDone <- true
	ptable.pull = pullResp.Pull.GetInfo()[0]
	return ptable, nil
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

func (ptable *PTable) Space(tablename string, addrs []string) (int64, int64, error) {
	// pull table ---> get partition master info state ---> calculate
	/*
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
	*/
	pull := ptable.pull
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
		pmaddrs := make(map[string]int)

		for _, partition := range pull.GetPartitions() {
			ip := partition.Master.GetIp()
			port := strconv.Itoa(int(partition.Master.GetPort()))
			pmaddrs[ip+":"+port] = 0
			//jobs <- job{addr, results}
		}
		for addr, _ := range pmaddrs {
			jobs <- &JobInfoCap{addr, results}
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

func (ptable *PTable) Stats(tablename string, addrs []string) (int64, int32, error) {
	/*
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
	*/
	pull := ptable.pull

	//	var masters []string

	var query int64
	var qps int32

	worker := runtime.NumCPU()
	//fmt.Println(worker)
	jobs := make(chan Jobber, worker)
	results := make(chan Result, len(pull.GetPartitions()))
	dones := make(chan struct{}, worker)

	go func() {
		pmaddrs := make(map[string]int)
		for _, partition := range pull.GetPartitions() {
			ip := partition.Master.GetIp()
			port := strconv.Itoa(int(partition.Master.GetPort()))
			pmaddrs[ip+":"+port] = 0

			//pmaddrs[ip] = port
			//jobs <- job{addr, results}
		}
		for addr, _ := range pmaddrs {
			jobs <- &JobInfoStats{addr, results}
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

func (ptable *PTable) Offset(tablename string, addrs []string) (map[int32]*client.SyncOffset, error) {
	unsyncoffset := make(map[int32]*client.SyncOffset)
	/*

		Mconn, err := NewConn(addrs)
		if err != nil {
			return unsyncoffset, err
		}

		pullResp, err := Mconn.PullTable(tablename)
		if err != nil {
			return unsyncoffset, err
		}
		Mconn.RecvDone <- true
		pull := pullResp.Pull.GetInfo()[0]
	*/
	pull := ptable.pull

	// 获取 各 master 的 offset, 再获取各 slave 的 offset,然后对 offset 进行比对

	//	var masters []string

	worker := runtime.NumCPU()
	//fmt.Println(worker)
	jobsMaster := make(chan Jobber, worker)
	resultsMaster := make(chan Result, len(pull.GetPartitions()))
	donesMaster := make(chan struct{}, worker)

	//add job
	go func() {
		pmaddrs := make(map[string]int)
		for _, partition := range pull.GetPartitions() {
			mip := partition.Master.GetIp()
			mport := strconv.Itoa(int(partition.Master.GetPort()))
			//ip + port 这种 map 方式不合理
			pmaddrs[mip+":"+mport] = 0
			for _, slave := range partition.Slaves {
				sip := slave.GetIp()
				sport := strconv.Itoa(int(slave.GetPort()))

				pmaddrs[sip+":"+sport] = 0
				//jobs <- job{addr, results}
			}
			//jobs <- job{addr, results}
		}
		for addr, _ := range pmaddrs {
			jobsMaster <- &JobOffset{addr, resultsMaster}
		}
		close(jobsMaster)
	}()

	//do master job
	for i := 0; i < worker; i++ {
		go doJob(jobsMaster, donesMaster, tablename)
	}

	dataMaster := awaitForCloseResult(donesMaster, resultsMaster, worker)
	alloffsets := make(map[string][]*client.SyncOffset)

	for _, d := range dataMaster {
		alloffsets[d.addr] = d.Offsets
	}

	for _, partition := range pull.GetPartitions() {
		mIp := partition.Master.GetIp()
		mPort := strconv.Itoa(int(partition.Master.GetPort()))
		maddr := mIp + ":" + mPort
		masteroffset := alloffsets[maddr][partition.GetId()]
		for _, slave := range partition.Slaves {
			sIp := slave.GetIp()
			sPort := strconv.Itoa(int(slave.GetPort()))
			saddr := sIp + ":" + sPort
			slaveoffset := alloffsets[saddr][partition.GetId()]
			if !reflect.DeepEqual(masteroffset, slaveoffset) {
				unsyncoffset[partition.GetId()] = slaveoffset
			} else {
				unsyncoffset[partition.GetId()] = nil
			}
		}

	}
	//fmt.Println(dataSlave)
	return unsyncoffset, nil
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

type JobOffset struct {
	addr   string
	result chan<- Result
}

type Result struct {
	addr    string
	Used    int64
	Remain  int64
	QPS     int32
	Query   int64
	Offsets []*client.SyncOffset
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

func (job *JobOffset) Do(tablename string) {
	//addr := partition.Master.GetIp() + ":" + strconv.Itoa(int(partition.Master.GetPort()))

	Nconn, errN := NewConn([]string{job.addr})
	if errN != nil {
		return
	}

	inforesp, err := Nconn.InfoPartition(tablename)
	Nconn.RecvDone <- true
	if err != nil {
		return
	}
	infoPart := inforesp.GetInfoPartition()
	var result Result
	//	result.Offsets = make(map[string][]*client.SyncOffset)
	for _, i := range infoPart {
		//fmt.Println("offset:", i.GetSyncOffset())
		result.Offsets = i.GetSyncOffset()
		result.addr = job.addr
		result.Query = 0
		result.QPS = 0
	}

	job.result <- result
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
