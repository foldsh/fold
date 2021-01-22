// This defines a service manifest. It is how the runtime is able to introspect
// the structure of a service and figure out what infrastructure is required
// in order to operate it.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.14.0
// source: manifest.proto

package manifest

import (
	proto "github.com/golang/protobuf/proto"
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

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type HttpMethod int32

const (
	HttpMethod_GET    HttpMethod = 0
	HttpMethod_PUT    HttpMethod = 1
	HttpMethod_POST   HttpMethod = 2
	HttpMethod_DELETE HttpMethod = 3
)

// Enum value maps for HttpMethod.
var (
	HttpMethod_name = map[int32]string{
		0: "GET",
		1: "PUT",
		2: "POST",
		3: "DELETE",
	}
	HttpMethod_value = map[string]int32{
		"GET":    0,
		"PUT":    1,
		"POST":   2,
		"DELETE": 3,
	}
)

func (x HttpMethod) Enum() *HttpMethod {
	p := new(HttpMethod)
	*p = x
	return p
}

func (x HttpMethod) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (HttpMethod) Descriptor() protoreflect.EnumDescriptor {
	return file_manifest_proto_enumTypes[0].Descriptor()
}

func (HttpMethod) Type() protoreflect.EnumType {
	return &file_manifest_proto_enumTypes[0]
}

func (x HttpMethod) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use HttpMethod.Descriptor instead.
func (HttpMethod) EnumDescriptor() ([]byte, []int) {
	return file_manifest_proto_rawDescGZIP(), []int{0}
}

// A manifest describing everything required to build and deploy a service.
type Manifest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The name of the service.
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// The version of the service.
	Version *Version `protobuf:"bytes,2,opt,name=version,proto3" json:"version,omitempty"`
	// Information required to build the service. This will be filled in
	// by the command line tool.
	BuildInfo *BuildInfo `protobuf:"bytes,3,opt,name=build_info,json=buildInfo,proto3" json:"build_info,omitempty"`
	// The routes defined by the router within the service.
	Routes []*Route `protobuf:"bytes,4,rep,name=routes,proto3" json:"routes,omitempty"`
}

func (x *Manifest) Reset() {
	*x = Manifest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_manifest_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Manifest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Manifest) ProtoMessage() {}

func (x *Manifest) ProtoReflect() protoreflect.Message {
	mi := &file_manifest_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Manifest.ProtoReflect.Descriptor instead.
func (*Manifest) Descriptor() ([]byte, []int) {
	return file_manifest_proto_rawDescGZIP(), []int{0}
}

func (x *Manifest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Manifest) GetVersion() *Version {
	if x != nil {
		return x.Version
	}
	return nil
}

func (x *Manifest) GetBuildInfo() *BuildInfo {
	if x != nil {
		return x.BuildInfo
	}
	return nil
}

func (x *Manifest) GetRoutes() []*Route {
	if x != nil {
		return x.Routes
	}
	return nil
}

type BuildInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The maintainer of the service.
	Maintainer string `protobuf:"bytes,1,opt,name=maintainer,proto3" json:"maintainer,omitempty"`
	// The image name.
	Image string `protobuf:"bytes,2,opt,name=image,proto3" json:"image,omitempty"`
	// The image tag.
	Tag string `protobuf:"bytes,3,opt,name=tag,proto3" json:"tag,omitempty"`
	// The path to the service, relative to the project root.
	Path string `protobuf:"bytes,4,opt,name=path,proto3" json:"path,omitempty"`
}

func (x *BuildInfo) Reset() {
	*x = BuildInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_manifest_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BuildInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BuildInfo) ProtoMessage() {}

func (x *BuildInfo) ProtoReflect() protoreflect.Message {
	mi := &file_manifest_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BuildInfo.ProtoReflect.Descriptor instead.
func (*BuildInfo) Descriptor() ([]byte, []int) {
	return file_manifest_proto_rawDescGZIP(), []int{1}
}

func (x *BuildInfo) GetMaintainer() string {
	if x != nil {
		return x.Maintainer
	}
	return ""
}

func (x *BuildInfo) GetImage() string {
	if x != nil {
		return x.Image
	}
	return ""
}

func (x *BuildInfo) GetTag() string {
	if x != nil {
		return x.Tag
	}
	return ""
}

func (x *BuildInfo) GetPath() string {
	if x != nil {
		return x.Path
	}
	return ""
}

// A SemVer version number.
type Version struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Major
	Major int32 `protobuf:"varint,1,opt,name=major,proto3" json:"major,omitempty"`
	// Minor
	Minor int32 `protobuf:"varint,2,opt,name=minor,proto3" json:"minor,omitempty"`
	// Patch
	Patch int32 `protobuf:"varint,3,opt,name=patch,proto3" json:"patch,omitempty"`
}

