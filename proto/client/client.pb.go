// Code generated by protoc-gen-go.
// source: client.proto
// DO NOT EDIT!

/*
Package client is a generated protocol buffer package.

It is generated from these files:
	client.proto

It has these top-level messages:
	Node
	SyncOffset
	CmdRequest
	CmdResponse
	BinlogSkip
	SyncRequest
*/
package client

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Type int32

const (
	Type_SYNC          Type = 0
	Type_SET           Type = 1
	Type_GET           Type = 2
	Type_DEL           Type = 3
	Type_INFOSTATS     Type = 4
	Type_INFOCAPACITY  Type = 5
	Type_INFOPARTITION Type = 6
)

var Type_name = map[int32]string{
	0: "SYNC",
	1: "SET",
	2: "GET",
	3: "DEL",
	4: "INFOSTATS",
	5: "INFOCAPACITY",
	6: "INFOPARTITION",
}
var Type_value = map[string]int32{
	"SYNC":          0,
	"SET":           1,
	"GET":           2,
	"DEL":           3,
	"INFOSTATS":     4,
	"INFOCAPACITY":  5,
	"INFOPARTITION": 6,
}

func (x Type) Enum() *Type {
	p := new(Type)
	*p = x
	return p
}
func (x Type) String() string {
	return proto.EnumName(Type_name, int32(x))
}
func (x *Type) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(Type_value, data, "Type")
	if err != nil {
		return err
	}
	*x = Type(value)
	return nil
}
func (Type) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type SyncType int32

const (
	SyncType_CMD  SyncType = 0
	SyncType_SKIP SyncType = 1
)

var SyncType_name = map[int32]string{
	0: "CMD",
	1: "SKIP",
}
var SyncType_value = map[string]int32{
	"CMD":  0,
	"SKIP": 1,
}

func (x SyncType) Enum() *SyncType {
	p := new(SyncType)
	*p = x
	return p
}
func (x SyncType) String() string {
	return proto.EnumName(SyncType_name, int32(x))
}
func (x *SyncType) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(SyncType_value, data, "SyncType")
	if err != nil {
		return err
	}
	*x = SyncType(value)
	return nil
}
func (SyncType) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

type StatusCode int32

const (
	StatusCode_kOk       StatusCode = 0
	StatusCode_kNotFound StatusCode = 1
	StatusCode_kWait     StatusCode = 2
	StatusCode_kError    StatusCode = 3
	StatusCode_kFallback StatusCode = 4
)

var StatusCode_name = map[int32]string{
	0: "kOk",
	1: "kNotFound",
	2: "kWait",
	3: "kError",
	4: "kFallback",
}
var StatusCode_value = map[string]int32{
	"kOk":       0,
	"kNotFound": 1,
	"kWait":     2,
	"kError":    3,
	"kFallback": 4,
}

func (x StatusCode) Enum() *StatusCode {
	p := new(StatusCode)
	*p = x
	return p
}
func (x StatusCode) String() string {
	return proto.EnumName(StatusCode_name, int32(x))
}
func (x *StatusCode) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(StatusCode_value, data, "StatusCode")
	if err != nil {
		return err
	}
	*x = StatusCode(value)
	return nil
}
func (StatusCode) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

type Node struct {
	Ip               *string `protobuf:"bytes,1,req,name=ip" json:"ip,omitempty"`
	Port             *int32  `protobuf:"varint,2,req,name=port" json:"port,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *Node) Reset()                    { *m = Node{} }
func (m *Node) String() string            { return proto.CompactTextString(m) }
func (*Node) ProtoMessage()               {}
func (*Node) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Node) GetIp() string {
	if m != nil && m.Ip != nil {
		return *m.Ip
	}
	return ""
}

func (m *Node) GetPort() int32 {
	if m != nil && m.Port != nil {
		return *m.Port
	}
	return 0
}

type SyncOffset struct {
	Filenum          *int32 `protobuf:"varint,1,req,name=filenum" json:"filenum,omitempty"`
	Offset           *int64 `protobuf:"varint,2,req,name=offset" json:"offset,omitempty"`
	Partition        *int32 `protobuf:"varint,3,opt,name=partition" json:"partition,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (m *SyncOffset) Reset()                    { *m = SyncOffset{} }
func (m *SyncOffset) String() string            { return proto.CompactTextString(m) }
func (*SyncOffset) ProtoMessage()               {}
func (*SyncOffset) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *SyncOffset) GetFilenum() int32 {
	if m != nil && m.Filenum != nil {
		return *m.Filenum
	}
	return 0
}

func (m *SyncOffset) GetOffset() int64 {
	if m != nil && m.Offset != nil {
		return *m.Offset
	}
	return 0
}

