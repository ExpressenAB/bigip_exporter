build() {
  echo "Building"
  build_windows
  build_linux
  build_docker $1
}

package() {
  echo "Packaging"
  package_windows $1
  package_linux $1
}

clean() {
  echo "Cleaning"
  clean_windows
  clean_linux
}

build_windows() {
  echo "Building windows"
  GOOS=windows GOARCH=amd64 go build -mod=vendor .
}

package_windows() {
  echo "Packaging windows"
  version=$1
  mkdir windows
  mv bigip_exporter.exe windows/
  cp LICENSE windows/
  tar -zcvf bigip_exporter-$version.windows-amd64.tar.gz windows
}

clean_windows() {
  echo "Cleaning windows"
  rm -rf windows
}

build_linux() {
  echo "Building linux"
  GOOS=linux GOARCH=amd64 go build -mod=vendor .
}

package_linux() {
  echo "Packaging linux"
  version=$1
  mkdir linux
  mv bigip_exporter linux/
  cp LICENSE linux/
  tar -zcvf bigip_exporter-$version.linux-amd64.tar.gz linux
}

clean_linux() {
  echo "Cleaning linux"
  rm -rf linux
}


fmt() {
  find . -not -path "./vendor/*" -name "*.go" -exec go fmt {} \;
}

build_docker() {
  version=$1
  id=$(docker build . | awk 'END {print $3}')
  if [[ $? != 0 ]]; then
    echo "Docker build failed"
    exit 1
  fi
  echo $id
  docker tag $id expressenab/bigip_exporter:$version
  docker tag $id expressenab/bigip_exporter:latest
  docker push expressenab/bigip_exporter:$version
  docker push expressenab/bigip_exporter:latest
}

if [[ $# == 0 ]]; then
  echo "At least one argument is needed"
elif [[ $1 == "fmt" ]]; then
  fmt
else
  build $1
  package $1
  clean
fi
