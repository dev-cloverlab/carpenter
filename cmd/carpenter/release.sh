#!/bin/sh

XC_ARCH=${XC_ARCH:-386 amd64}
XC_OS=${XC_OS:-linux darwin windows}

rm -rf pkg/
rm -rf pkg-tmp/

gox \
    -os="${XC_OS}" \
    -arch="${XC_ARCH}" \
    -output "pkg-tmp/{{.OS}}_{{.Arch}}/{{.Dir}}"

mkdir pkg/
for file in `\find ./pkg-tmp -type f`; do
    FILE_NAME=${file##*/}
    FILE_PATH=${file%/*}
    `zip -jq ./pkg/${FILE_NAME%.*}_${FILE_PATH#*/*/}.zip ${file}`
done

ghr -u dev-cloverlab ${1:-v1.0.0} pkg/

rm -rf pkg/
rm -rf pkg-tmp/