func (m *SyncOffset) GetPartition() int32 {
	if m != nil && m.Partition != nil {
		return *m.Partition
	}
	return 0
}

type CmdRequest struct {
	Type             *Type            `protobuf:"varint,1,req,name=type,enum=client.Type" json:"type,omitempty"`
	Sync             *CmdRequest_Sync `protobuf:"bytes,2,opt,name=sync" json:"sync,omitempty"`
	Set              *CmdRequest_Set  `protobuf:"bytes,3,opt,name=set" json:"set,omitempty"`
	Get              *CmdRequest_Get  `protobuf:"bytes,4,opt,name=get" json:"get,omitempty"`
	Del              *CmdRequest_Del  `protobuf:"bytes,5,opt,name=del" json:"del,omitempty"`
	Info             *CmdRequest_Info `protobuf:"bytes,6,opt,name=info" json:"info,omitempty"`
	XXX_unrecognized []byte           `json:"-"`
}

func (m *CmdRequest) Reset()                    { *m = CmdRequest{} }
func (m *CmdRequest) String() string            { return proto.CompactTextString(m) }
func (*CmdRequest) ProtoMessage()               {}
func (*CmdRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *CmdRequest) GetType() Type {
	if m != nil && m.Type != nil {
		return *m.Type
	}
	return Type_SYNC
}

func (m *CmdRequest) GetSync() *CmdRequest_Sync {
	if m != nil {
		return m.Sync
	}
	return nil
}

func (m *CmdRequest) GetSet() *CmdRequest_Set {
	if m != nil {
		return m.Set
	}
	return nil
}

func (m *CmdRequest) GetGet() *CmdRequest_Get {
	if m != nil {
		return m.Get
	}
	return nil
}

func (m *CmdRequest) GetDel() *CmdRequest_Del {
	if m != nil {
		return m.Del
	}
	return nil
}

func (m *CmdRequest) GetInfo() *CmdRequest_Info {
	if m != nil {
		return m.Info
	}
	return nil
}

// Sync
type CmdRequest_Sync struct {
	Node             *Node       `protobuf:"bytes,1,req,name=node" json:"node,omitempty"`
	TableName        *string     `protobuf:"bytes,2,req,name=table_name,json=tableName" json:"table_name,omitempty"`
	SyncOffset       *SyncOffset `protobuf:"bytes,3,req,name=sync_offset,json=syncOffset" json:"sync_offset,omitempty"`
	XXX_unrecognized []byte      `json:"-"`
}

func (m *CmdRequest_Sync) Reset()                    { *m = CmdRequest_Sync{} }
func (m *CmdRequest_Sync) String() string            { return proto.CompactTextString(m) }
func (*CmdRequest_Sync) ProtoMessage()               {}
func (*CmdRequest_Sync) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2, 0} }

func (m *CmdRequest_Sync) GetNode() *Node {
	if m != nil {
		return m.Node
	}
	return nil
}

func (m *CmdRequest_Sync) GetTableName() string {
	if m != nil && m.TableName != nil {
		return *m.TableName
	}
	return ""
}

func (m *CmdRequest_Sync) GetSyncOffset() *SyncOffset {
	if m != nil {
		return m.SyncOffset
	}
	return nil
}

type CmdRequest_Set struct {
	TableName        *string `protobuf:"bytes,1,req,name=table_name,json=tableName" json:"table_name,omitempty"`
	Key              *string `protobuf:"bytes,2,req,name=key" json:"key,omitempty"`
	Value            []byte  `protobuf:"bytes,3,req,name=value" json:"value,omitempty"`
	Uuid             *string `protobuf:"bytes,4,opt,name=uuid" json:"uuid,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *CmdRequest_Set) Reset()                    { *m = CmdRequest_Set{} }
func (m *CmdRequest_Set) String() string            { return proto.CompactTextString(m) }
func (*CmdRequest_Set) ProtoMessage()               {}
func (*CmdRequest_Set) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2, 1} }

func (m *CmdRequest_Set) GetTableName() string {
	if m != nil && m.TableName != nil {
		return *m.TableName
	}
	return ""
}

func (m *CmdRequest_Set) GetKey() string {
	if m != nil && m.Key != nil {
		return *m.Key
	}
	return ""
}

func (m *CmdRequest_Set) GetValue() []byte {
	if m != nil {
		return m.Value
	}
	return nil
}

func (m *CmdRequest_Set) GetUuid() string {
	if m != nil && m.Uuid != nil {
		return *m.Uuid
	}
	return ""
}

type CmdRequest_Get struct {
	TableName        *string `protobuf:"bytes,1,req,name=table_name,json=tableName" json:"table_name,omitempty"`
	Key              *string `protobuf:"bytes,2,req,name=key" json:"key,omitempty"`
	Uuid             *string `protobuf:"bytes,3,opt,name=uuid" json:"uuid,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *CmdRequest_Get) Reset()                    { *m = CmdRequest_Get{} }
func (m *CmdRequest_Get) String() string            { return proto.CompactTextString(m) }
func (*CmdRequest_Get) ProtoMessage()               {}
func (*CmdRequest_Get) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2, 2} }

