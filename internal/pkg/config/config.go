package config

import (
	"flag"
	"strings"
)

//Config represents config model
type Config struct {
	Iface          string
	Port           uint
	LogMetaHeaders string
	LoggerLevel    string
	ColorOutput    bool
	JSONOutput     bool
	ProtoPaths     string
	ProtoFiles     []string
	ProtoSetFiles   []string
}

var config *Config

var (
	iface          = flag.String("i", "eth0", "Interface to get packets from")
	port           = flag.Uint("p", 80, "TCP port")
	logMetaHeaders = flag.String("meta", "*", "Comma separated list of properties meta info for http2")
	loggerLevel    = flag.String("log-level", "info", "Logger level")
	colorOutput    = flag.Bool("color", false, "Output with color")
	jsonOutput     = flag.Bool("json", false, "Json output")
	protoPaths     = flag.String("proto-path", "", "Paths with proto files")
	protoFiles     = flag.String("proto-files", "", "Names of proto files")
	protoSetFile   = flag.String("proto-set", "", "Names of files containing encoded FileDescriptorSet.")
)

func splitByComma(str string) []string {
	if str == "" {
		return make([]string, 0)
	} else {
		return strings.Split(str, ",")
	}
}

//Init inits config
func Init() {
	flag.Parse()

	config = &Config{
		*iface,
		*port,
		*logMetaHeaders,
		*loggerLevel,
		*colorOutput,
		*jsonOutput,
		*protoPaths,
		splitByComma(*protoFiles),
		splitByComma(*protoSetFile),
	}
}

//GetConfig returns config
func GetConfig() *Config {
	return config
}

//GetLogMetaHeaders ...
func (config *Config) GetLogMetaHeaders() map[string]struct{} {
	result := make(map[string]struct{})

	logMetaHeaders := strings.TrimSpace(config.LogMetaHeaders)
	metaHeaders := strings.Split(logMetaHeaders, ",")

	for _, metaHeader := range metaHeaders {
		result[metaHeader] = struct{}{}
	}

	return result
}
