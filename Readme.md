# PKG-DIFF

## Description
pkg-diff returns packages name, that need to be built, based on their dependencies that were change in the last commit.

This is helpful in the CD process when some of your code has changed and you need to know which package to build.

## Examples
Output packages that need to be built based on the changes of the latest commit.

In plain format. Prints only the changed packages' names, one per line:
```bash
pkg-diff  -pkg-path myPkgPath

// Output:
myservice/app/services/component-monitoring
```

In JSON format:
```bash
pkg-diff -pkg-path myPkgPath -output=json | jq

// Output:
{
  "packages": {
    "myservice/app/services/component-monitoring": {
      "files": [
        "/projects/myservice/app/services/component-monitoring/routes.go",
        "/projects/myservice/app/services/component-monitoring/main.go",
        "/projects/myservice/app/services/component-monitoring/helpers.go",
        "/projects/myservice/app/services/component-monitoring/handlers.go"
      ],
      "name": "myservice/app/services/component-monitoring",
      "dependencies": [
        "myservice/pkg/web",
        "myservice/pkg/web/checkgrp",
        "myservice/internal/component-monitoring",
        "myservice/pkg/component",
        "myservice/pkg/logger",
        "myservice/pkg/tools"
      ],
      "isMain": true
    }
  }
}
```

## For Developers

Makefile commands:
```
build the project:
    make build
```

```
run local:
    ./.dist/pkg-structure -pkg-path <path_to_golang_project>
```