func (m *CmdRequest_Get) GetTableName() string {
	if m != nil && m.TableName != nil {
		return *m.TableName
	}
	return ""
}

func (m *CmdRequest_Get) GetKey() string {
	if m != nil && m.Key != nil {
		return *m.Key
	}
	return ""
}

func (m *CmdRequest_Get) GetUuid() string {
	if m != nil && m.Uuid != nil {
		return *m.Uuid
	}
	return ""
}

// Delete
type CmdRequest_Del struct {
	TableName        *string `protobuf:"bytes,1,req,name=table_name,json=tableName" json:"table_name,omitempty"`
	Key              *string `protobuf:"bytes,2,req,name=key" json:"key,omitempty"`
	Uuid             *string `protobuf:"bytes,3,opt,name=uuid" json:"uuid,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *CmdRequest_Del) Reset()                    { *m = CmdRequest_Del{} }
func (m *CmdRequest_Del) String() string            { return proto.CompactTextString(m) }
func (*CmdRequest_Del) ProtoMessage()               {}
func (*CmdRequest_Del) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2, 3} }

func (m *CmdRequest_Del) GetTableName() string {
	if m != nil && m.TableName != nil {
		return *m.TableName
	}
	return ""
}

func (m *CmdRequest_Del) GetKey() string {
	if m != nil && m.Key != nil {
		return *m.Key
	}
	return ""
}

func (m *CmdRequest_Del) GetUuid() string {
	if m != nil && m.Uuid != nil {
		return *m.Uuid
	}
	return ""
}

type CmdRequest_Info struct {
	TableName        *string `protobuf:"bytes,1,opt,name=table_name,json=tableName" json:"table_name,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *CmdRequest_Info) Reset()                    { *m = CmdRequest_Info{} }
func (m *CmdRequest_Info) String() string            { return proto.CompactTextString(m) }
func (*CmdRequest_Info) ProtoMessage()               {}
func (*CmdRequest_Info) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2, 4} }

func (m *CmdRequest_Info) GetTableName() string {
	if m != nil && m.TableName != nil {
		return *m.TableName
	}
	return ""
}

type CmdResponse struct {
	Type             *Type                        `protobuf:"varint,1,req,name=type,enum=client.Type" json:"type,omitempty"`
	Code             *StatusCode                  `protobuf:"varint,2,req,name=code,enum=client.StatusCode" json:"code,omitempty"`
	Msg              *string                      `protobuf:"bytes,3,opt,name=msg" json:"msg,omitempty"`
	Sync             *CmdResponse_Sync            `protobuf:"bytes,4,opt,name=sync" json:"sync,omitempty"`
	Get              *CmdResponse_Get             `protobuf:"bytes,5,opt,name=get" json:"get,omitempty"`
	Redirect         *Node                        `protobuf:"bytes,6,opt,name=redirect" json:"redirect,omitempty"`
	InfoStats        []*CmdResponse_InfoStats     `protobuf:"bytes,7,rep,name=info_stats,json=infoStats" json:"info_stats,omitempty"`
	InfoCapacity     []*CmdResponse_InfoCapacity  `protobuf:"bytes,8,rep,name=info_capacity,json=infoCapacity" json:"info_capacity,omitempty"`
	InfoPartition    []*CmdResponse_InfoPartition `protobuf:"bytes,9,rep,name=info_partition,json=infoPartition" json:"info_partition,omitempty"`
	XXX_unrecognized []byte                       `json:"-"`
}

func (m *CmdResponse) Reset()                    { *m = CmdResponse{} }
func (m *CmdResponse) String() string            { return proto.CompactTextString(m) }
func (*CmdResponse) ProtoMessage()               {}
func (*CmdResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *CmdResponse) GetType() Type {
	if m != nil && m.Type != nil {
		return *m.Type
	}
	return Type_SYNC
}

func (m *CmdResponse) GetCode() StatusCode {
	if m != nil && m.Code != nil {
		return *m.Code
	}
	return StatusCode_kOk
}

func (m *CmdResponse) GetMsg() string {
	if m != nil && m.Msg != nil {
		return *m.Msg
	}
	return ""
}

func (m *CmdResponse) GetSync() *CmdResponse_Sync {
	if m != nil {
		return m.Sync
	}
	return nil
}

