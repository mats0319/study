$path = Get-Location

Set-Location $PSScriptRoot
    if (!(Test-Path "./ui/dist")) {
        echo "please build web code first"
        return
    }

    if (Test-Path "./contract_manage") {
        Remove-Item "./contract_manage/*" -Recurse
    } else {
        mkdir "./contract_manage"
    }

    # binary executable file
    go build -o "contract_manage.exe"

    Move-Item "contract_manage.exe" -Destination "./contract_manage/contract_manage.exe"

    # manual
    Copy-Item "manual.txt" -Destination "./contract_manage/manual.txt"

    # web code
    mkdir "./contract_manage/ui"
    Copy-Item "./ui/dist/" -Destination "./contract_manage/ui/" -Recurse

Set-Location $path

# windows not allow run ps script:
# (admin start)Set-ExecutionPolicy RemoteSigned
