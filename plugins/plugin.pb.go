// Code generated by protoc-gen-go.
// source: plugin.proto
// DO NOT EDIT!

/*
Package openapic_v1 is a generated protocol buffer package.

It is generated from these files:
	plugin.proto

It has these top-level messages:
	Version
	PluginRequest
	PluginResponse
	Wrapper
*/
package openapic_v1

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

// The version number of OpenAPI compiler.
type Version struct {
	Major int32 `protobuf:"varint,1,opt,name=major" json:"major,omitempty"`
	Minor int32 `protobuf:"varint,2,opt,name=minor" json:"minor,omitempty"`
	Patch int32 `protobuf:"varint,3,opt,name=patch" json:"patch,omitempty"`
	// A suffix for alpha, beta or rc release, e.g., "alpha-1", "rc2". It should
	// be empty for mainline stable releases.
	Suffix string `protobuf:"bytes,4,opt,name=suffix" json:"suffix,omitempty"`
}

func (m *Version) Reset()                    { *m = Version{} }
func (m *Version) String() string            { return proto.CompactTextString(m) }
func (*Version) ProtoMessage()               {}
func (*Version) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

// An encoded PluginRequest is written to the plugin's stdin.
type PluginRequest struct {
	// The OpenAPI descriptions that were explicitly listed on the command line.
	// The specifications will appear in the order they are specified to openapic.
	Wrapper []*Wrapper `protobuf:"bytes,1,rep,name=wrapper" json:"wrapper,omitempty"`
	// The plugin parameter passed on the command-line.
	Parameter string `protobuf:"bytes,2,opt,name=parameter" json:"parameter,omitempty"`
	// The version number of openapi compiler.
	CompilerVersion *Version `protobuf:"bytes,3,opt,name=compiler_version,json=compilerVersion" json:"compiler_version,omitempty"`
}

func (m *PluginRequest) Reset()                    { *m = PluginRequest{} }
func (m *PluginRequest) String() string            { return proto.CompactTextString(m) }
func (*PluginRequest) ProtoMessage()               {}
func (*PluginRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *PluginRequest) GetWrapper() []*Wrapper {
	if m != nil {
		return m.Wrapper
	}
	return nil
}

func (m *PluginRequest) GetCompilerVersion() *Version {
	if m != nil {
		return m.CompilerVersion
	}
	return nil
}

// The plugin writes an encoded PluginResponse to stdout.
type PluginResponse struct {
	// Error message.  If non-empty, the plugin failed.
	// The plugin process should exit with status code zero
	// even if it reports an error in this way.
	//
	// This should be used to indicate errors which prevent the plugin from
	// operating as intended.  Errors which indicate a problem in openapic
	// itself -- such as the input Document being unparseable -- should be
	// reported by writing a message to stderr and exiting with a non-zero
	// status code.
	Error []string `protobuf:"bytes,1,rep,name=error" json:"error,omitempty"`
	// text output
	Text string `protobuf:"bytes,2,opt,name=text" json:"text,omitempty"`
}