func (m *CmdResponse) GetGet() *CmdResponse_Get {
	if m != nil {
		return m.Get
	}
	return nil
}

func (m *CmdResponse) GetRedirect() *Node {
	if m != nil {
		return m.Redirect
	}
	return nil
}

func (m *CmdResponse) GetInfoStats() []*CmdResponse_InfoStats {
	if m != nil {
		return m.InfoStats
	}
	return nil
}

func (m *CmdResponse) GetInfoCapacity() []*CmdResponse_InfoCapacity {
	if m != nil {
		return m.InfoCapacity
	}
	return nil
}

func (m *CmdResponse) GetInfoPartition() []*CmdResponse_InfoPartition {
	if m != nil {
		return m.InfoPartition
	}
	return nil
}

type CmdResponse_Sync struct {
	TableName        *string     `protobuf:"bytes,1,req,name=table_name,json=tableName" json:"table_name,omitempty"`
	SyncOffset       *SyncOffset `protobuf:"bytes,2,req,name=sync_offset,json=syncOffset" json:"sync_offset,omitempty"`
	XXX_unrecognized []byte      `json:"-"`
}

func (m *CmdResponse_Sync) Reset()                    { *m = CmdResponse_Sync{} }
func (m *CmdResponse_Sync) String() string            { return proto.CompactTextString(m) }
func (*CmdResponse_Sync) ProtoMessage()               {}
func (*CmdResponse_Sync) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3, 0} }

func (m *CmdResponse_Sync) GetTableName() string {
	if m != nil && m.TableName != nil {
		return *m.TableName
	}
	return ""
}

func (m *CmdResponse_Sync) GetSyncOffset() *SyncOffset {
	if m != nil {
		return m.SyncOffset
	}
	return nil
}

type CmdResponse_Get struct {
	Value            []byte `protobuf:"bytes,1,opt,name=value" json:"value,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (m *CmdResponse_Get) Reset()                    { *m = CmdResponse_Get{} }
func (m *CmdResponse_Get) String() string            { return proto.CompactTextString(m) }
func (*CmdResponse_Get) ProtoMessage()               {}
func (*CmdResponse_Get) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3, 1} }

func (m *CmdResponse_Get) GetValue() []byte {
	if m != nil {
		return m.Value
	}
	return nil
}

// InfoStats
type CmdResponse_InfoStats struct {
	TableName        *string `protobuf:"bytes,1,req,name=table_name,json=tableName" json:"table_name,omitempty"`
	TotalQuerys      *int64  `protobuf:"varint,2,req,name=total_querys,json=totalQuerys" json:"total_querys,omitempty"`
	Qps              *int32  `protobuf:"varint,3,req,name=qps" json:"qps,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *CmdResponse_InfoStats) Reset()                    { *m = CmdResponse_InfoStats{} }
func (m *CmdResponse_InfoStats) String() string            { return proto.CompactTextString(m) }
func (*CmdResponse_InfoStats) ProtoMessage()               {}
func (*CmdResponse_InfoStats) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3, 2} }

func (m *CmdResponse_InfoStats) GetTableName() string {
	if m != nil && m.TableName != nil {
		return *m.TableName
	}
	return ""
}

func (m *CmdResponse_InfoStats) GetTotalQuerys() int64 {
	if m != nil && m.TotalQuerys != nil {
		return *m.TotalQuerys
	}
	return 0
}

func (m *CmdResponse_InfoStats) GetQps() int32 {
	if m != nil && m.Qps != nil {
		return *m.Qps
	}
	return 0
}

