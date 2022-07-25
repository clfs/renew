# renew
Update go install-ed binaries.

`renew` only depends on the Go standard library.

## Installation
Install `renew`:

```text
go install github.com/clfs/renew@latest
```

Uninstall `renew`:
```text
$ rm $(which renew)
```

## Usage
Run `renew` for help:

```text
$ renew
Usage of renew:
  -list
        list go install-ed binaries
  -skip
        skip binaries that fail to update
  -update string
        update a single binary to latest
  -update-all
        update all binaries to latest
```

Print a tab-delimited list of go install-ed binaries:

```text
$ renew -list
renew	github.com/clfs/renew
staticcheck	honnef.co/go/tools/cmd/staticcheck
tfsec	github.com/aquasecurity/tfsec/cmd/tfsec
```

Update a go-installed binary to `@latest`:

```text
$ renew -update staticcheck
==== staticcheck ====
[+] updated!
```

Update all go-installed binaries to `@latest`. Optionally, add `-skip` to
continue updating other binaries even if some updates fail.

```text
$ renew -update-all
==== renew
[+] updated!
==== staticcheck
[+] updated!
==== tfsec
[+] updated!
```