func (x *Version) Reset() {
	*x = Version{}
	if protoimpl.UnsafeEnabled {
		mi := &file_manifest_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Version) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Version) ProtoMessage() {}

func (x *Version) ProtoReflect() protoreflect.Message {
	mi := &file_manifest_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Version.ProtoReflect.Descriptor instead.
func (*Version) Descriptor() ([]byte, []int) {
	return file_manifest_proto_rawDescGZIP(), []int{2}
}

func (x *Version) GetMajor() int32 {
	if x != nil {
		return x.Major
	}
	return 0
}

func (x *Version) GetMinor() int32 {
	if x != nil {
		return x.Minor
	}
	return 0
}

func (x *Version) GetPatch() int32 {
	if x != nil {
		return x.Patch
	}
	return 0
}

// A route defined by a service.
type Route struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The HTTP method which this Route handles.
	HttpMethod HttpMethod `protobuf:"varint,1,opt,name=http_method,json=httpMethod,proto3,enum=manifest.HttpMethod" json:"http_method,omitempty"`
	// The identifier of the handler that will handle the route.
	Handler string `protobuf:"bytes,2,opt,name=handler,proto3" json:"handler,omitempty"`
	// The path pattern that will match this route.
	PathSpec string `protobuf:"bytes,3,opt,name=path_spec,json=pathSpec,proto3" json:"path_spec,omitempty"`
}

func (x *Route) Reset() {
	*x = Route{}
	if protoimpl.UnsafeEnabled {
		mi := &file_manifest_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Route) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Route) ProtoMessage() {}

