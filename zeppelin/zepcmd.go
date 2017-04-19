package zeppelin

import (
	"errors"
	"time"

	"github.com/tinytub/zep-cli/proto/ZPMeta"
	"github.com/tinytub/zep-cli/proto/client"

	"github.com/golang/protobuf/proto"
)

//TODO: node 节点由外部选择
func (c *Connection) PullTable(tablename string) (*ZPMeta.MetaCmdResponse, error) {
	c.mu.Lock()
	cmd, err := c.MakeCmdPull(tablename, "", 0)
	if err != nil {
		logger.Info("marshal proto error", err)
	}
	//	c.Send(cmd)
	c.Send(cmd)
	data, err := c.getData("meta")
	if err != nil {
		return data.(*ZPMeta.MetaCmdResponse), err
	}
	c.mu.Unlock()
	return data.(*ZPMeta.MetaCmdResponse), nil
}

func (c *Connection) MakeCmdPull(name string, node string, port int32) ([]byte, error) {
	var raw_cmd *ZPMeta.MetaCmd
	switch {
	case name != "":
		raw_cmd = &ZPMeta.MetaCmd{
			Type: ZPMeta.Type_PULL.Enum(),
			Pull: &ZPMeta.MetaCmd_Pull{Name: &name},
		}
	case node != "" && port != 0:

		raw_cmd = &ZPMeta.MetaCmd{
			Type: ZPMeta.Type_PULL.Enum(),
			Pull: &ZPMeta.MetaCmd_Pull{Node: &ZPMeta.Node{Ip: &node,
				Port: &port},
			},
		}
	}

	/*
		raw_cmd := &ZPMeta.MetaCmd{
			Type: ZPMeta.Type_PULL.Enum(),
			Pull: &ZPMeta.MetaCmd_Pull{Name: &arg},
		}
	*/

	return proto.Marshal(raw_cmd)
}

func (c *Connection) PullNode(node string, port int32) (*ZPMeta.MetaCmdResponse, error) {
	c.mu.Lock()
	//cmd, err := c.MakeCmdPullnode(node, port)
	cmd, err := c.MakeCmdPull("", node, port)
	if err != nil {
		logger.Info("marshal proto error", err)
	}
	//	c.Send(cmd)
	c.Send(cmd)
	data, err := c.getData("meta")
	if err != nil {
		return data.(*ZPMeta.MetaCmdResponse), err
	}

	c.mu.Unlock()
	return data.(*ZPMeta.MetaCmdResponse), nil
}

func (c *Connection) MakeCmdPullnode(node string, port int32) ([]byte, error) {

	raw_cmd := &ZPMeta.MetaCmd{
		Type: ZPMeta.Type_PULL.Enum(),
		Pull: &ZPMeta.MetaCmd_Pull{Node: &ZPMeta.Node{Ip: &node,
			Port: &port},
		},
	}
	/*
			type Node struct {
			Ip               *string `protobuf:"bytes,1,req,name=ip" json:"ip,omitempty"`
			Port             *int32  `protobuf:"varint,2,req,name=port" json:"port,omitempty"`
			XXX_unrecognized []byte  `json:"-"`
		}
	*/

	return proto.Marshal(raw_cmd)
}

func (c *Connection) ListTable() (*ZPMeta.MetaCmdResponse, error) {

	c.mu.Lock()
	//conn := NewConn("10.203.11.74:19221")
	cmd, err := c.MakeCmdListTable()
	if err != nil {
		logger.Info("marshal proto error", err)
	}
	//	cmd2, _ := c.MakeCmdPull("testtable")
	//	c.Data = make(chan []byte, 2)

	c.Send(cmd)
	data, err := c.getData("meta")
	if err != nil {
		return data.(*ZPMeta.MetaCmdResponse), err
	}

	/*
		c.Send(cmd2)
		data2 := c.getData()
		logger.Info(data2)
	*/
	c.mu.Unlock()
	return data.(*ZPMeta.MetaCmdResponse), nil

}

func (c *Connection) MakeCmdListTable() ([]byte, error) {

	raw_cmd := &ZPMeta.MetaCmd{
		Type: ZPMeta.Type_LISTTABLE.Enum(),
	}

	return proto.Marshal(raw_cmd)
}

