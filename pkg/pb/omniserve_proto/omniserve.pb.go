// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v5.27.3
// source: omniserve_proto/omniserve.proto

package omniserve_proto

import (
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

type FileChunk struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ProjectCode string `protobuf:"bytes,1,opt,name=project_code,json=projectCode,proto3" json:"project_code,omitempty"`
	FilePath    string `protobuf:"bytes,2,opt,name=file_path,json=filePath,proto3" json:"file_path,omitempty"`
	Content     []byte `protobuf:"bytes,3,opt,name=content,proto3" json:"content,omitempty"`
}

func (x *FileChunk) Reset() {
	*x = FileChunk{}
	if protoimpl.UnsafeEnabled {
		mi := &file_omniserve_proto_omniserve_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FileChunk) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FileChunk) ProtoMessage() {}

func (x *FileChunk) ProtoReflect() protoreflect.Message {
	mi := &file_omniserve_proto_omniserve_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FileChunk.ProtoReflect.Descriptor instead.
func (*FileChunk) Descriptor() ([]byte, []int) {
	return file_omniserve_proto_omniserve_proto_rawDescGZIP(), []int{0}
}

func (x *FileChunk) GetProjectCode() string {
	if x != nil {
		return x.ProjectCode
	}
	return ""
}

func (x *FileChunk) GetFilePath() string {
	if x != nil {
		return x.FilePath
	}
	return ""
}

func (x *FileChunk) GetContent() []byte {
	if x != nil {
		return x.Content
	}
	return nil
}

type PushResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Message string `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *PushResponse) Reset() {
	*x = PushResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_omniserve_proto_omniserve_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PushResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PushResponse) ProtoMessage() {}

func (x *PushResponse) ProtoReflect() protoreflect.Message {
	mi := &file_omniserve_proto_omniserve_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PushResponse.ProtoReflect.Descriptor instead.
func (*PushResponse) Descriptor() ([]byte, []int) {
	return file_omniserve_proto_omniserve_proto_rawDescGZIP(), []int{1}
}

func (x *PushResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

var File_omniserve_proto_omniserve_proto protoreflect.FileDescriptor

var file_omniserve_proto_omniserve_proto_rawDesc = []byte{
	0x0a, 0x1f, 0x6f, 0x6d, 0x6e, 0x69, 0x73, 0x65, 0x72, 0x76, 0x65, 0x5f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2f, 0x6f, 0x6d, 0x6e, 0x69, 0x73, 0x65, 0x72, 0x76, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x09, 0x6f, 0x6d, 0x6e, 0x69, 0x73, 0x65, 0x72, 0x76, 0x65, 0x22, 0x65, 0x0a, 0x09,
	0x46, 0x69, 0x6c, 0x65, 0x43, 0x68, 0x75, 0x6e, 0x6b, 0x12, 0x21, 0x0a, 0x0c, 0x70, 0x72, 0x6f,
	0x6a, 0x65, 0x63, 0x74, 0x5f, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0b, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x1b, 0x0a, 0x09,
	0x66, 0x69, 0x6c, 0x65, 0x5f, 0x70, 0x61, 0x74, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x08, 0x66, 0x69, 0x6c, 0x65, 0x50, 0x61, 0x74, 0x68, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f, 0x6e,
	0x74, 0x65, 0x6e, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74,
	0x65, 0x6e, 0x74, 0x22, 0x28, 0x0a, 0x0c, 0x50, 0x75, 0x73, 0x68, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x32, 0x4b, 0x0a,
	0x09, 0x4f, 0x6d, 0x6e, 0x69, 0x53, 0x65, 0x72, 0x76, 0x65, 0x12, 0x3e, 0x0a, 0x09, 0x50, 0x75,
	0x73, 0x68, 0x46, 0x69, 0x6c, 0x65, 0x73, 0x12, 0x14, 0x2e, 0x6f, 0x6d, 0x6e, 0x69, 0x73, 0x65,
	0x72, 0x76, 0x65, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x43, 0x68, 0x75, 0x6e, 0x6b, 0x1a, 0x17, 0x2e,
	0x6f, 0x6d, 0x6e, 0x69, 0x73, 0x65, 0x72, 0x76, 0x65, 0x2e, 0x50, 0x75, 0x73, 0x68, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x28, 0x01, 0x42, 0x3d, 0x5a, 0x3b, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x65, 0x6d, 0x75, 0x74, 0x68, 0x69, 0x61,
	0x6e, 0x69, 0x6d, 0x62, 0x69, 0x74, 0x68, 0x69, 0x2f, 0x4f, 0x6d, 0x6e, 0x69, 0x53, 0x65, 0x72,
	0x76, 0x65, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x70, 0x62, 0x2f, 0x6f, 0x6d, 0x6e, 0x69, 0x73, 0x65,
	0x72, 0x76, 0x65, 0x5f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_omniserve_proto_omniserve_proto_rawDescOnce sync.Once
	file_omniserve_proto_omniserve_proto_rawDescData = file_omniserve_proto_omniserve_proto_rawDesc
)

func file_omniserve_proto_omniserve_proto_rawDescGZIP() []byte {
	file_omniserve_proto_omniserve_proto_rawDescOnce.Do(func() {
		file_omniserve_proto_omniserve_proto_rawDescData = protoimpl.X.CompressGZIP(file_omniserve_proto_omniserve_proto_rawDescData)
	})
	return file_omniserve_proto_omniserve_proto_rawDescData
}

var file_omniserve_proto_omniserve_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_omniserve_proto_omniserve_proto_goTypes = []interface{}{
	(*FileChunk)(nil),    // 0: omniserve.FileChunk
	(*PushResponse)(nil), // 1: omniserve.PushResponse
}
var file_omniserve_proto_omniserve_proto_depIdxs = []int32{
	0, // 0: omniserve.OmniServe.PushFiles:input_type -> omniserve.FileChunk
	1, // 1: omniserve.OmniServe.PushFiles:output_type -> omniserve.PushResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_omniserve_proto_omniserve_proto_init() }
func file_omniserve_proto_omniserve_proto_init() {
	if File_omniserve_proto_omniserve_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_omniserve_proto_omniserve_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FileChunk); i {
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
		file_omniserve_proto_omniserve_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PushResponse); i {
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
			RawDescriptor: file_omniserve_proto_omniserve_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_omniserve_proto_omniserve_proto_goTypes,
		DependencyIndexes: file_omniserve_proto_omniserve_proto_depIdxs,
		MessageInfos:      file_omniserve_proto_omniserve_proto_msgTypes,
	}.Build()
	File_omniserve_proto_omniserve_proto = out.File
	file_omniserve_proto_omniserve_proto_rawDesc = nil
	file_omniserve_proto_omniserve_proto_goTypes = nil
	file_omniserve_proto_omniserve_proto_depIdxs = nil
}
