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
	"fmt"
	"os"
	"strconv"
)

var status = map[int32]string{
	0: "up",
	1: "down",
}

func ListNode(addrs []string) {
	conn := NewConn(addrs)
	//conn.mu.Lock()
	data, _ := conn.ListNode()
	if data.Code.String() != "OK" {
		fmt.Println(*data.Msg)
		os.Exit(0)
	}
	nodes := data.ListNode.Nodes.Nodes
	for _, node := range nodes {
		fmt.Printf("IP: %s, Port: %d, Status: %s\n", *node.Node.Ip, *node.Node.Port, status[*node.Status])
	}
	conn.RecvDone <- true
	return
}

func ListMeta(addrs []string) {
	conn := NewConn(addrs)
	//conn.mu.Lock()
	data, _ := conn.ListMeta()
	if data.Code.String() != "OK" {
		fmt.Println(*data.Msg)
		os.Exit(0)
	}
	metas := data.ListMeta.Nodes
	fmt.Println("Leader:")
	fmt.Printf("IP: %s, Port: %d\n", *metas.Leader.Ip, *metas.Leader.Port)
	fmt.Println("Followers:")
	for _, follower := range metas.Followers {
		fmt.Printf("IP: %s, Port: %d\n", *follower.Ip, *follower.Port)
	}
	conn.RecvDone <- true
	return
}

func ListTable(addrs []string) {
	conn := NewConn(addrs)
	//conn.mu.Lock()
	data, _ := conn.ListTable()
	if data.Code.String() != "OK" {
		fmt.Println(*data.Msg)
		os.Exit(0)
	}
	tables := data.ListTable.Tables.Name
	for _, table := range tables {
		fmt.Println(table)
	}
	conn.RecvDone <- true
	return
}

func CreateTable(name string, num int32, addrs []string) {
	conn := NewConn(addrs)
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

func Set(tablename string, key string, value string, addrs []string) {
	Mconn := NewConn(addrs)
	fmt.Println(Mconn)
	tableinfo, _ := Mconn.PullTable(tablename)

	partcount := len(tableinfo.Pull.Info[0].Partitions)
	//./src/zp_table.cc:  int par_num = std::hash<std::string>()(key) % partitions_.size();

	/* 动态链接库的编译方法
	gcc -c chash.cc -std=c++11
	ar rv libchash.a chash.o
	mv libchash.a ../lib
	测试
	g++ -o chash chash.cc -std=c++11
	*/
	fmt.Println(C.chash(C.CString(key)))
	parNum := uint(C.chash(C.CString(key))) % uint(partcount)
	fmt.Println(parNum)
	nodemaster := tableinfo.Pull.Info[0].Partitions[parNum-1].Master
	//nodemaster.GetIp() + ":" + strconv.Itoa(int(nodemaster.GetPort()))

	var nodeaddrs []string
	nodeaddrs = append(nodeaddrs, nodemaster.GetIp()+":"+strconv.Itoa(int(nodemaster.GetPort())))
	fmt.Println(nodeaddrs)
	Nconn := NewConn(nodeaddrs)
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
	getresp, err := Nconn.Get(tablename, key)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(getresp)
	/*
		//conn.mu.Lock()
		val, _ := getBytes(value)
		//fmt.Println([]byte(value))
		fmt.Println(val)
		pNum := getPartition(conn, tablename)
		fmt.Println(pNum)
		//data, _ := conn.Set(tablename, key, val)
		//data, _ := conn.Set(tablename, key, []byte(value))
		//fmt.Println(data)
		conn.RecvDone <- true
		return
	*/
}

func locationNode() {}

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
