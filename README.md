# go-diff-image

A diff tool for images.


## Installation
```
$ go get -u github.com/xshoji/go-diff-image/...
```

## Usage
```bash
go run cmd/diff-image/main.go -h
Usage of /var/folders/_q/dpw924t12bj25568xfxcd2wm0000gn/T/go-build1976644086/b001/exe/main:
    -a, --after-image-path string
    	[required] Image path (after)
    -b, --before-image-path string
    	[required] Image path (before)
    -f, --failure-if-diff-exists
    	To be failure if diff exists.
    -h, --help
    	Show help
    -o, --output string
    	[required] Output path

go run cmd/diff-image/main.go -o=/tmp/d.png -b=/tmp/s1.png -a=/tmp/s2.png
```
## Example

![example](https://raw.githubusercontent.com/xshoji/go-diff-image/master/example.png)
