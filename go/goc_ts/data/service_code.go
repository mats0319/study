package data

var ServiceCode_Template = `
import { axiosWrapper } from "./config"
import { {{ $structures }} } from "./{{ $filename }}.go"
{{ $serviceCode_Utils }}

class {{ $filenameBig }}Axios {{{ $serviceCode_Requests }}}

export const {{ $filename }}Axios: {{ $filenameBig }}Axios = new {{ $filenameBig }}Axios()
`

var ServiceCode_Request = "\n" +
	"{{ $indentation }}public {{ $serviceNameSmall }}({{ $paramsWithType }}): " +
	"Promise<{{ $serviceName }}Res> {{{ $serviceCode_ReqStruct }}\n" +
	"{{ $indentation }}{{ $indentation }}return axiosWrapper.post(\"{{ $serviceURI }}\"{{ $requestParams }})\n" +
	"{{ $indentation }}}\n"

var ServiceCode_ReqStruct = "\n" +
	"{{ $indentation }}{{ $indentation }}let req: {{ $serviceName }}Req = {\n" +
	"{{ $indentation }}{{ $indentation }}{{ $indentation }}{{ $paramsWithInit }}" +
	"{{ $indentation }}{{ $indentation }}}\n"
