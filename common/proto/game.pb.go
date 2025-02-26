// Code generated by protoc-gen-go. DO NOT EDIT.
// source: game.proto

package pb

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// 排行榜结构体
type PbRank struct {
	Member               string   `protobuf:"bytes,1,opt,name=Member,proto3" json:"Member,omitempty"`
	Score                int64    `protobuf:"varint,2,opt,name=Score,proto3" json:"Score,omitempty"`
	Rank                 int32    `protobuf:"varint,3,opt,name=Rank,proto3" json:"Rank,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PbRank) Reset()         { *m = PbRank{} }
func (m *PbRank) String() string { return proto.CompactTextString(m) }
func (*PbRank) ProtoMessage()    {}
func (*PbRank) Descriptor() ([]byte, []int) {
	return fileDescriptor_38fc58335341d769, []int{0}
}

func (m *PbRank) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PbRank.Unmarshal(m, b)
}
func (m *PbRank) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PbRank.Marshal(b, m, deterministic)
}
func (m *PbRank) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PbRank.Merge(m, src)
}
func (m *PbRank) XXX_Size() int {
	return xxx_messageInfo_PbRank.Size(m)
}
func (m *PbRank) XXX_DiscardUnknown() {
	xxx_messageInfo_PbRank.DiscardUnknown(m)
}

var xxx_messageInfo_PbRank proto.InternalMessageInfo

func (m *PbRank) GetMember() string {
	if m != nil {
		return m.Member
	}
	return ""
}

func (m *PbRank) GetScore() int64 {
	if m != nil {
		return m.Score
	}
	return 0
}

func (m *PbRank) GetRank() int32 {
	if m != nil {
		return m.Rank
	}
	return 0
}

func init() {
	proto.RegisterType((*PbRank)(nil), "pb.PbRank")
}

func init() { proto.RegisterFile("game.proto", fileDescriptor_38fc58335341d769) }

var fileDescriptor_38fc58335341d769 = []byte{
	// 106 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x4a, 0x4f, 0xcc, 0x4d,
	0xd5, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x2a, 0x48, 0x52, 0xf2, 0xe2, 0x62, 0x0b, 0x48,
	0x0a, 0x4a, 0xcc, 0xcb, 0x16, 0x12, 0xe3, 0x62, 0xf3, 0x4d, 0xcd, 0x4d, 0x4a, 0x2d, 0x92, 0x60,
	0x54, 0x60, 0xd4, 0xe0, 0x0c, 0x82, 0xf2, 0x84, 0x44, 0xb8, 0x58, 0x83, 0x93, 0xf3, 0x8b, 0x52,
	0x25, 0x98, 0x14, 0x18, 0x35, 0x98, 0x83, 0x20, 0x1c, 0x21, 0x21, 0x2e, 0x16, 0x90, 0x2e, 0x09,
	0x66, 0x05, 0x46, 0x0d, 0xd6, 0x20, 0x30, 0x3b, 0x89, 0x0d, 0x6c, 0xac, 0x31, 0x20, 0x00, 0x00,
	0xff, 0xff, 0x5f, 0x77, 0x65, 0x2d, 0x64, 0x00, 0x00, 0x00,
}
