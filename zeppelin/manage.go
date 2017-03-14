package zeppelin

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
)

var status = map[int32]string{
	0: "up",
	1: "down",
}

func ListNode() {
	conn := NewConn()
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

func ListMeta() {
	conn := NewConn()
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

func ListTable() {
	conn := NewConn()
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

func CreateTable(name string, num int32) {
	conn := NewConn()
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

func Set(tablename string, key string, value string) {
	conn := NewConn()
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
