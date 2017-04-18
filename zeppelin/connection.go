package zeppelin

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"net"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/tinytub/zep-cli/proto/ZPMeta"
	"github.com/tinytub/zep-cli/proto/client"

	"github.com/golang/protobuf/proto"
	logging "github.com/op/go-logging"
)

var logger = logging.MustGetLogger("zeppelinecore")

// 连接池 参考 https://github.com/dbehnke/influxdb-gb/blob/4ee5f88c4f5233e3b98d7885e1347582aacdff70/vendor/src/github.com/influxdb/influxdb/cluster/client_pool.go

/*
type NodeConnMap struct {
	Stream map[string]*Connection
}

type MetaConnMap struct {
	Stream []*Connection
}
*/

type Connection struct {
	//Conn *net.TCPConn
	Conn net.Conn
	//	MetaCMD   *ZPMeta.MetaCmd
	//	ClientCMD *client.CmdRequest
	Data       chan []byte
	ok         int
	HasRequest chan bool
	RecvDone   chan bool
	mu         sync.Mutex
}

/*
func (c *Connection) NewMultiNodeConn(nodelist *ZPMeta.Nodes) *NodeConnMap {
	nodeConnMap := &NodeConnMap{}
	nodeConnMap.Stream = make(map[string]*Connection)
	for _, node := range nodelist.Nodes {
		addr := *node.Node.Ip + ":" + strconv.Itoa(int(*node.Node.Port))
		tcpAddr, err := net.ResolveTCPAddr("tcp4", addr)
		logger.Info(tcpAddr)
		if err != nil {
			logger.Info("tcp resolve err:", err)
			return nodeConnMap
		}
		conn, err := net.DialTCP("tcp", nil, tcpAddr)
		if err != nil {
			logger.Info("tcp conn err:", err)
			return nodeConnMap
		}
		logger.Info(conn)

		cData := make(chan []byte, 4)
		recvDone := make(chan bool)
		hasRequest := make(chan bool)
		nodeConnMap.Stream[addr] = &Connection{Conn: conn, Data: cData, HasRequest: hasRequest, RecvDone: recvDone}
	}
	return nodeConnMap
}
*/
/*
//func (c *Connection) NewMultiMetaConn(nodelist *ZPMeta.MetaNodes) *MetaConnMap {
func NewMultiMetaConn(addrs []string) *MetaConnMap {
	//入配置文件后入库
	addrs = []string{"10.203.11.73:19221", "10.203.11.75:19221", "10.203.11.74:19221"}
	metaConnMap := &MetaConnMap{}
	for _, meta := range addrs {
		mConn, err := NewConn(meta)
		if err != nil {
			logger.Info("conn error")
			break
		}
		metaConnMap.Stream = append(metaConnMap.Stream, mConn)
	}
	/*
		for _, node := range nodelist.Followers {
			addr := *node.Ip + ":" + strconv.Itoa(int(*node.Port))
			mFolConn, _ := NewConn(addr)
			metaConnMap.Stream = append(metaConnMap.Stream, mFolConn)
		}

		mLeadAddr := *nodelist.Leader.Ip + ":" + strconv.Itoa(int(*nodelist.Leader.Port))
		mLeadConn, _ := NewConn(mLeadAddr)
		metaConnMap.Stream = append(metaConnMap.Stream, mLeadConn)
*/
/*
	return metaConnMap
}
*/

//./src/zp_table.cc:  int par_num = std::hash<std::string>()(key) % partitions_.size();

//func NewConn(ads []string) *Connection {
func NewConn(addrs []string) (*Connection, error) {

	//addrs := []string{"10.203.11.76:9221"}
	c := &Connection{}
	for _, addr := range addrs {
		conn, err := c.newConn(addr)
		if err != nil {
			logger.Info("bad conn, continue:", err)
			continue
		}
		go conn.Recv()
		return conn, nil
	}
	return c, errors.New("all bad conn")
}

func (c *Connection) newConn(addr string) (*Connection, error) {
	fConn := &Connection{}

	tcpAddr, err := net.ResolveTCPAddr("tcp4", addr)
	logger.Info(tcpAddr)
	if err != nil {
		logger.Info("tcp resolve err:", err)
		return nil, err
	}
	//conn, err := net.DialTCP("tcp", nil, tcpAddr)
	conn, err := net.DialTimeout("tcp", addr, 5000*time.Millisecond)
	if err != nil {
		logger.Info("tcp conn err:", err)
		return nil, err
	}

	fConn.Data = make(chan []byte, 4)
	fConn.RecvDone = make(chan bool)
	fConn.HasRequest = make(chan bool)
	fConn.Conn = conn

	return fConn, nil
}

// 这里不对... 得重新整
func (c *Connection) NodeConns() map[string]*Connection {
	nodes, _ := c.ListNode()
	//nodeConnMap := conn.NewMultiNodeConn(nodes.ListNode.Nodes)
	nodemap := make(map[string]*Connection)
	for _, node := range nodes.ListNode.Nodes.Nodes {
		addr := *node.Node.Ip + ":" + strconv.Itoa(int(*node.Node.Port))
		tcpAddr, err := net.ResolveTCPAddr("tcp4", addr)
		logger.Info(tcpAddr)
		if err != nil {
			logger.Info("tcp resolve err:", err)
			return nodemap
		}
		conn, err := net.DialTCP("tcp", nil, tcpAddr)
		if err != nil {
			logger.Info("tcp conn err:", err)
			return nodemap
		}
		logger.Info(conn)

		cData := make(chan []byte, 4)
		recvDone := make(chan bool)
		hasRequest := make(chan bool)
		nodemap[addr] = &Connection{Conn: conn, Data: cData, HasRequest: hasRequest, RecvDone: recvDone}
	}
	return nodemap
}

