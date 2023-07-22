<h1>pkg-structure</h1>
<p>pkg-structure will map the dependencies of the golang packages your project is using.</p>
<p>This is helpful in the CD process when some of your code has changed and you need to know which package has changed and it's dependencies to build the correct application</p>

```
build:
    make build
```

```
    run:
        ./.dist/pkg-structure -pkg-path <path_to_golang_project>
```