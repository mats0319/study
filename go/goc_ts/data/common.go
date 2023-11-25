package data

var Copyright = []byte("// Generate File, Should not Edit.\n" +
	"// Author: mario. https://github.com/mats9693\n" +
	"// Version: " + Version + "\n")
var Version = "goc_ts v0.2.1"

const RequestMessageSuffix = "Req"
const ResponseMessageSuffix = "Res"

const ServiceFileSuffix = ".http.ts"
const MessageFileSuffix = ".go.ts"

// MessageFieldType go type - ts type
var MessageFieldType = map[string]string{
	"string": "string",
	"int32":  "number",
	"bool": "boolean",
}

// MessageFieldZeroValue ts type - zero value
var MessageFieldZeroValue = map[string]string{
	"string": `""`,
	"number": `0`,
	"boolean": `false`,
}

var indentation = 4 // useless, default value in flag

func SetIndentation(number int) {
	indentation = number
}

func GetIndentation(level int) []byte {
	res := make([]byte, 0, indentation*level)
	for i := 0; i < indentation*level; i++ {
		res = append(res, ' ')
	}
	return res
}