func (c *Connection) ListNode() (*ZPMeta.MetaCmdResponse, error) {
	cmd, err := c.MakeCmdListNode()
	if err != nil {
		logger.Info("marshal proto error", err)
	}
	c.Send(cmd)
	c.mu.Lock()

	data, err := c.getData("meta")
	if err != nil {
		return data.(*ZPMeta.MetaCmdResponse), err
	}

	c.mu.Unlock()

	return data.(*ZPMeta.MetaCmdResponse), nil
}

func (c *Connection) MakeCmdListNode() ([]byte, error) {

	raw_cmd := &ZPMeta.MetaCmd{
		Type: ZPMeta.Type_LISTNODE.Enum(),
	}

	return proto.Marshal(raw_cmd)
}

func (c *Connection) ListMeta() (*ZPMeta.MetaCmdResponse, error) {
	cmd, err := c.MakeCmdListMeta()
	if err != nil {
		logger.Info("marshal proto error", err)
	}
	logger.Info("listmeta")
	c.mu.Lock()

	c.Send(cmd)
	data, err := c.getData("meta")
	if err != nil {
		return data.(*ZPMeta.MetaCmdResponse), err
	}

	c.mu.Unlock()

	return data.(*ZPMeta.MetaCmdResponse), nil
}

func (c *Connection) MakeCmdListMeta() ([]byte, error) {

	raw_cmd := &ZPMeta.MetaCmd{
		Type: ZPMeta.Type_LISTMETA.Enum(),
	}

	return proto.Marshal(raw_cmd)
}

func (c *Connection) CreateTable(name string, num int32) (*ZPMeta.MetaCmdResponse, error) {
	cmd, err := c.MakeCmdCreateTable(name, num)
	if err != nil {
		logger.Info("marshal proto error", err)
	}
	c.mu.Lock()

	c.Send(cmd)
	data, err := c.getData("meta")
	if err != nil {
		return data.(*ZPMeta.MetaCmdResponse), err
	}

	c.mu.Unlock()

	return data.(*ZPMeta.MetaCmdResponse), nil
}

func (c *Connection) MakeCmdCreateTable(name string, num int32) ([]byte, error) {

	raw_cmd := &ZPMeta.MetaCmd{
		Type: ZPMeta.Type_INIT.Enum(),
		Init: &ZPMeta.MetaCmd_Init{
			Name: &name,
			Num:  &num,
		},
	}
	return proto.Marshal(raw_cmd)
}

// client
func (c *Connection) InfoStats(tablename string) (*client.CmdResponse, error) {
	c.mu.Lock()
	cmd, err := c.MakeCmdInfoStats(tablename)
	if err != nil {
		logger.Info("marshal proto error", err)
	}
	c.Send(cmd)
	data, err := c.getData("node")
	if err != nil {
		return data.(*client.CmdResponse), err
	}

	c.mu.Unlock()
	return data.(*client.CmdResponse), nil
}

/*
func (c *Connection) TotalQPS(nodeconns map[string]*Connection, table string) int {
	var totalQPS int
	for _, nConn := range nodeconns {
		go nConn.Recv()
		data, err := nConn.InfoStats(table)
		if err == nil {
			infoStats := data.GetInfoStats()
			totalQPS = totalQPS + int(*infoStats[0].Qps)
		}
		nConn.RecvDone <- true
	}
	return totalQPS

}
*/

func (c *Connection) MakeCmdInfoStats(tablename string) ([]byte, error) {
	logger.Info("tablename is:", tablename)
	raw_cmd := &client.CmdRequest{
		Type: client.Type_INFOSTATS.Enum(),
		Info: &client.CmdRequest_Info{TableName: &tablename},
	}

	return proto.Marshal(raw_cmd)
}

func (c *Connection) InfoCapacity(tablename string) (*client.CmdResponse, error) {
	c.mu.Lock()
	cmd, err := c.MakeCmdInfoCapacity(tablename)
	if err != nil {
		logger.Info("marshal proto error", err)
	}
	c.Send(cmd)
	data, err := c.getData("node")
	if err != nil {
		return data.(*client.CmdResponse), err
	}

	c.mu.Unlock()
	return data.(*client.CmdResponse), nil
}

func (c *Connection) MakeCmdInfoCapacity(tablename string) ([]byte, error) {
	//logger.Info("tablename is:", tablename)
	raw_cmd := &client.CmdRequest{
		Type: client.Type_INFOCAPACITY.Enum(),
		Info: &client.CmdRequest_Info{TableName: &tablename},
	}

	return proto.Marshal(raw_cmd)
}

