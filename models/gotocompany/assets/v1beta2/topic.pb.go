// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        (unknown)
// source: gotocompany/assets/v1beta2/topic.proto

package assetsv1beta2

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	structpb "google.golang.org/protobuf/types/known/structpb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Topic is resource that represents a logical group of messages
// in message bus like kafka, pubsub, pulsar etc.
type Topic struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The metrics of the topic.
	// For an example check out topic profile schema.
	Profile *TopicProfile `protobuf:"bytes,1,opt,name=profile,proto3" json:"profile,omitempty"`
	// The schema of the topic.
	// For an example check out topic schema.
	Schema *TopicSchema `protobuf:"bytes,2,opt,name=schema,proto3" json:"schema,omitempty"`
	// List of attributes the model has.
	Attributes *structpb.Struct `protobuf:"bytes,10,opt,name=attributes,proto3" json:"attributes,omitempty"`
	// The timestamp of the topic's creation.
	CreateTime *timestamppb.Timestamp `protobuf:"bytes,101,opt,name=create_time,json=createTime,proto3" json:"create_time,omitempty"`
	// The timestamp when the topic was last modified.
	UpdateTime *timestamppb.Timestamp `protobuf:"bytes,102,opt,name=update_time,json=updateTime,proto3" json:"update_time,omitempty"`
}

func (x *Topic) Reset() {
	*x = Topic{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gotocompany_assets_v1beta2_topic_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Topic) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Topic) ProtoMessage() {}

func (x *Topic) ProtoReflect() protoreflect.Message {
	mi := &file_gotocompany_assets_v1beta2_topic_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Topic.ProtoReflect.Descriptor instead.
func (*Topic) Descriptor() ([]byte, []int) {
	return file_gotocompany_assets_v1beta2_topic_proto_rawDescGZIP(), []int{0}
}

func (x *Topic) GetProfile() *TopicProfile {
	if x != nil {
		return x.Profile
	}
	return nil
}

func (x *Topic) GetSchema() *TopicSchema {
	if x != nil {
		return x.Schema
	}
	return nil
}

func (x *Topic) GetAttributes() *structpb.Struct {
	if x != nil {
		return x.Attributes
	}
	return nil
}

func (x *Topic) GetCreateTime() *timestamppb.Timestamp {
	if x != nil {
		return x.CreateTime
	}
	return nil
}

func (x *Topic) GetUpdateTime() *timestamppb.Timestamp {
	if x != nil {
		return x.UpdateTime
	}
	return nil
}

// TopicProfile is the profile of the topic.
type TopicProfile struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The throughput of the topic.
	// Example: `1m/minute`.
	Throughput string `protobuf:"bytes,1,opt,name=throughput,proto3" json:"throughput,omitempty"`
	// The number of partitions in the topic.
	// Example: `12`.
	NumberOfPartitions int64 `protobuf:"varint,2,opt,name=number_of_partitions,json=numberOfPartitions,proto3" json:"number_of_partitions,omitempty"`
}

func (x *TopicProfile) Reset() {
	*x = TopicProfile{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gotocompany_assets_v1beta2_topic_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TopicProfile) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TopicProfile) ProtoMessage() {}

func (x *TopicProfile) ProtoReflect() protoreflect.Message {
	mi := &file_gotocompany_assets_v1beta2_topic_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TopicProfile.ProtoReflect.Descriptor instead.
func (*TopicProfile) Descriptor() ([]byte, []int) {
	return file_gotocompany_assets_v1beta2_topic_proto_rawDescGZIP(), []int{1}
}

func (x *TopicProfile) GetThroughput() string {
	if x != nil {
		return x.Throughput
	}
	return ""
}

func (x *TopicProfile) GetNumberOfPartitions() int64 {
	if x != nil {
		return x.NumberOfPartitions
	}
	return 0
}

// TopicSchema represents a schema for message bus.
// It is facet used to specify the schema of a message bus.
type TopicSchema struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SchemaUrl string `protobuf:"bytes,1,opt,name=schema_url,json=schemaUrl,proto3" json:"schema_url,omitempty"`
	Format    string `protobuf:"bytes,2,opt,name=format,proto3" json:"format,omitempty"`
}

