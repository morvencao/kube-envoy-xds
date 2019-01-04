// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: envoy/config/transport_socket/capture/v2alpha/capture.proto

package v2

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"
import core "github.com/morvencao/kube-envoy-xds/envoy/api/v2/core"

import io "io"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

// File format.
type FileSink_Format int32

const (
	// Binary proto format as per :ref:`Trace
	// <envoy_api_msg_data.tap.v2alpha.Trace>`.
	FileSink_PROTO_BINARY FileSink_Format = 0
	// Text proto format as per :ref:`Trace
	// <envoy_api_msg_data.tap.v2alpha.Trace>`.
	FileSink_PROTO_TEXT FileSink_Format = 1
)

var FileSink_Format_name = map[int32]string{
	0: "PROTO_BINARY",
	1: "PROTO_TEXT",
}
var FileSink_Format_value = map[string]int32{
	"PROTO_BINARY": 0,
	"PROTO_TEXT":   1,
}

func (x FileSink_Format) String() string {
	return proto.EnumName(FileSink_Format_name, int32(x))
}
func (FileSink_Format) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_capture_910b340f585e95f6, []int{0, 0}
}

// File sink.
//
// .. warning::
//
//   The current file sink implementation buffers the entire trace in memory
//   prior to writing. This will OOM for long lived sockets and/or where there
//   is a large amount of traffic on the socket.
type FileSink struct {
	// Path prefix. The output file will be of the form <path_prefix>_<id>.pb, where <id> is an
	// identifier distinguishing the recorded trace for individual socket instances (the Envoy
	// connection ID).
	PathPrefix           string          `protobuf:"bytes,1,opt,name=path_prefix,json=pathPrefix,proto3" json:"path_prefix,omitempty"`
	Format               FileSink_Format `protobuf:"varint,2,opt,name=format,proto3,enum=envoy.config.transport_socket.capture.v2alpha.FileSink_Format" json:"format,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *FileSink) Reset()         { *m = FileSink{} }
func (m *FileSink) String() string { return proto.CompactTextString(m) }
func (*FileSink) ProtoMessage()    {}
func (*FileSink) Descriptor() ([]byte, []int) {
	return fileDescriptor_capture_910b340f585e95f6, []int{0}
}
func (m *FileSink) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *FileSink) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_FileSink.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (dst *FileSink) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FileSink.Merge(dst, src)
}
func (m *FileSink) XXX_Size() int {
	return m.Size()
}
func (m *FileSink) XXX_DiscardUnknown() {
	xxx_messageInfo_FileSink.DiscardUnknown(m)
}

var xxx_messageInfo_FileSink proto.InternalMessageInfo

func (m *FileSink) GetPathPrefix() string {
	if m != nil {
		return m.PathPrefix
	}
	return ""
}

func (m *FileSink) GetFormat() FileSink_Format {
	if m != nil {
		return m.Format
	}
	return FileSink_PROTO_BINARY
}

// Configuration for capture transport socket. This wraps another transport socket, providing the
// ability to interpose and record in plain text any traffic that is surfaced to Envoy.
type Capture struct {
	// Types that are valid to be assigned to SinkSelector:
	//	*Capture_FileSink
	SinkSelector isCapture_SinkSelector `protobuf_oneof:"sink_selector"`
	// The underlying transport socket being wrapped.
	TransportSocket      *core.TransportSocket `protobuf:"bytes,2,opt,name=transport_socket,json=transportSocket,proto3" json:"transport_socket,omitempty"`
	XXX_NoUnkeyedLiteral struct{}              `json:"-"`
	XXX_unrecognized     []byte                `json:"-"`
	XXX_sizecache        int32                 `json:"-"`
}

func (m *Capture) Reset()         { *m = Capture{} }
func (m *Capture) String() string { return proto.CompactTextString(m) }
func (*Capture) ProtoMessage()    {}
func (*Capture) Descriptor() ([]byte, []int) {
	return fileDescriptor_capture_910b340f585e95f6, []int{1}
}
func (m *Capture) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Capture) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Capture.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (dst *Capture) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Capture.Merge(dst, src)
}
func (m *Capture) XXX_Size() int {
	return m.Size()
}
func (m *Capture) XXX_DiscardUnknown() {
	xxx_messageInfo_Capture.DiscardUnknown(m)
}

var xxx_messageInfo_Capture proto.InternalMessageInfo

type isCapture_SinkSelector interface {
	isCapture_SinkSelector()
	MarshalTo([]byte) (int, error)
	Size() int
}

type Capture_FileSink struct {
	FileSink *FileSink `protobuf:"bytes,1,opt,name=file_sink,json=fileSink,proto3,oneof"`
}

func (*Capture_FileSink) isCapture_SinkSelector() {}

func (m *Capture) GetSinkSelector() isCapture_SinkSelector {
	if m != nil {
		return m.SinkSelector
	}
	return nil
}

func (m *Capture) GetFileSink() *FileSink {
	if x, ok := m.GetSinkSelector().(*Capture_FileSink); ok {
		return x.FileSink
	}
	return nil
}

func (m *Capture) GetTransportSocket() *core.TransportSocket {
	if m != nil {
		return m.TransportSocket
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*Capture) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _Capture_OneofMarshaler, _Capture_OneofUnmarshaler, _Capture_OneofSizer, []interface{}{
		(*Capture_FileSink)(nil),
	}
}

func _Capture_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*Capture)
	// sink_selector
	switch x := m.SinkSelector.(type) {
	case *Capture_FileSink:
		_ = b.EncodeVarint(1<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.FileSink); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("Capture.SinkSelector has unexpected type %T", x)
	}
	return nil
}

func _Capture_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*Capture)
	switch tag {
	case 1: // sink_selector.file_sink
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(FileSink)
		err := b.DecodeMessage(msg)
		m.SinkSelector = &Capture_FileSink{msg}
		return true, err
	default:
		return false, nil
	}
}

func _Capture_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*Capture)
	// sink_selector
	switch x := m.SinkSelector.(type) {
	case *Capture_FileSink:
		s := proto.Size(x.FileSink)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

func init() {
	proto.RegisterType((*FileSink)(nil), "envoy.config.transport_socket.capture.v2alpha.FileSink")
	proto.RegisterType((*Capture)(nil), "envoy.config.transport_socket.capture.v2alpha.Capture")
	proto.RegisterEnum("envoy.config.transport_socket.capture.v2alpha.FileSink_Format", FileSink_Format_name, FileSink_Format_value)
}
func (m *FileSink) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *FileSink) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.PathPrefix) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintCapture(dAtA, i, uint64(len(m.PathPrefix)))
		i += copy(dAtA[i:], m.PathPrefix)
	}
	if m.Format != 0 {
		dAtA[i] = 0x10
		i++
		i = encodeVarintCapture(dAtA, i, uint64(m.Format))
	}
	if m.XXX_unrecognized != nil {
		i += copy(dAtA[i:], m.XXX_unrecognized)
	}
	return i, nil
}

func (m *Capture) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Capture) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.SinkSelector != nil {
		nn1, err := m.SinkSelector.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += nn1
	}
	if m.TransportSocket != nil {
		dAtA[i] = 0x12
		i++
		i = encodeVarintCapture(dAtA, i, uint64(m.TransportSocket.Size()))
		n2, err := m.TransportSocket.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n2
	}
	if m.XXX_unrecognized != nil {
		i += copy(dAtA[i:], m.XXX_unrecognized)
	}
	return i, nil
}

func (m *Capture_FileSink) MarshalTo(dAtA []byte) (int, error) {
	i := 0
	if m.FileSink != nil {
		dAtA[i] = 0xa
		i++
		i = encodeVarintCapture(dAtA, i, uint64(m.FileSink.Size()))
		n3, err := m.FileSink.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n3
	}
	return i, nil
}
func encodeVarintCapture(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *FileSink) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.PathPrefix)
	if l > 0 {
		n += 1 + l + sovCapture(uint64(l))
	}
	if m.Format != 0 {
		n += 1 + sovCapture(uint64(m.Format))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *Capture) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.SinkSelector != nil {
		n += m.SinkSelector.Size()
	}
	if m.TransportSocket != nil {
		l = m.TransportSocket.Size()
		n += 1 + l + sovCapture(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *Capture_FileSink) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.FileSink != nil {
		l = m.FileSink.Size()
		n += 1 + l + sovCapture(uint64(l))
	}
	return n
}

func sovCapture(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozCapture(x uint64) (n int) {
	return sovCapture(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *FileSink) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowCapture
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: FileSink: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: FileSink: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PathPrefix", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCapture
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthCapture
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.PathPrefix = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Format", wireType)
			}
			m.Format = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCapture
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Format |= (FileSink_Format(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipCapture(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthCapture
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Capture) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowCapture
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Capture: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Capture: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field FileSink", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCapture
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthCapture
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			v := &FileSink{}
			if err := v.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			m.SinkSelector = &Capture_FileSink{v}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TransportSocket", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCapture
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthCapture
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.TransportSocket == nil {
				m.TransportSocket = &core.TransportSocket{}
			}
			if err := m.TransportSocket.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipCapture(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthCapture
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipCapture(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowCapture
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
					return 0, ErrIntOverflowCapture
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
			return iNdEx, nil
		case 1:
			iNdEx += 8
			return iNdEx, nil
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowCapture
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
			iNdEx += length
			if length < 0 {
				return 0, ErrInvalidLengthCapture
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowCapture
					}
					if iNdEx >= l {
						return 0, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					innerWire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				innerWireType := int(innerWire & 0x7)
				if innerWireType == 4 {
					break
				}
				next, err := skipCapture(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
			}
			return iNdEx, nil
		case 4:
			return iNdEx, nil
		case 5:
			iNdEx += 4
			return iNdEx, nil
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
	}
	panic("unreachable")
}

var (
	ErrInvalidLengthCapture = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowCapture   = fmt.Errorf("proto: integer overflow")
)

func init() {
	proto.RegisterFile("envoy/config/transport_socket/capture/v2alpha/capture.proto", fileDescriptor_capture_910b340f585e95f6)
}

var fileDescriptor_capture_910b340f585e95f6 = []byte{
	// 328 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x91, 0x41, 0x4b, 0xc3, 0x30,
	0x18, 0x86, 0x97, 0x1d, 0xe6, 0xf6, 0x4d, 0xb7, 0xd2, 0xd3, 0x10, 0x99, 0xa3, 0xa7, 0x21, 0x98,
	0x40, 0x3c, 0x78, 0x10, 0x04, 0x27, 0x0e, 0x3d, 0xe8, 0x46, 0x57, 0x86, 0x7a, 0x29, 0x59, 0x49,
	0x5d, 0x58, 0x6d, 0x42, 0x1a, 0x8b, 0xfe, 0x2a, 0x7f, 0x83, 0x37, 0x8f, 0xfe, 0x04, 0xd9, 0x2f,
	0x91, 0xa5, 0xed, 0xc1, 0xdd, 0xf4, 0x98, 0x97, 0xf7, 0x7d, 0xbe, 0x37, 0xdf, 0x07, 0x67, 0x3c,
	0xcd, 0xe5, 0x1b, 0x89, 0x64, 0x1a, 0x8b, 0x27, 0x62, 0x34, 0x4b, 0x33, 0x25, 0xb5, 0x09, 0x33,
	0x19, 0xad, 0xb8, 0x21, 0x11, 0x53, 0xe6, 0x45, 0x73, 0x92, 0x53, 0x96, 0xa8, 0x25, 0xab, 0xde,
	0x58, 0x69, 0x69, 0xa4, 0x7b, 0x6c, 0xc3, 0xb8, 0x08, 0xe3, 0xed, 0x30, 0xae, 0xcc, 0x65, 0x78,
	0xff, 0xa0, 0x98, 0xc5, 0x94, 0x20, 0x39, 0x25, 0x91, 0xd4, 0x9c, 0x2c, 0x58, 0x56, 0xc2, 0xbc,
	0x77, 0x04, 0xcd, 0xb1, 0x48, 0xf8, 0x4c, 0xa4, 0x2b, 0xf7, 0x10, 0xda, 0x8a, 0x99, 0x65, 0xa8,
	0x34, 0x8f, 0xc5, 0x6b, 0x0f, 0x0d, 0xd0, 0xb0, 0xe5, 0xc3, 0x46, 0x9a, 0x5a, 0xc5, 0x9d, 0x43,
	0x23, 0x96, 0xfa, 0x99, 0x99, 0x5e, 0x7d, 0x80, 0x86, 0x1d, 0x7a, 0x8e, 0xff, 0xd4, 0x05, 0x57,
	0x93, 0xf0, 0xd8, 0x52, 0xfc, 0x92, 0xe6, 0x1d, 0x41, 0xa3, 0x50, 0x5c, 0x07, 0x76, 0xa7, 0xfe,
	0x24, 0x98, 0x84, 0xa3, 0x9b, 0xbb, 0x0b, 0xff, 0xc1, 0xa9, 0xb9, 0x1d, 0x80, 0x42, 0x09, 0xae,
	0xee, 0x03, 0x07, 0x79, 0x1f, 0x08, 0x76, 0x2e, 0x0b, 0xae, 0x3b, 0x87, 0x56, 0x2c, 0x12, 0x1e,
	0x66, 0x22, 0x5d, 0xd9, 0xba, 0x6d, 0x7a, 0xfa, 0xcf, 0x4a, 0xd7, 0x35, 0xbf, 0x19, 0x57, 0x8b,
	0xb8, 0x05, 0x67, 0x3b, 0x68, 0x7f, 0xdc, 0xa6, 0x5e, 0x89, 0x67, 0x4a, 0xe0, 0x9c, 0xe2, 0xcd,
	0x3a, 0x71, 0x50, 0x59, 0x67, 0xd6, 0xe9, 0x77, 0xcd, 0x6f, 0x61, 0xd4, 0x85, 0xbd, 0x4d, 0xc3,
	0x30, 0xe3, 0x09, 0x8f, 0x8c, 0xd4, 0x23, 0xe7, 0x73, 0xdd, 0x47, 0x5f, 0xeb, 0x3e, 0xfa, 0x5e,
	0xf7, 0xd1, 0x63, 0x3d, 0xa7, 0x8b, 0x86, 0x3d, 0xc7, 0xc9, 0x4f, 0x00, 0x00, 0x00, 0xff, 0xff,
	0xf7, 0xc5, 0xf0, 0xf4, 0x1a, 0x02, 0x00, 0x00,
}