// InfoCapacity
type CmdResponse_InfoCapacity struct {
	TableName        *string `protobuf:"bytes,1,req,name=table_name,json=tableName" json:"table_name,omitempty"`
	Used             *int64  `protobuf:"varint,2,req,name=used" json:"used,omitempty"`
	Remain           *int64  `protobuf:"varint,3,req,name=remain" json:"remain,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *CmdResponse_InfoCapacity) Reset()                    { *m = CmdResponse_InfoCapacity{} }
func (m *CmdResponse_InfoCapacity) String() string            { return proto.CompactTextString(m) }
func (*CmdResponse_InfoCapacity) ProtoMessage()               {}
func (*CmdResponse_InfoCapacity) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3, 3} }

func (m *CmdResponse_InfoCapacity) GetTableName() string {
	if m != nil && m.TableName != nil {
		return *m.TableName
	}
	return ""
}

func (m *CmdResponse_InfoCapacity) GetUsed() int64 {
	if m != nil && m.Used != nil {
		return *m.Used
	}
	return 0
}

func (m *CmdResponse_InfoCapacity) GetRemain() int64 {
	if m != nil && m.Remain != nil {
		return *m.Remain
	}
	return 0
}

// InfoPartition
type CmdResponse_InfoPartition struct {
	TableName        *string       `protobuf:"bytes,1,req,name=table_name,json=tableName" json:"table_name,omitempty"`
	SyncOffset       []*SyncOffset `protobuf:"bytes,2,rep,name=sync_offset,json=syncOffset" json:"sync_offset,omitempty"`
	XXX_unrecognized []byte        `json:"-"`
}

func (m *CmdResponse_InfoPartition) Reset()                    { *m = CmdResponse_InfoPartition{} }
func (m *CmdResponse_InfoPartition) String() string            { return proto.CompactTextString(m) }
func (*CmdResponse_InfoPartition) ProtoMessage()               {}
func (*CmdResponse_InfoPartition) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3, 4} }

func (m *CmdResponse_InfoPartition) GetTableName() string {
	if m != nil && m.TableName != nil {
		return *m.TableName
	}
	return ""
}

func (m *CmdResponse_InfoPartition) GetSyncOffset() []*SyncOffset {
	if m != nil {
		return m.SyncOffset
	}
	return nil
}

type BinlogSkip struct {
	TableName        *string `protobuf:"bytes,1,req,name=table_name,json=tableName" json:"table_name,omitempty"`
	PartitionId      *int32  `protobuf:"varint,2,req,name=partition_id,json=partitionId" json:"partition_id,omitempty"`
	Gap              *int64  `protobuf:"varint,3,req,name=gap" json:"gap,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *BinlogSkip) Reset()                    { *m = BinlogSkip{} }
func (m *BinlogSkip) String() string            { return proto.CompactTextString(m) }
func (*BinlogSkip) ProtoMessage()               {}
func (*BinlogSkip) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *BinlogSkip) GetTableName() string {
	if m != nil && m.TableName != nil {
		return *m.TableName
	}
	return ""
}

func (m *BinlogSkip) GetPartitionId() int32 {
	if m != nil && m.PartitionId != nil {
		return *m.PartitionId
	}
	return 0
}

func (m *BinlogSkip) GetGap() int64 {
	if m != nil && m.Gap != nil {
		return *m.Gap
	}
	return 0
}

type SyncRequest struct {
	SyncType         *SyncType   `protobuf:"varint,1,req,name=sync_type,json=syncType,enum=client.SyncType" json:"sync_type,omitempty"`
	Epoch            *int64      `protobuf:"varint,2,req,name=epoch" json:"epoch,omitempty"`
	From             *Node       `protobuf:"bytes,3,req,name=from" json:"from,omitempty"`
	SyncOffset       *SyncOffset `protobuf:"bytes,4,req,name=sync_offset,json=syncOffset" json:"sync_offset,omitempty"`
	Request          *CmdRequest `protobuf:"bytes,5,opt,name=request" json:"request,omitempty"`
	BinlogSkip       *BinlogSkip `protobuf:"bytes,6,opt,name=binlog_skip,json=binlogSkip" json:"binlog_skip,omitempty"`
	XXX_unrecognized []byte      `json:"-"`
}

func (m *SyncRequest) Reset()                    { *m = SyncRequest{} }
func (m *SyncRequest) String() string            { return proto.CompactTextString(m) }
func (*SyncRequest) ProtoMessage()               {}
func (*SyncRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *SyncRequest) GetSyncType() SyncType {
	if m != nil && m.SyncType != nil {
		return *m.SyncType
	}
	return SyncType_CMD
}

func (m *SyncRequest) GetEpoch() int64 {
	if m != nil && m.Epoch != nil {
		return *m.Epoch
	}
	return 0
}

func (m *SyncRequest) GetFrom() *Node {
	if m != nil {
		return m.From
	}
	return nil
}

func (m *SyncRequest) GetSyncOffset() *SyncOffset {
	if m != nil {
		return m.SyncOffset
	}
	return nil
}

func (m *SyncRequest) GetRequest() *CmdRequest {
	if m != nil {
		return m.Request
	}
	return nil
}

func (m *SyncRequest) GetBinlogSkip() *BinlogSkip {
	if m != nil {
		return m.BinlogSkip
	}
	return nil
}

