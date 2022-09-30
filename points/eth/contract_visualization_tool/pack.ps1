$path = Get-Location

Set-Location $PSScriptRoot
    if (!(Test-Path "./ui/dist")) {
        echo "please build web code first"
        return
    }

    if (Test-Path "./contract_visualization_tool") {
        Remove-Item "./contract_visualization_tool/*" -Recurse
    } else {
        mkdir "./contract_visualization_tool"
    }

    # binary executable file
    go build -o "contract_visualization_tool.exe"

    Move-Item "contract_visualization_tool.exe" -Destination "./contract_visualization_tool/contract_visualization_tool.exe"

    # manual
    Copy-Item "manual.md" -Destination "./contract_visualization_tool/manual.md"

    # web code
    mkdir "./contract_visualization_tool/ui"
    Copy-Item "./ui/dist/" -Destination "./contract_visualization_tool/ui/" -Recurse

Set-Location $path

# windows not allow run ps script:
# (admin start)Set-ExecutionPolicy RemoteSigned