func (c *Connection) Send(data []byte) error {

	buff := make([]byte, 4+len(data))
	binary.BigEndian.PutUint32(buff[0:4], uint32(len(data)))

	copy(buff[4:], data)

	_, err := c.Conn.Write(buff)

	if err != nil {
		logger.Info("write err:", err)
		return err
	}

	logger.Info("tcp write done")
	c.HasRequest <- true
	return nil
}

func (c *Connection) Recv() {
	i := 1
	for {
		logger.Info("recv", i)
		select {
		case <-c.HasRequest:
			buf := make([]byte, 4)
			var size int32
			var data []byte
			reader := bufio.NewReader(c.Conn)
			//logger.Info(c.Conn, i)
			// 前四个字节一般表示网络包大小
			if count, err := io.ReadFull(reader, buf); err == nil && count == 4 {

				//logger.Info("reading", i)
				sbuf := bytes.NewBuffer(buf)
				// 先读数据长度
				binary.Read(sbuf, binary.BigEndian, &size)

				// 固定解析出的数据的存储空间
				data = make([]byte, size)

				// 整段读取
				count, err := io.ReadFull(reader, data)

				if err != nil {
					if err == syscall.EPIPE {
						logger.Info("io read err:", err)
						c.Conn.Close()
					}
					return
				}
				// 确认数据长度和从tcp 协议中获取的长度是否相同
				if count != int(size) {
					logger.Info("wrong count")
				}
				/*
					newdata := &ZPMeta.MetaCmdResponse{}
					// protobuf 解析
					err = proto.Unmarshal(data, newdata)
					if err != nil {
						logger.Info("unmarshaling error: ", err)
					}

					logger.Info(newdata)
				*/

				//c.Data = make(chan []byte, 1)
				c.Data <- data
				//logger.Info("send to channel done!", i)
				//logger.Info("reading done!", i)
			}
		case <-c.RecvDone:
			logger.Info("this recieve worker done", i)
			c.Conn.Close()
			return
			//					default:
			//						logger.Info("waiting request or done", i)
		case <-time.After(5000 * time.Millisecond):
			logger.Info("waiting for request or done io in 5 second")
			return
			// 这里还应该在 case 一个 停止的 sigal, 或者看要不要设置超时.

		}
		i = i + 1
	}
}

//func (c *Connection) ProtoUnserialize(data []byte, tag string) *ZPMeta.MetaCmdResponse {
func (c *Connection) ProtoUnserialize(data []byte, tag string) interface{} {
	if tag == "meta" {
		newdata := &ZPMeta.MetaCmdResponse{}
		// protobuf 解析
		err := proto.Unmarshal(data, newdata)
		if err != nil {
			logger.Info("unmarshaling error: ", err)
		}
		//	logger.Info(newdata)
		return newdata
	} else if tag == "node" {
		newdata := &client.CmdResponse{}
		// protobuf 解析
		err := proto.Unmarshal(data, newdata)
		if err != nil {
			logger.Info("unmarshaling error: ", err)
		}
		logger.Info(newdata)
		return newdata
	}
	return nil
}

/*
func InitZepConn(neednode bool) (*NodeConnMap, *MetaConnMap) { // 创建一个消息 Test
	addrs := []string{"10.203.11.73:19221", "10.203.11.75:19221", "10.203.11.74:19221"}
	var stream []Connection
	//metaMap := &MetaConnMap{Stream: stream}
	var conn *Connection
	var err error
	for _, addr := range addrs {
		conn, err = NewConn(addr)
		if err == nil {
			stream = append(stream, *conn)
			logger.Info("meta conn ok")
			break
		}
	}
	//conn.Data = make(chan []byte, 4)
	//defer conn.Conn.Close()
	//defer close(conn.RecvDone)
	//defer close(conn.Data)

	go conn.Recv()

	if neednode == true {
		nodes, err := conn.ListNode()
		nodeConnMap := conn.NewMultiNodeConn(nodes.ListNode.Nodes)
		logger.Info(nodeConnMap, &nodeConnMap)
		return nodeConnMap, metaConnMap
	}

	metas, err := conn.ListMeta()

	conn.RecvDone <- true

	if err != nil {
		logger.Info(err)
	}

	metaConnMap := conn.NewMultiMetaConn(metas.ListMeta.Nodes)
	logger.Info(metaConnMap, &metaConnMap)
	/*
		for _, node := range nodes.ListNode.Nodes.Nodes {
			nodeaddr := *node.Node.Ip + ":" + strconv.Itoa(int(*node.Node.Port))
			//logger.Info(nodeConnMap.Stream[nodeaddr])

			logger.Info(nodeaddr)
			c := nodeConnMap.Stream[nodeaddr]
			go c.Recv()
			a, _ := c.InfoStats("testtable")
			logger.Info(a)
			c.RecvDone <- true
		}
		//logger.Info(nodes)
*/
/*
	return nil, metaConnMap
}
*/

/*
func (m *MetaConnMap) GetMeta() (*Connection, error) {
	m.Strean
	for _, conn := range m.Stream {
		if conn.ok == 0 {
			return conn, nil
		}
	}

	logger.Error("all meta failed, return default meta")
	return NewConn("10.203.11.73:19221")
}
*/
