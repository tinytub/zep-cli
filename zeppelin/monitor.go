package zeppelin

import "fmt"

//	"github.com/tinytub/zeppeline-monitor/zeppelin"

func NodeStats() {
	conn := NewConn()
	//conn.mu.Lock()
	data, _ := conn.ListNode()
	//fmt.Println(data)
	logger.Info(data)
	nodes := data.ListNode.Nodes.Nodes
	fmt.Println("node numric:", len(nodes))
	for _, node := range nodes {
		fmt.Println(node.Node)
	}
	//conn.Send(data)
	//time.Sleep(time.Second * 10)
	//conn.mu.Unlock()
	conn.RecvDone <- true
	return
}

func MetaStats() {
	conn := NewConn()
	//conn.mu.Lock()
	data, _ := conn.ListMeta()
	metas := data.ListMeta.Nodes
	fmt.Println("leader:", metas.Leader)
	fmt.Println("followers:", metas.Followers)
	/*
		for _, node := range data.ListMeta.Nodes {
			continue
		}
	*/
	//fmt.Println(data)
	conn.RecvDone <- true
	return
}

func TableStats() {
	conn := NewConn()
	//conn.mu.Lock()
	data, _ := conn.ListTable()
	tables := data.ListTable.Tables.Name
	fmt.Println(tables)
	for _, table := range tables {
		tableinfo, _ := conn.PullTable(table)
		for _, n := range tableinfo.Pull.Info {
			fmt.Printf("tablename: %s, partition num: %d\n", *n.Name, len(n.Partitions))

		}
		//fmt.Println(tableinfo)
	}
	conn.RecvDone <- true
	return
}