func init() {
	proto.RegisterType((*Node)(nil), "client.Node")
	proto.RegisterType((*SyncOffset)(nil), "client.SyncOffset")
	proto.RegisterType((*CmdRequest)(nil), "client.CmdRequest")
	proto.RegisterType((*CmdRequest_Sync)(nil), "client.CmdRequest.Sync")
	proto.RegisterType((*CmdRequest_Set)(nil), "client.CmdRequest.Set")
	proto.RegisterType((*CmdRequest_Get)(nil), "client.CmdRequest.Get")
	proto.RegisterType((*CmdRequest_Del)(nil), "client.CmdRequest.Del")
	proto.RegisterType((*CmdRequest_Info)(nil), "client.CmdRequest.Info")
	proto.RegisterType((*CmdResponse)(nil), "client.CmdResponse")
	proto.RegisterType((*CmdResponse_Sync)(nil), "client.CmdResponse.Sync")
	proto.RegisterType((*CmdResponse_Get)(nil), "client.CmdResponse.Get")
	proto.RegisterType((*CmdResponse_InfoStats)(nil), "client.CmdResponse.InfoStats")
	proto.RegisterType((*CmdResponse_InfoCapacity)(nil), "client.CmdResponse.InfoCapacity")
	proto.RegisterType((*CmdResponse_InfoPartition)(nil), "client.CmdResponse.InfoPartition")
	proto.RegisterType((*BinlogSkip)(nil), "client.BinlogSkip")
	proto.RegisterType((*SyncRequest)(nil), "client.SyncRequest")
	proto.RegisterEnum("client.Type", Type_name, Type_value)
	proto.RegisterEnum("client.SyncType", SyncType_name, SyncType_value)
	proto.RegisterEnum("client.StatusCode", StatusCode_name, StatusCode_value)
}

