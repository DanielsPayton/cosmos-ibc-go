// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: ibc/applications/middleware/fee/v1/genesis.proto

package types

import (
	fmt "fmt"
	_ "github.com/cosmos/ibc-go/modules/core/04-channel/types"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// GenesisState defines the fee middleware genesis state
type GenesisState struct {
	// A mapping of packets -> escrowed fees
	PacketsFees []*IdentifiedPacketFee `protobuf:"bytes,1,rep,name=packets_fees,json=packetsFees,proto3" json:"packets_fees,omitempty" yaml:"packets_fees"`
}

func (m *GenesisState) Reset()         { *m = GenesisState{} }
func (m *GenesisState) String() string { return proto.CompactTextString(m) }
func (*GenesisState) ProtoMessage()    {}
func (*GenesisState) Descriptor() ([]byte, []int) {
	return fileDescriptor_6a4f6c8d78b5dfbf, []int{0}
}
func (m *GenesisState) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *GenesisState) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_GenesisState.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *GenesisState) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GenesisState.Merge(m, src)
}
func (m *GenesisState) XXX_Size() int {
	return m.Size()
}
func (m *GenesisState) XXX_DiscardUnknown() {
	xxx_messageInfo_GenesisState.DiscardUnknown(m)
}

var xxx_messageInfo_GenesisState proto.InternalMessageInfo

func (m *GenesisState) GetPacketsFees() []*IdentifiedPacketFee {
	if m != nil {
		return m.PacketsFees
	}
	return nil
}

func init() {
	proto.RegisterType((*GenesisState)(nil), "ibc.applications.middleware.fee.v1.GenesisState")
}

func init() {
	proto.RegisterFile("ibc/applications/middleware/fee/v1/genesis.proto", fileDescriptor_6a4f6c8d78b5dfbf)
}

var fileDescriptor_6a4f6c8d78b5dfbf = []byte{
	// 290 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x90, 0xb1, 0x4e, 0xf3, 0x30,
	0x1c, 0xc4, 0x1b, 0x7d, 0xd2, 0x37, 0xa4, 0x9d, 0x0a, 0x12, 0xa8, 0x83, 0x81, 0x4c, 0x0c, 0xd4,
	0xa6, 0x61, 0x40, 0x30, 0x76, 0x28, 0x42, 0x62, 0x40, 0xb0, 0xb1, 0x20, 0xc7, 0xb9, 0xa4, 0x16,
	0x49, 0x1c, 0xe5, 0xef, 0x16, 0x75, 0xe6, 0x05, 0x78, 0x2c, 0xc6, 0x8e, 0x4c, 0x08, 0x25, 0x6f,
	0xc0, 0x13, 0x20, 0x27, 0x95, 0xe8, 0x06, 0xdb, 0x59, 0xfe, 0xdd, 0xf9, 0xce, 0xfe, 0xa9, 0x8e,
	0x94, 0x90, 0x65, 0x99, 0x69, 0x25, 0xad, 0x36, 0x05, 0x89, 0x5c, 0xc7, 0x71, 0x86, 0x67, 0x59,
	0x41, 0x24, 0x80, 0x58, 0x4e, 0x44, 0x8a, 0x02, 0xa4, 0x89, 0x97, 0x95, 0xb1, 0x66, 0x18, 0xe8,
	0x48, 0xf1, 0x6d, 0x07, 0xff, 0x71, 0xf0, 0x04, 0xe0, 0xcb, 0xc9, 0x68, 0x37, 0x35, 0xa9, 0x69,
	0x71, 0xe1, 0x54, 0xe7, 0x1c, 0x9d, 0xfc, 0xe1, 0x2d, 0x17, 0xd0, 0xd1, 0x47, 0x8e, 0x56, 0xa6,
	0x82, 0x50, 0x73, 0x59, 0x14, 0xc8, 0xdc, 0xf5, 0x46, 0x76, 0x48, 0xf0, 0xe2, 0xf9, 0x83, 0xab,
	0xae, 0xdc, 0xbd, 0x95, 0x16, 0x43, 0xf2, 0x07, 0xa5, 0x54, 0x4f, 0xb0, 0xf4, 0x98, 0x00, 0xb4,
	0xef, 0x1d, 0xfe, 0x3b, 0xee, 0x87, 0xe7, 0xfc, 0xf7, 0xca, 0xfc, 0x3a, 0x46, 0x61, 0x75, 0xa2,
	0x11, 0xdf, 0xb6, 0x09, 0x33, 0x60, 0xba, 0xf7, 0xf5, 0x71, 0xb0, 0xb3, 0x92, 0x79, 0x76, 0x19,
	0x6c, 0xc7, 0x06, 0x77, 0xfd, 0xcd, 0x71, 0x06, 0xd0, 0xf4, 0xe6, 0xad, 0x66, 0xde, 0xba, 0x66,
	0xde, 0x67, 0xcd, 0xbc, 0xd7, 0x86, 0xf5, 0xd6, 0x0d, 0xeb, 0xbd, 0x37, 0xac, 0xf7, 0x10, 0xa6,
	0xda, 0xce, 0x17, 0x11, 0x57, 0x26, 0x17, 0xca, 0x50, 0x6e, 0x48, 0xe8, 0x48, 0x8d, 0x53, 0x23,
	0x72, 0x13, 0x2f, 0x32, 0x90, 0xfb, 0x0d, 0x12, 0xe1, 0xc5, 0xd8, 0xad, 0xb7, 0xab, 0x12, 0x14,
	0xfd, 0x6f, 0xa7, 0x9d, 0x7d, 0x07, 0x00, 0x00, 0xff, 0xff, 0x27, 0xbb, 0x9c, 0x19, 0x99, 0x01,
	0x00, 0x00,
}

func (m *GenesisState) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GenesisState) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *GenesisState) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.PacketsFees) > 0 {
		for iNdEx := len(m.PacketsFees) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.PacketsFees[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func encodeVarintGenesis(dAtA []byte, offset int, v uint64) int {
	offset -= sovGenesis(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *GenesisState) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.PacketsFees) > 0 {
		for _, e := range m.PacketsFees {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	return n
}

func sovGenesis(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozGenesis(x uint64) (n int) {
	return sovGenesis(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *GenesisState) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGenesis
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: GenesisState: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GenesisState: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PacketsFees", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.PacketsFees = append(m.PacketsFees, &IdentifiedPacketFee{})
			if err := m.PacketsFees[len(m.PacketsFees)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGenesis(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGenesis
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipGenesis(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowGenesis
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthGenesis
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupGenesis
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthGenesis
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthGenesis        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowGenesis          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupGenesis = fmt.Errorf("proto: unexpected end of group")
)
