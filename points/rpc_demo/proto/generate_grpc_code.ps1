$path = Get-Location

Set-Location $PSScriptRoot

    if (!(Test-Path "./impl/")) {
        mkdir "./impl"
    }

    # proto files:
    # out path is relative on proto file
    protoc --go_out=./impl --go_opt=paths=source_relative `
    --go-grpc_out=./impl --go-grpc_opt=paths=source_relative `
    user.proto

Set-Location $path

# windows not allow run ps script:
# (admin start)Set-ExecutionPolicy RemoteSigned
