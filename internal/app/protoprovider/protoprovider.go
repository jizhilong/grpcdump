package protoprovider

import (
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	descpb "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/jhump/protoreflect/desc"
	"io/ioutil"

	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/jhump/protoreflect/dynamic"
)

var protoPaths = make(map[string]*ProtoMethod)

//ProtoMethod represents method of grpc
type ProtoMethod struct {
	Request  *dynamic.Message
	Response *dynamic.Message
}

//GetProtoByPath return ProtoMethod by path
func GetProtoByPath(path string) (*ProtoMethod, bool) {
	if protoMethod, ok := protoPaths[path]; ok {
		return protoMethod, true
	}

	return nil, false
}

//Init ...
func Init(importPaths string, protoFiles []string, protoSetFiles []string) error {
	var fileDescriptors []*desc.FileDescriptor
	var err error
	if len(protoSetFiles) != 0 {
		if len(protoFiles) != 0 {
			return errors.New("cant use both -proto-set and -proto-files flags")
		}
		fileDescriptors, err = initFromProtoSet(protoSetFiles)
	} else if len(protoFiles) == 0 {
		return errors.New("flag -proto-set or -proto-files required")
	} else {
		fileDescriptors, err = initFromProtoFiles(importPaths, protoFiles)
	}

	if err != nil {
		return err
	}

	if len(fileDescriptors) < 1 {
		return errors.New("Not found proto messages")
	}

	for _, parsedFile := range fileDescriptors {
		for _, service := range parsedFile.GetServices() {
			for _, method := range service.GetMethods() {
				protoPaths["/"+method.GetService().GetFullyQualifiedName()+"/"+method.GetName()] = &ProtoMethod{
					Request:  dynamic.NewMessage(method.GetInputType()),
					Response: dynamic.NewMessage(method.GetOutputType()),
				}
			}
		}
	}

	return nil
}

func initFromProtoSet(fileNames []string) ([]*desc.FileDescriptor, error) {
	files := &descpb.FileDescriptorSet{}
	for _, fileName := range fileNames {
		b, err := ioutil.ReadFile(fileName)
		if err != nil {
			return nil, fmt.Errorf("could not load protoset file %q: %v", fileName, err)
		}
		var fs descpb.FileDescriptorSet
		err = proto.Unmarshal(b, &fs)
		if err != nil {
			return nil, fmt.Errorf("could not parse contents of protoset file %q: %v", fileName, err)
		}
		files.File = append(files.File, fs.File...)
	}

	fileDescriptors, err := desc.CreateFileDescriptorsFromSet(files)
	if err != nil {
		return nil, err
	}

	result := make([]*desc.FileDescriptor, 0)
	for _, value := range fileDescriptors{
		result = append(result, value)
	}
	return result, nil
}

func initFromProtoFiles(importPaths string, protoFiles []string) ([]*desc.FileDescriptor, error) {
	fileNames, err := protoparse.ResolveFilenames([]string{importPaths}, protoFiles...)
	if err != nil {
		return nil, err
	}
	p := protoparse.Parser{
		ImportPaths:           []string{importPaths},
		InferImportPaths:      len(importPaths) == 0,
		IncludeSourceCodeInfo: true,
	}
	parsedFiles, err := p.ParseFiles(fileNames...)

	return parsedFiles, err
}