func (m *PluginResponse) Reset()                    { *m = PluginResponse{} }
func (m *PluginResponse) String() string            { return proto.CompactTextString(m) }
func (*PluginResponse) ProtoMessage()               {}
func (*PluginResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

type Wrapper struct {
	// The filename or URL of the wrapped description
	Name string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	// The version of the OpenAPI specification that is used by the wrapped description.
	Version string `protobuf:"bytes,2,opt,name=version" json:"version,omitempty"`
	// Must be a valid serialized protocol buffer of the appropriate OpenAPI specification version.
	Value []byte `protobuf:"bytes,3,opt,name=value,proto3" json:"value,omitempty"`
}

func (m *Wrapper) Reset()                    { *m = Wrapper{} }
func (m *Wrapper) String() string            { return proto.CompactTextString(m) }
func (*Wrapper) ProtoMessage()               {}
func (*Wrapper) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func init() {
	proto.RegisterType((*Version)(nil), "openapic.v1.Version")
	proto.RegisterType((*PluginRequest)(nil), "openapic.v1.PluginRequest")
	proto.RegisterType((*PluginResponse)(nil), "openapic.v1.PluginResponse")
	proto.RegisterType((*Wrapper)(nil), "openapic.v1.Wrapper")
}

func init() { proto.RegisterFile("plugin.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 275 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x5c, 0x91, 0xbb, 0x4e, 0xec, 0x30,
	0x10, 0x86, 0x95, 0xb3, 0x97, 0x28, 0xb3, 0x7b, 0x00, 0x59, 0x2b, 0x94, 0x82, 0x02, 0xa5, 0xda,
	0x2a, 0x12, 0xd0, 0xd1, 0xf0, 0x04, 0x48, 0xc8, 0x05, 0x94, 0xc8, 0x44, 0xb3, 0x60, 0x94, 0xd8,
	0x83, 0xed, 0x84, 0x7d, 0x1a, 0x9e, 0x95, 0xf8, 0x26, 0x10, 0xdd, 0xfc, 0xdf, 0x4c, 0x32, 0xdf,
	0x24, 0xb0, 0xa5, 0x7e, 0x7c, 0x95, 0xaa, 0x25, 0xa3, 0x9d, 0x66, 0x1b, 0x4d, 0xa8, 0x04, 0xc9,
	0xae, 0x9d, 0xae, 0x9a, 0x0e, 0xca, 0x47, 0x34, 0x56, 0x6a, 0xc5, 0x76, 0xb0, 0x1a, 0xc4, 0xbb,
	0x36, 0x75, 0x71, 0x59, 0xec, 0x57, 0x3c, 0x86, 0x40, 0xa5, 0x9a, 0xe9, 0xbf, 0x44, 0x7d, 0xf0,
	0x94, 0x84, 0xeb, 0xde, 0xea, 0x45, 0xa4, 0x21, 0xb0, 0x73, 0x58, 0xdb, 0xf1, 0x70, 0x90, 0xc7,
	0x7a, 0x39, 0xe3, 0x8a, 0xa7, 0xd4, 0x7c, 0x15, 0xf0, 0xff, 0x21, 0x28, 0x70, 0xfc, 0x18, 0xd1,
	0x3a, 0xd6, 0x42, 0xf9, 0x69, 0x04, 0x11, 0xfa, 0x6d, 0x8b, 0xfd, 0xe6, 0x7a, 0xd7, 0xfe, 0xb2,
	0x6a, 0x9f, 0x62, 0x8f, 0xe7, 0x21, 0x76, 0x01, 0x15, 0x09, 0x23, 0x06, 0x74, 0x18, 0x4d, 0x2a,
	0xfe, 0x03, 0xd8, 0x1d, 0x9c, 0x75, 0x7a, 0x20, 0xd9, 0xa3, 0x79, 0x9e, 0xe2, 0x35, 0x41, 0xec,
	0xef, 0x6b, 0xd3, 0xa5, 0xfc, 0x34, 0x4f, 0x27, 0xd0, 0xdc, 0xc2, 0x49, 0xf6, 0xb3, 0xa4, 0x95,
	0x45, 0x7f, 0x20, 0x1a, 0xa3, 0xa3, 0x5e, 0xc5, 0x63, 0x60, 0x0c, 0x96, 0x0e, 0x8f, 0x2e, 0x19,
	0x84, 0xba, 0xb9, 0x87, 0x32, 0xe9, 0xfa, 0xb6, 0x9a, 0x95, 0xc2, 0x07, 0x9c, 0xdb, 0xbe, 0x66,
	0x35, 0x94, 0x59, 0x29, 0x3e, 0x95, 0xa3, 0x5f, 0x31, 0x89, 0x7e, 0xc4, 0xa0, 0xba, 0xe5, 0x31,
	0xbc, 0xac, 0xc3, 0x4f, 0xba, 0xf9, 0x0e, 0x00, 0x00, 0xff, 0xff, 0xcf, 0x68, 0xf0, 0x89, 0xb4,
	0x01, 0x00, 0x00,
}
