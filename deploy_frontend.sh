#!/bin/sh
rm -rf ./static
mkdir templates
tar -xvf dist.tar.xz
mv ./dist/static ./static
mv ./dist/index.html ./templates/index.html
rm -rf dist
rm -rf dist.tar.xz
sed -i 's/^/{{define "index"}}/' ./templates/index.html
sed -i 's/$/{{end}}/' ./templates/index.html