func (x *TopicSchema) Reset() {
	*x = TopicSchema{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gotocompany_assets_v1beta2_topic_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TopicSchema) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TopicSchema) ProtoMessage() {}

func (x *TopicSchema) ProtoReflect() protoreflect.Message {
	mi := &file_gotocompany_assets_v1beta2_topic_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TopicSchema.ProtoReflect.Descriptor instead.
func (*TopicSchema) Descriptor() ([]byte, []int) {
	return file_gotocompany_assets_v1beta2_topic_proto_rawDescGZIP(), []int{2}
}

func (x *TopicSchema) GetSchemaUrl() string {
	if x != nil {
		return x.SchemaUrl
	}
	return ""
}

func (x *TopicSchema) GetFormat() string {
	if x != nil {
		return x.Format
	}
	return ""
}

var File_gotocompany_assets_v1beta2_topic_proto protoreflect.FileDescriptor

var file_gotocompany_assets_v1beta2_topic_proto_rawDesc = []byte{
	0x0a, 0x26, 0x67, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6d, 0x70, 0x61, 0x6e, 0x79, 0x2f, 0x61, 0x73,
	0x73, 0x65, 0x74, 0x73, 0x2f, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x32, 0x2f, 0x74, 0x6f, 0x70,
	0x69, 0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x1a, 0x67, 0x6f, 0x74, 0x6f, 0x63, 0x6f,
	0x6d, 0x70, 0x61, 0x6e, 0x79, 0x2e, 0x61, 0x73, 0x73, 0x65, 0x74, 0x73, 0x2e, 0x76, 0x31, 0x62,
	0x65, 0x74, 0x61, 0x32, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0xbf, 0x02, 0x0a, 0x05, 0x54, 0x6f, 0x70, 0x69, 0x63, 0x12, 0x42, 0x0a,
	0x07, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x28,
	0x2e, 0x67, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6d, 0x70, 0x61, 0x6e, 0x79, 0x2e, 0x61, 0x73, 0x73,
	0x65, 0x74, 0x73, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x32, 0x2e, 0x54, 0x6f, 0x70, 0x69,
	0x63, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x52, 0x07, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c,
	0x65, 0x12, 0x3f, 0x0a, 0x06, 0x73, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x27, 0x2e, 0x67, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6d, 0x70, 0x61, 0x6e, 0x79, 0x2e,
	0x61, 0x73, 0x73, 0x65, 0x74, 0x73, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x32, 0x2e, 0x54,
	0x6f, 0x70, 0x69, 0x63, 0x53, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x52, 0x06, 0x73, 0x63, 0x68, 0x65,
	0x6d, 0x61, 0x12, 0x37, 0x0a, 0x0a, 0x61, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x73,
	0x18, 0x0a, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x75, 0x63, 0x74, 0x52,
	0x0a, 0x61, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x73, 0x12, 0x3b, 0x0a, 0x0b, 0x63,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x65, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x0a, 0x63, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x3b, 0x0a, 0x0b, 0x75, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x66, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x54, 0x69, 0x6d, 0x65, 0x22, 0x60, 0x0a, 0x0c, 0x54, 0x6f, 0x70, 0x69, 0x63, 0x50, 0x72,
	0x6f, 0x66, 0x69, 0x6c, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x74, 0x68, 0x72, 0x6f, 0x75, 0x67, 0x68,
	0x70, 0x75, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x74, 0x68, 0x72, 0x6f, 0x75,
	0x67, 0x68, 0x70, 0x75, 0x74, 0x12, 0x30, 0x0a, 0x14, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x5f,
	0x6f, 0x66, 0x5f, 0x70, 0x61, 0x72, 0x74, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x12, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x4f, 0x66, 0x50, 0x61, 0x72,
	0x74, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x22, 0x44, 0x0a, 0x0b, 0x54, 0x6f, 0x70, 0x69, 0x63,
	0x53, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x63, 0x68, 0x65, 0x6d, 0x61,
	0x5f, 0x75, 0x72, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x63, 0x68, 0x65,
	0x6d, 0x61, 0x55, 0x72, 0x6c, 0x12, 0x16, 0x0a, 0x06, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x42, 0x59, 0x0a,
	0x16, 0x63, 0x6f, 0x6d, 0x2e, 0x67, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6d, 0x70, 0x61, 0x6e, 0x79,
	0x2e, 0x61, 0x73, 0x73, 0x65, 0x74, 0x73, 0x42, 0x0a, 0x54, 0x6f, 0x70, 0x69, 0x63, 0x50, 0x72,
	0x6f, 0x74, 0x6f, 0x5a, 0x33, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x67, 0x6f, 0x74, 0x6f, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x6e, 0x2f, 0x61, 0x73, 0x73, 0x65,
	0x74, 0x73, 0x2f, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x32, 0x3b, 0x61, 0x73, 0x73, 0x65, 0x74,
	0x73, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x32, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_gotocompany_assets_v1beta2_topic_proto_rawDescOnce sync.Once
	file_gotocompany_assets_v1beta2_topic_proto_rawDescData = file_gotocompany_assets_v1beta2_topic_proto_rawDesc
)

func file_gotocompany_assets_v1beta2_topic_proto_rawDescGZIP() []byte {
	file_gotocompany_assets_v1beta2_topic_proto_rawDescOnce.Do(func() {
		file_gotocompany_assets_v1beta2_topic_proto_rawDescData = protoimpl.X.CompressGZIP(file_gotocompany_assets_v1beta2_topic_proto_rawDescData)
	})
	return file_gotocompany_assets_v1beta2_topic_proto_rawDescData
}

var file_gotocompany_assets_v1beta2_topic_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_gotocompany_assets_v1beta2_topic_proto_goTypes = []interface{}{
	(*Topic)(nil),                 // 0: gotocompany.assets.v1beta2.Topic
	(*TopicProfile)(nil),          // 1: gotocompany.assets.v1beta2.TopicProfile
	(*TopicSchema)(nil),           // 2: gotocompany.assets.v1beta2.TopicSchema
	(*structpb.Struct)(nil),       // 3: google.protobuf.Struct
	(*timestamppb.Timestamp)(nil), // 4: google.protobuf.Timestamp
}
var file_gotocompany_assets_v1beta2_topic_proto_depIdxs = []int32{
	1, // 0: gotocompany.assets.v1beta2.Topic.profile:type_name -> gotocompany.assets.v1beta2.TopicProfile
	2, // 1: gotocompany.assets.v1beta2.Topic.schema:type_name -> gotocompany.assets.v1beta2.TopicSchema
	3, // 2: gotocompany.assets.v1beta2.Topic.attributes:type_name -> google.protobuf.Struct
	4, // 3: gotocompany.assets.v1beta2.Topic.create_time:type_name -> google.protobuf.Timestamp
	4, // 4: gotocompany.assets.v1beta2.Topic.update_time:type_name -> google.protobuf.Timestamp
	5, // [5:5] is the sub-list for method output_type
	5, // [5:5] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_gotocompany_assets_v1beta2_topic_proto_init() }
func file_gotocompany_assets_v1beta2_topic_proto_init() {
	if File_gotocompany_assets_v1beta2_topic_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_gotocompany_assets_v1beta2_topic_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Topic); i {
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
		file_gotocompany_assets_v1beta2_topic_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TopicProfile); i {
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
		file_gotocompany_assets_v1beta2_topic_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TopicSchema); i {
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
			RawDescriptor: file_gotocompany_assets_v1beta2_topic_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_gotocompany_assets_v1beta2_topic_proto_goTypes,
		DependencyIndexes: file_gotocompany_assets_v1beta2_topic_proto_depIdxs,
		MessageInfos:      file_gotocompany_assets_v1beta2_topic_proto_msgTypes,
	}.Build()
	File_gotocompany_assets_v1beta2_topic_proto = out.File
	file_gotocompany_assets_v1beta2_topic_proto_rawDesc = nil
	file_gotocompany_assets_v1beta2_topic_proto_goTypes = nil
	file_gotocompany_assets_v1beta2_topic_proto_depIdxs = nil
}