func init() { proto.RegisterFile("client.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 907 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xac, 0x55, 0xdd, 0x6e, 0xeb, 0x44,
	0x10, 0x6e, 0xfc, 0x93, 0xd6, 0x93, 0xb4, 0x32, 0x2b, 0x74, 0xb0, 0x02, 0x95, 0x72, 0x22, 0x81,
	0x42, 0x39, 0xf4, 0xa2, 0xdc, 0x72, 0x53, 0xd2, 0xb4, 0xe4, 0x00, 0x69, 0xd9, 0x44, 0x42, 0x45,
	0x48, 0x39, 0xae, 0xbd, 0x09, 0x2b, 0x3b, 0xb6, 0x6b, 0x6f, 0x90, 0x22, 0xf1, 0x08, 0xbc, 0x06,
	0x8f, 0xc4, 0xfb, 0xa0, 0x19, 0x6f, 0x9c, 0x9f, 0xd3, 0xa8, 0x05, 0x9d, 0xbb, 0x99, 0xf5, 0xe7,
	0x6f, 0xc6, 0xdf, 0x7c, 0xb3, 0x86, 0x66, 0x10, 0x4b, 0x91, 0xa8, 0xf3, 0x2c, 0x4f, 0x55, 0xca,
	0xea, 0x65, 0xd6, 0x39, 0x03, 0x6b, 0x98, 0x86, 0x82, 0x9d, 0x80, 0x21, 0x33, 0xaf, 0xd6, 0x36,
	0xba, 0x0e, 0x37, 0x64, 0xc6, 0x18, 0x58, 0x59, 0x9a, 0x2b, 0xcf, 0x68, 0x1b, 0x5d, 0x9b, 0x53,
	0xdc, 0xf9, 0x0d, 0x60, 0xb4, 0x4c, 0x82, 0xdb, 0xe9, 0xb4, 0x10, 0x8a, 0x79, 0x70, 0x38, 0x95,
	0xb1, 0x48, 0x16, 0x73, 0x7a, 0xcd, 0xe6, 0xab, 0x94, 0xbd, 0x82, 0x7a, 0x4a, 0x18, 0x7a, 0xdb,
	0xe4, 0x3a, 0x63, 0x9f, 0x81, 0x93, 0xf9, 0xb9, 0x92, 0x4a, 0xa6, 0x89, 0x67, 0xb6, 0x6b, 0x5d,
	0x9b, 0xaf, 0x0f, 0x3a, 0x7f, 0xdb, 0x00, 0xbd, 0x79, 0xc8, 0xc5, 0xe3, 0x42, 0x14, 0x8a, 0xb5,
	0xc1, 0x52, 0xcb, 0x4c, 0x10, 0xf7, 0xc9, 0x45, 0xf3, 0x5c, 0x77, 0x3f, 0x5e, 0x66, 0x82, 0xd3,
	0x13, 0xf6, 0x15, 0x58, 0xc5, 0x32, 0x09, 0x3c, 0xa3, 0x5d, 0xeb, 0x36, 0x2e, 0x3e, 0x59, 0x21,
	0xd6, 0x1c, 0xe7, 0xd8, 0x2d, 0x27, 0x10, 0xeb, 0x82, 0x89, 0x0d, 0x99, 0x84, 0x7d, 0xf5, 0x14,
	0x56, 0x28, 0x8e, 0x10, 0x44, 0xce, 0x84, 0xf2, 0xac, 0xbd, 0xc8, 0x1b, 0x44, 0xce, 0x4a, 0x64,
	0x28, 0x62, 0xcf, 0xde, 0x8b, 0xbc, 0x12, 0x31, 0x47, 0x08, 0xb6, 0x2a, 0x93, 0x69, 0xea, 0xd5,
	0xf7, 0xb6, 0x3a, 0x48, 0xa6, 0x29, 0x27, 0x50, 0xeb, 0x4f, 0xb0, 0xb0, 0x71, 0x54, 0x20, 0x49,
	0xc3, 0x52, 0x81, 0xc6, 0x5a, 0x01, 0x1c, 0x17, 0xa7, 0x27, 0xec, 0x14, 0x40, 0xf9, 0x0f, 0xb1,
	0x98, 0x24, 0xfe, 0x5c, 0x90, 0xd8, 0x0e, 0x77, 0xe8, 0x64, 0xe8, 0xcf, 0x05, 0xfb, 0x06, 0x1a,
	0xf8, 0xed, 0x13, 0x3d, 0x0c, 0x93, 0x78, 0xd8, 0x8a, 0x67, 0x3d, 0x4a, 0x0e, 0x45, 0x15, 0xb7,
	0xde, 0x81, 0x39, 0x12, 0x6a, 0x87, 0xba, 0xb6, 0x4b, 0xed, 0x82, 0x19, 0x89, 0xa5, 0x2e, 0x89,
	0x21, 0xfb, 0x18, 0xec, 0x3f, 0xfc, 0x78, 0x21, 0xa8, 0x4c, 0x93, 0x97, 0x09, 0xda, 0x68, 0xb1,
	0x90, 0x21, 0xa9, 0xe9, 0x70, 0x8a, 0x5b, 0x6f, 0xc1, 0xbc, 0xf9, 0x3f, 0x15, 0x56, 0x5c, 0xe6,
	0x36, 0xd7, 0x95, 0x88, 0x3f, 0x0c, 0xd7, 0xe7, 0x60, 0xe1, 0x14, 0xde, 0x23, 0xab, 0x6d, 0x91,
	0x75, 0xfe, 0xa9, 0x43, 0x83, 0x06, 0x57, 0x64, 0x69, 0x52, 0x88, 0x17, 0x18, 0xf5, 0x0b, 0xb0,
	0x02, 0x1c, 0xa4, 0x41, 0x88, 0xf5, 0x00, 0x94, 0xaf, 0x16, 0x45, 0x8f, 0xc6, 0x89, 0xcf, 0xb1,
	0xcd, 0x79, 0x31, 0xd3, 0x3d, 0x61, 0xc8, 0xde, 0x68, 0x8b, 0x97, 0x66, 0xf4, 0xb6, 0x7c, 0x53,
	0x96, 0xdf, 0xf4, 0xf8, 0x97, 0xa5, 0x73, 0xed, 0x27, 0x4c, 0xa6, 0xc1, 0x1b, 0xd6, 0x3d, 0xca,
	0x45, 0x28, 0x73, 0x11, 0x28, 0x6d, 0xca, 0x6d, 0x7f, 0x55, 0x4f, 0xd9, 0xb7, 0x00, 0xe8, 0xca,
	0x49, 0xa1, 0x7c, 0x55, 0x78, 0x87, 0x6d, 0xb3, 0xdb, 0xb8, 0x38, 0x7d, 0x8a, 0x1b, 0xb5, 0xc3,
	0x4f, 0x2a, 0xb8, 0x23, 0x57, 0x21, 0xeb, 0xc3, 0x31, 0xbd, 0x1d, 0xf8, 0x99, 0x1f, 0x48, 0xb5,
	0xf4, 0x8e, 0x88, 0xa0, 0xbd, 0x8f, 0xa0, 0xa7, 0x71, 0xbc, 0x29, 0x37, 0x32, 0xf6, 0x3d, 0x9c,
	0x10, 0xcd, 0xfa, 0xfa, 0x70, 0x88, 0xe7, 0xf5, 0x3e, 0x9e, 0xbb, 0x15, 0x90, 0x53, 0xfd, 0x2a,
	0x6d, 0xfd, 0xaa, 0x97, 0xeb, 0x19, 0xc7, 0xec, 0xac, 0x8e, 0xf1, 0xa2, 0xd5, 0xf9, 0xb4, 0x34,
	0x76, 0xb5, 0x09, 0x68, 0x9d, 0xd5, 0x26, 0xb4, 0x26, 0xe0, 0x54, 0x0a, 0x3d, 0x57, 0xfd, 0x35,
	0x34, 0x55, 0xaa, 0xfc, 0x78, 0xf2, 0xb8, 0x10, 0xf9, 0xb2, 0xd0, 0xd7, 0x68, 0x83, 0xce, 0x7e,
	0xa6, 0x23, 0xf4, 0xca, 0x63, 0x56, 0xd0, 0xb2, 0xd9, 0x1c, 0xc3, 0xd6, 0x3d, 0x34, 0x37, 0x15,
	0x7c, 0xae, 0x06, 0x6e, 0x40, 0x21, 0x42, 0xcd, 0x4d, 0x31, 0x5e, 0xdc, 0xb9, 0x98, 0xfb, 0x32,
	0x21, 0x5e, 0x93, 0xeb, 0xac, 0x15, 0xc0, 0xf1, 0x96, 0xa8, 0xff, 0x59, 0x3d, 0xf3, 0x79, 0xf5,
	0x3a, 0xef, 0x00, 0xbe, 0x93, 0x49, 0x9c, 0xce, 0x46, 0x91, 0xcc, 0x5e, 0xa0, 0x50, 0xe5, 0x85,
	0x89, 0x0c, 0xf5, 0x6f, 0xaa, 0x51, 0x9d, 0x0d, 0x42, 0x54, 0x68, 0xe6, 0x67, 0xfa, 0x4b, 0x30,
	0xec, 0xfc, 0x65, 0x40, 0x83, 0xd6, 0x45, 0xff, 0x62, 0xbe, 0x06, 0x87, 0xda, 0xdc, 0x58, 0x5f,
	0x77, 0xb3, 0x49, 0x5a, 0xe1, 0xa3, 0x42, 0x47, 0x38, 0x57, 0x91, 0xa5, 0xc1, 0xef, 0x5a, 0xb2,
	0x32, 0xc1, 0xf5, 0x9f, 0xe6, 0xe9, 0x5c, 0xdf, 0xae, 0x3b, 0xb7, 0x34, 0x3e, 0xd9, 0x55, 0xc3,
	0x7a, 0x89, 0x97, 0xd8, 0x1b, 0x38, 0xcc, 0xcb, 0x36, 0xf5, 0x3e, 0xb3, 0xf7, 0x7f, 0x1a, 0x7c,
	0x05, 0xc1, 0x12, 0x0f, 0xa4, 0xdd, 0xa4, 0x88, 0x64, 0xa6, 0x37, 0xba, 0x7a, 0x63, 0x2d, 0x2b,
	0x87, 0x87, 0x2a, 0x3e, 0x9b, 0x80, 0x45, 0xdf, 0x75, 0x04, 0xd6, 0xe8, 0x7e, 0xd8, 0x73, 0x0f,
	0xd8, 0x21, 0x98, 0xa3, 0xfe, 0xd8, 0xad, 0x61, 0x70, 0xd3, 0x1f, 0xbb, 0x06, 0x06, 0x57, 0xfd,
	0x1f, 0x5d, 0x93, 0x1d, 0x83, 0x33, 0x18, 0x5e, 0xdf, 0x8e, 0xc6, 0x97, 0xe3, 0x91, 0x6b, 0x31,
	0x17, 0x9a, 0x98, 0xf6, 0x2e, 0xef, 0x2e, 0x7b, 0x83, 0xf1, 0xbd, 0x6b, 0xb3, 0x8f, 0xe0, 0x18,
	0x4f, 0xee, 0x2e, 0xf9, 0x78, 0x30, 0x1e, 0xdc, 0x0e, 0xdd, 0xfa, 0xd9, 0x29, 0x1c, 0xad, 0x64,
	0x44, 0xa2, 0xde, 0x4f, 0x57, 0xee, 0x01, 0x55, 0xfb, 0x61, 0x70, 0xe7, 0xd6, 0xce, 0xde, 0x02,
	0xac, 0xaf, 0x40, 0x04, 0x44, 0xb7, 0x91, 0x7b, 0x80, 0x95, 0xa2, 0x61, 0xaa, 0xae, 0xd3, 0x45,
	0x12, 0xba, 0x35, 0xe6, 0x80, 0x1d, 0xfd, 0xe2, 0x4b, 0xe5, 0x1a, 0x0c, 0xa0, 0x1e, 0xf5, 0xf3,
	0x3c, 0xcd, 0xcb, 0x7e, 0xa2, 0x6b, 0x3f, 0x8e, 0x1f, 0xfc, 0x20, 0x72, 0xad, 0x7f, 0x03, 0x00,
	0x00, 0xff, 0xff, 0x57, 0xaa, 0x50, 0xd6, 0xdd, 0x08, 0x00, 0x00,
}
