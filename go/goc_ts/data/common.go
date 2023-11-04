package data

var Copyright = []byte("// Generate File, Should not Edit.\n" +
	"// Author: mario. https://github.com/mats9693\n" +
	"// Version: goc_ts " + Version + "\n")
var Version = "v0.2.0"

const RequestMessageSuffix = "Req"
const ResponseMessageSuffix = "Res"

const ServiceFileSuffix = ".http.ts"
const MessageFileSuffix = ".go.ts"

// MessageFieldType go type - ts type
var MessageFieldType = map[string]string{
	"string": "string",
	"int32":  "number",
}

// MessageFieldZeroValue ts type - zero value
var MessageFieldZeroValue = map[string]string{
	"string": `""`,
	"number": `0`,
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