func (c *Connection) InfoPartition(tablename string) (*client.CmdResponse, error) {
	c.mu.Lock()
	cmd, err := c.MakeCmdInfoPartition(tablename)
	if err != nil {
		logger.Info("marshal proto error", err)
	}
	c.Send(cmd)
	data, err := c.getData("node")
	if err != nil {
		return data.(*client.CmdResponse), err
	}

	c.mu.Unlock()
	return data.(*client.CmdResponse), nil
}

func (c *Connection) MakeCmdInfoPartition(tablename string) ([]byte, error) {
	//logger.Info("tablename is:", tablename)
	raw_cmd := &client.CmdRequest{
		Type: client.Type_INFOPARTITION.Enum(),
		Info: &client.CmdRequest_Info{TableName: &tablename},
	}

	return proto.Marshal(raw_cmd)
}

/*
func (c *Connection) Ping() bool {
	c.Send(&PingPacket{})

	select {
	case _, ok := <-c.pongs:
		return ok
	case <-time.After(500 * time.Millisecond):
		return false
	}
}
*/

func (c *Connection) Set(tablename string, key string, value []byte) (*client.CmdResponse, error) {
	cmd, err := c.MakeCmdSet(tablename, key, value)
	if err != nil {
		logger.Info("marshal proto error", err)
	}
	c.mu.Lock()

	c.Send(cmd)
	data, err := c.getData("node")
	if err != nil {
		return data.(*client.CmdResponse), err
	}

	c.mu.Unlock()

	return data.(*client.CmdResponse), nil
}

func (c *Connection) MakeCmdSet(tablename string, key string, value []byte) ([]byte, error) {
	raw_cmd := &client.CmdRequest{
		Type: client.Type_SET.Enum(),
		Set: &client.CmdRequest_Set{
			TableName: &tablename,
			Key:       &key,
			Value:     value,
		},
	}
	return proto.Marshal(raw_cmd)
}

func (c *Connection) Get(tablename string, key string) (*client.CmdResponse, error) {
	cmd, err := c.MakeCmdGet(tablename, key)
	if err != nil {
		logger.Info("marshal proto error", err)
	}
	c.mu.Lock()

	c.Send(cmd)
	data, err := c.getData("node")
	if err != nil {
		return data.(*client.CmdResponse), err
	}

	c.mu.Unlock()

	return data.(*client.CmdResponse), nil
}

func (c *Connection) MakeCmdGet(tablename string, key string) ([]byte, error) {
	raw_cmd := &client.CmdRequest{
		Type: client.Type_GET.Enum(),
		Get: &client.CmdRequest_Get{
			TableName: &tablename,
			Key:       &key,
			//Value:     value,
		},
	}
	return proto.Marshal(raw_cmd)
}

func (c *Connection) Ping() (*ZPMeta.MetaCmdResponse, error) {
	cmd, err := c.MakeCmdPing()
	if err != nil {
		logger.Info("marshal proto error", err)
	}
	c.Send(cmd)
	c.mu.Lock()
	data, err := c.getData("meta")
	if err != nil {
		return data.(*ZPMeta.MetaCmdResponse), err
	}

	c.mu.Unlock()
	return data.(*ZPMeta.MetaCmdResponse), nil
}

func (c *Connection) MakeCmdPing() ([]byte, error) {
	raw_cmd := &ZPMeta.MetaCmd{
		Type: ZPMeta.Type_PING.Enum(),
	}
	return proto.Marshal(raw_cmd)
}

//func (c *Connection) getData(tag string) *ZPMeta.MetaCmdResponse {
func (c *Connection) getData(tag string) (interface{}, error) {
	//TODO 这里可以加 retry
	timeout := time.After(1 * time.Second)
	//	tick := time.Tick(500 * time.Millisecond)

	for {
		select {
		/*
			case <-tick:
				rawdata := <-c.Data
				fmt.Println(time.Now())
		*/
		case rawdata := <-c.Data:
			if rawdata != nil {
				newdata := c.ProtoUnserialize(rawdata, tag)
				//close(c.RecvDone)
				return newdata, nil
			}
		case <-timeout:
			logger.Info("time out 1 second")
			nildata := c.ProtoUnserialize(nil, tag)
			return nildata, errors.New("time out in 1 second")
			/*
				case <-time.After(5000 * time.Millisecond):
					logger.Info("time out 5000ms")
					//		return &ZPMeta.MetaCmdResponse{}
					//return nil
					continue
			*/
		}
	}
}

func (c *Connection) NilStruct(tag string) {

}