func (x *Route) ProtoReflect() protoreflect.Message {
	mi := &file_manifest_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Route.ProtoReflect.Descriptor instead.
func (*Route) Descriptor() ([]byte, []int) {
	return file_manifest_proto_rawDescGZIP(), []int{3}
}

func (x *Route) GetHttpMethod() HttpMethod {
	if x != nil {
		return x.HttpMethod
	}
	return HttpMethod_GET
}

func (x *Route) GetHandler() string {
	if x != nil {
		return x.Handler
	}
	return ""
}

func (x *Route) GetPathSpec() string {
	if x != nil {
		return x.PathSpec
	}
	return ""
}

var File_manifest_proto protoreflect.FileDescriptor

var file_manifest_proto_rawDesc = []byte{
	0x0a, 0x0e, 0x6d, 0x61, 0x6e, 0x69, 0x66, 0x65, 0x73, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x08, 0x6d, 0x61, 0x6e, 0x69, 0x66, 0x65, 0x73, 0x74, 0x22, 0xa8, 0x01, 0x0a, 0x08, 0x4d,
	0x61, 0x6e, 0x69, 0x66, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x2b, 0x0a, 0x07, 0x76,
	0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x6d,
	0x61, 0x6e, 0x69, 0x66, 0x65, 0x73, 0x74, 0x2e, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x52,
	0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x32, 0x0a, 0x0a, 0x62, 0x75, 0x69, 0x6c,
	0x64, 0x5f, 0x69, 0x6e, 0x66, 0x6f, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x6d,
	0x61, 0x6e, 0x69, 0x66, 0x65, 0x73, 0x74, 0x2e, 0x42, 0x75, 0x69, 0x6c, 0x64, 0x49, 0x6e, 0x66,
	0x6f, 0x52, 0x09, 0x62, 0x75, 0x69, 0x6c, 0x64, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x27, 0x0a, 0x06,
	0x72, 0x6f, 0x75, 0x74, 0x65, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x6d,
	0x61, 0x6e, 0x69, 0x66, 0x65, 0x73, 0x74, 0x2e, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x52, 0x06, 0x72,
	0x6f, 0x75, 0x74, 0x65, 0x73, 0x22, 0x67, 0x0a, 0x09, 0x42, 0x75, 0x69, 0x6c, 0x64, 0x49, 0x6e,
	0x66, 0x6f, 0x12, 0x1e, 0x0a, 0x0a, 0x6d, 0x61, 0x69, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65, 0x72,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x6d, 0x61, 0x69, 0x6e, 0x74, 0x61, 0x69, 0x6e,
	0x65, 0x72, 0x12, 0x14, 0x0a, 0x05, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x74, 0x61, 0x67, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x74, 0x61, 0x67, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x61,
	0x74, 0x68, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x70, 0x61, 0x74, 0x68, 0x22, 0x4b,
	0x0a, 0x07, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x14, 0x0a, 0x05, 0x6d, 0x61, 0x6a,
	0x6f, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x6d, 0x61, 0x6a, 0x6f, 0x72, 0x12,
	0x14, 0x0a, 0x05, 0x6d, 0x69, 0x6e, 0x6f, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05,
	0x6d, 0x69, 0x6e, 0x6f, 0x72, 0x12, 0x14, 0x0a, 0x05, 0x70, 0x61, 0x74, 0x63, 0x68, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x70, 0x61, 0x74, 0x63, 0x68, 0x22, 0x75, 0x0a, 0x05, 0x52,
	0x6f, 0x75, 0x74, 0x65, 0x12, 0x35, 0x0a, 0x0b, 0x68, 0x74, 0x74, 0x70, 0x5f, 0x6d, 0x65, 0x74,
	0x68, 0x6f, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x14, 0x2e, 0x6d, 0x61, 0x6e, 0x69,
	0x66, 0x65, 0x73, 0x74, 0x2e, 0x48, 0x74, 0x74, 0x70, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x52,
	0x0a, 0x68, 0x74, 0x74, 0x70, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x68,
	0x61, 0x6e, 0x64, 0x6c, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x68, 0x61,
	0x6e, 0x64, 0x6c, 0x65, 0x72, 0x12, 0x1b, 0x0a, 0x09, 0x70, 0x61, 0x74, 0x68, 0x5f, 0x73, 0x70,
	0x65, 0x63, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x70, 0x61, 0x74, 0x68, 0x53, 0x70,
	0x65, 0x63, 0x2a, 0x34, 0x0a, 0x0a, 0x48, 0x74, 0x74, 0x70, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64,
	0x12, 0x07, 0x0a, 0x03, 0x47, 0x45, 0x54, 0x10, 0x00, 0x12, 0x07, 0x0a, 0x03, 0x50, 0x55, 0x54,
	0x10, 0x01, 0x12, 0x08, 0x0a, 0x04, 0x50, 0x4f, 0x53, 0x54, 0x10, 0x02, 0x12, 0x0a, 0x0a, 0x06,
	0x44, 0x45, 0x4c, 0x45, 0x54, 0x45, 0x10, 0x03, 0x42, 0x21, 0x5a, 0x1f, 0x67, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x66, 0x6f, 0x6c, 0x64, 0x73, 0x68, 0x2f, 0x66, 0x6f,
	0x6c, 0x64, 0x2f, 0x6d, 0x61, 0x6e, 0x69, 0x66, 0x65, 0x73, 0x74, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_manifest_proto_rawDescOnce sync.Once
	file_manifest_proto_rawDescData = file_manifest_proto_rawDesc
)

func file_manifest_proto_rawDescGZIP() []byte {
	file_manifest_proto_rawDescOnce.Do(func() {
		file_manifest_proto_rawDescData = protoimpl.X.CompressGZIP(file_manifest_proto_rawDescData)
	})
	return file_manifest_proto_rawDescData
}

var file_manifest_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_manifest_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_manifest_proto_goTypes = []interface{}{
	(HttpMethod)(0),   // 0: manifest.HttpMethod
	(*Manifest)(nil),  // 1: manifest.Manifest
	(*BuildInfo)(nil), // 2: manifest.BuildInfo
	(*Version)(nil),   // 3: manifest.Version
	(*Route)(nil),     // 4: manifest.Route
}
var file_manifest_proto_depIdxs = []int32{
	3, // 0: manifest.Manifest.version:type_name -> manifest.Version
	2, // 1: manifest.Manifest.build_info:type_name -> manifest.BuildInfo
	4, // 2: manifest.Manifest.routes:type_name -> manifest.Route
	0, // 3: manifest.Route.http_method:type_name -> manifest.HttpMethod
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_manifest_proto_init() }
func file_manifest_proto_init() {
	if File_manifest_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_manifest_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Manifest); i {
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
		file_manifest_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BuildInfo); i {
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
		file_manifest_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Version); i {
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
		file_manifest_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Route); i {
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
			RawDescriptor: file_manifest_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_manifest_proto_goTypes,
		DependencyIndexes: file_manifest_proto_depIdxs,
		EnumInfos:         file_manifest_proto_enumTypes,
		MessageInfos:      file_manifest_proto_msgTypes,
	}.Build()
	File_manifest_proto = out.File
	file_manifest_proto_rawDesc = nil
	file_manifest_proto_goTypes = nil
	file_manifest_proto_depIdxs = nil
}
