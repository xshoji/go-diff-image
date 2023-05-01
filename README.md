# go-diff-image

A diff tool for images.



## Installation

```
$ go get github.com/xshoji/go-diff-image
```



## Usage

```bash
./diff-image -h
Usage of ./diff-image:
    -a, --after-image-path string
    	[required] Image path (after)
    -b, --before-image-path string
    	[required] Image path (before)
    -e, --error-if-difference
    	Be regarded as an error (status code 1) if difference exists
    -h, --help
    	Show help
    -n, --not-output-if-same
    	Not output difference image if inputs are same
    -o, --output string
    	[required] Output path of difference image

./diff-image -o=/tmp/d.png -b=/tmp/s1.png -a=/tmp/s2.png
```



## Example

![example](https://raw.githubusercontent.com/xshoji/go-diff-image/master/example.png)




## Build

```
# Build
APP=/tmp/diff-image; MAIN=cmd/diff-image/main.go; go build -ldflags="-s -w" -trimpath -o "${APP}" "${MAIN}"; chmod +x "${APP}"

# Croess compile by goreleaser
goreleaser --snapshot --skip-publish --rm-dist
```



## Release

Release flow of this repository is integrated with github action.
Git tag pushing triggers release job.

```
# Release
git tag v0.0.2 && git push --tags



# Delete tag
echo "v0.0.1" |xargs -I{} bash -c "git tag -d {} && git push origin :{}"

# Delete tag and recreate new tag and push
echo "v0.0.1" |xargs -I{} bash -c "git tag -d {} && git push origin :{}; git tag {} -m \"Release beta version.\"; git push --tags"

```
