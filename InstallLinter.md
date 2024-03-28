# Installation

To install the golangci-lint binary on Windows run the following commands in **git bash**
(not in a cmd or PoweShell terminal):

```
# binary will be $(go env GOPATH)/bin/golangci-lint
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.57.1
```

See [golangci-ci install documentation](https://golangci-lint.run/welcome/install/) for more information about installation.

After installation of the golangci-lint binary, open an **elevated** PowerShell and run the following commands:

```
# This is the default installation path for git for Windows,
# you may need to change it to your own custom installation path 
$Env:Path += ";C:\Program Files\Git\usr\bin"
[Environment]::SetEnvironmentVariable("Path", $env:Path, [System.EnvironmentVariableTarget]::Machine)
```

The linter should now be installed, you can check the version of the linter by opening any terminal and running:

```
golangci-lint --version
```

# Usage
To run the linter, run the following command:

```
golangci-lint run
```

The linter will output the found issues in the CLI and make a golangci-lint-report.json file.
You can format this json using any (online) software you like.

To let the linter fix the found issues it is able to fix, open a terminal in the root of the InfoSec-Agent repository and run:

```
golangci-lint run --fix
```

This command will change the code files, make sure to inspect the changes before commiting/pushing them.
Not all found issues can be fixed automatically, you may need to fix some issues yourself.
