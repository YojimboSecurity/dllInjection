# To check the correctness of this file install psscriptanalyzer with the following
#   choco install psscriptanalyzer
# To check the file run the following
#   Invoke-ScriptAnalyzer -Path .\make.ps1

function Build {
    <#
    .SYNOPSIS
        Build DLLInjection executable
    .DESCRIPTION
        Builds the executable for DLLInjection
    .EXAMPLE
        PS C:\> Build

        This builds the executable
    .INPUTS
        BIN_NAME (String): The name of the executable
        GIT_COMMIT (String): The git commit hash
        GIT_DIRTY (String): If there are any unconnited changes
        BUILD_DATE (String): The date of this build
    .OUTPUTS
        None
    .NOTES
        This is used to build DLLInjection
    #>
    param (
        [String]
        $VERSION=$($($(Select-String -Path .\version\version.go -Pattern "const Version") -split " ",4)[3] -replace '"',''),
        [String]
        $BIN_NAME="dllInjection",
        [String]
        $GIT_COMMIT=$(git rev-parse --short HEAD),
        [String]
        $GIT_DIRTY="",
        [String]
        $BUILD_DATE=$(Get-Date -Format "yyyy.MM.dd-HH:mm")
    )
    begin{
        if ($(git status --porcelain)){
            $GIT_DIRTY="+CHANGES"
        }
    }
    process{
        Write-Output "Building: $BIN_NAME $VERSION"
        Write-Output "Git Commit: $GIT_COMMIT$GIT_DIRTY"
        Write-Output "Build Date: $BUILD_DATE"

        go build -ldflags "-X github.com/YojimboSecurity/dllInjection/version.GitCommit=$GIT_COMMIT$GIT_DIRTY -X github.com/YojimboSecurity/dllInjection/version.BuildDate=$BUILD_DATE" -o bin/$BIN_NAME.exe
    }
    end{}
}

function help{
    param()
    begin{}
    process{
        Write-Output "make.ps1 is used to help users with common tasks such as
building an executable, cleaning up executables, running tests, and
getting dependencies.

Usage:
    make.ps1 [command]

Available Commands:
    help        Prints this message
    build       Build execuable
    test        Run tests
    clean       Cleanup executables
    get-deps    Get dependencies

"
    }
    end{}
}

switch -Exact ($args[0]) {
    "build" {
        Build
    }
    "test" {
        Write-Output "Testing..."
        go test -v ./...
    }
    "clean" {
        Remove-Item .\bin\dllInjection.exe -ErrorAction SilentlyContinue
    }
    "get-deps" {
        Write-Output "Getting Dependencies..."
        go mod tidy
        go mod vendor
    }
    "help" {
        help
    }
    Default {
        if ($args){
            Write-Output "$args is not a valid argument"
            help
            exit 1
        }
        Build
    }
}