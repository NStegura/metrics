// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v3.6.1
// source: metricsapi.proto

package api

import (
	empty "github.com/golang/protobuf/ptypes/empty"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type MetricType int32

const (
	MetricType_GAUGE   MetricType = 0
	MetricType_COUNTER MetricType = 1
)

// Enum value maps for MetricType.
var (
	MetricType_name = map[int32]string{
		0: "GAUGE",
		1: "COUNTER",
	}
	MetricType_value = map[string]int32{
		"GAUGE":   0,
		"COUNTER": 1,
	}
)

func (x MetricType) Enum() *MetricType {
	p := new(MetricType)
	*p = x
	return p
}

func (x MetricType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (MetricType) Descriptor() protoreflect.EnumDescriptor {
	return file_metricsapi_proto_enumTypes[0].Descriptor()
}

func (MetricType) Type() protoreflect.EnumType {
	return &file_metricsapi_proto_enumTypes[0]
}

func (x MetricType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use MetricType.Descriptor instead.
func (MetricType) EnumDescriptor() ([]byte, []int) {
	return file_metricsapi_proto_rawDescGZIP(), []int{0}
}

type MetricsList struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Metrics []*Metric `protobuf:"bytes,1,rep,name=metrics,proto3" json:"metrics,omitempty"`
}

func (x *MetricsList) Reset() {
	*x = MetricsList{}
	if protoimpl.UnsafeEnabled {
		mi := &file_metricsapi_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MetricsList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MetricsList) ProtoMessage() {}

func (x *MetricsList) ProtoReflect() protoreflect.Message {
	mi := &file_metricsapi_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MetricsList.ProtoReflect.Descriptor instead.
func (*MetricsList) Descriptor() ([]byte, []int) {
	return file_metricsapi_proto_rawDescGZIP(), []int{0}
}

func (x *MetricsList) GetMetrics() []*Metric {
	if x != nil {
		return x.Metrics
	}
	return nil
}

type Metric struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id    string     `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Mtype MetricType `protobuf:"varint,2,opt,name=mtype,proto3,enum=metricsapi.MetricType" json:"mtype,omitempty"`
	Value float64    `protobuf:"fixed64,3,opt,name=value,proto3" json:"value,omitempty"`
	Delta int64      `protobuf:"varint,4,opt,name=delta,proto3" json:"delta,omitempty"`
}

func (x *Metric) Reset() {
	*x = Metric{}
	if protoimpl.UnsafeEnabled {
		mi := &file_metricsapi_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Metric) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Metric) ProtoMessage() {}

func (x *Metric) ProtoReflect() protoreflect.Message {
	mi := &file_metricsapi_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Metric.ProtoReflect.Descriptor instead.
func (*Metric) Descriptor() ([]byte, []int) {
	return file_metricsapi_proto_rawDescGZIP(), []int{1}
}

func (x *Metric) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Metric) GetMtype() MetricType {
	if x != nil {
		return x.Mtype
	}
	return MetricType_GAUGE
}

func (x *Metric) GetValue() float64 {
	if x != nil {
		return x.Value
	}
	return 0
}

func (x *Metric) GetDelta() int64 {
	if x != nil {
		return x.Delta
	}
	return 0
}

type UpdateResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Message string `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *UpdateResponse) Reset() {
	*x = UpdateResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_metricsapi_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateResponse) ProtoMessage() {}

func (x *UpdateResponse) ProtoReflect() protoreflect.Message {
	mi := &file_metricsapi_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateResponse.ProtoReflect.Descriptor instead.
func (*UpdateResponse) Descriptor() ([]byte, []int) {
	return file_metricsapi_proto_rawDescGZIP(), []int{2}
}

func (x *UpdateResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type Pong struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Pong bool `protobuf:"varint,1,opt,name=pong,proto3" json:"pong,omitempty"`
}

func (x *Pong) Reset() {
	*x = Pong{}
	if protoimpl.UnsafeEnabled {
		mi := &file_metricsapi_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Pong) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Pong) ProtoMessage() {}

func (x *Pong) ProtoReflect() protoreflect.Message {
	mi := &file_metricsapi_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Pong.ProtoReflect.Descriptor instead.
func (*Pong) Descriptor() ([]byte, []int) {
	return file_metricsapi_proto_rawDescGZIP(), []int{3}
}

func (x *Pong) GetPong() bool {
	if x != nil {
		return x.Pong
	}
	return false
}

var File_metricsapi_proto protoreflect.FileDescriptor

var file_metricsapi_proto_rawDesc = []byte{
	0x0a, 0x10, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x61, 0x70, 0x69, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x0a, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x61, 0x70, 0x69, 0x1a, 0x1b,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f,
	0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x3b, 0x0a, 0x0b, 0x4d,
	0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x2c, 0x0a, 0x07, 0x6d, 0x65,
	0x74, 0x72, 0x69, 0x63, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x6d, 0x65,
	0x74, 0x72, 0x69, 0x63, 0x73, 0x61, 0x70, 0x69, 0x2e, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x52,
	0x07, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x22, 0x72, 0x0a, 0x06, 0x4d, 0x65, 0x74, 0x72,
	0x69, 0x63, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02,
	0x69, 0x64, 0x12, 0x2c, 0x0a, 0x05, 0x6d, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0e, 0x32, 0x16, 0x2e, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x61, 0x70, 0x69, 0x2e, 0x4d,
	0x65, 0x74, 0x72, 0x69, 0x63, 0x54, 0x79, 0x70, 0x65, 0x52, 0x05, 0x6d, 0x74, 0x79, 0x70, 0x65,
	0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x01, 0x52,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x64, 0x65, 0x6c, 0x74, 0x61, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x64, 0x65, 0x6c, 0x74, 0x61, 0x22, 0x2a, 0x0a, 0x0e,
	0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18,
	0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x1a, 0x0a, 0x04, 0x50, 0x6f, 0x6e, 0x67,
	0x12, 0x12, 0x0a, 0x04, 0x70, 0x6f, 0x6e, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x04,
	0x70, 0x6f, 0x6e, 0x67, 0x2a, 0x24, 0x0a, 0x0a, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x54, 0x79,
	0x70, 0x65, 0x12, 0x09, 0x0a, 0x05, 0x47, 0x41, 0x55, 0x47, 0x45, 0x10, 0x00, 0x12, 0x0b, 0x0a,
	0x07, 0x43, 0x4f, 0x55, 0x4e, 0x54, 0x45, 0x52, 0x10, 0x01, 0x32, 0x8e, 0x01, 0x0a, 0x0a, 0x4d,
	0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x41, 0x70, 0x69, 0x12, 0x49, 0x0a, 0x10, 0x55, 0x70, 0x64,
	0x61, 0x74, 0x65, 0x41, 0x6c, 0x6c, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x12, 0x17, 0x2e,
	0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x61, 0x70, 0x69, 0x2e, 0x4d, 0x65, 0x74, 0x72, 0x69,
	0x63, 0x73, 0x4c, 0x69, 0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73,
	0x61, 0x70, 0x69, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x22, 0x00, 0x12, 0x35, 0x0a, 0x07, 0x47, 0x65, 0x74, 0x50, 0x69, 0x6e, 0x67, 0x12,
	0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x10, 0x2e, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63,
	0x73, 0x61, 0x70, 0x69, 0x2e, 0x50, 0x6f, 0x6e, 0x67, 0x22, 0x00, 0x42, 0x21, 0x5a, 0x1f, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x4e, 0x53, 0x74, 0x65, 0x67, 0x75,
	0x72, 0x61, 0x2f, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x2f, 0x61, 0x70, 0x69, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_metricsapi_proto_rawDescOnce sync.Once
	file_metricsapi_proto_rawDescData = file_metricsapi_proto_rawDesc
)

func file_metricsapi_proto_rawDescGZIP() []byte {
	file_metricsapi_proto_rawDescOnce.Do(func() {
		file_metricsapi_proto_rawDescData = protoimpl.X.CompressGZIP(file_metricsapi_proto_rawDescData)
	})
	return file_metricsapi_proto_rawDescData
}

var file_metricsapi_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_metricsapi_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_metricsapi_proto_goTypes = []any{
	(MetricType)(0),        // 0: metricsapi.MetricType
	(*MetricsList)(nil),    // 1: metricsapi.MetricsList
	(*Metric)(nil),         // 2: metricsapi.Metric
	(*UpdateResponse)(nil), // 3: metricsapi.UpdateResponse
	(*Pong)(nil),           // 4: metricsapi.Pong
	(*empty.Empty)(nil),    // 5: google.protobuf.Empty
}
var file_metricsapi_proto_depIdxs = []int32{
	2, // 0: metricsapi.MetricsList.metrics:type_name -> metricsapi.Metric
	0, // 1: metricsapi.Metric.mtype:type_name -> metricsapi.MetricType
	1, // 2: metricsapi.MetricsApi.UpdateAllMetrics:input_type -> metricsapi.MetricsList
	5, // 3: metricsapi.MetricsApi.GetPing:input_type -> google.protobuf.Empty
	3, // 4: metricsapi.MetricsApi.UpdateAllMetrics:output_type -> metricsapi.UpdateResponse
	4, // 5: metricsapi.MetricsApi.GetPing:output_type -> metricsapi.Pong
	4, // [4:6] is the sub-list for method output_type
	2, // [2:4] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_metricsapi_proto_init() }
func file_metricsapi_proto_init() {
	if File_metricsapi_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_metricsapi_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*MetricsList); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_metricsapi_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*Metric); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_metricsapi_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*UpdateResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_metricsapi_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*Pong); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_metricsapi_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_metricsapi_proto_goTypes,
		DependencyIndexes: file_metricsapi_proto_depIdxs,
		EnumInfos:         file_metricsapi_proto_enumTypes,
		MessageInfos:      file_metricsapi_proto_msgTypes,
	}.Build()
	File_metricsapi_proto = out.File
	file_metricsapi_proto_rawDesc = nil
	file_metricsapi_proto_goTypes = nil
	file_metricsapi_proto_depIdxs = nil
}
