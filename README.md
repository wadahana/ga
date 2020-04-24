# ga

How to build :

````
$ mkdir ga
$ cd ga
$ export GA=$PWD

# build vcapkey generator
$ mkdir third_party
$ cd $GA/third_party
$ git clone https://github.com/wadahana/MEmuVCapKey.git
$ cd $GA/third_party/MEmuVCapKey
$ ./build.sh

# build libvpx
$ cd $GA/third_party/
$ mkdir libvpx
$ wget https://github.com/webmproject/libvpx/archive/v1.8.2/libvpx-1.8.2.tar.gz
$ tar xzvf libvpx-1.8.2.tar.gz
$ pushd libvpx-1.8.2
$ ./configure --prefix=../libvpx --target=x86_64-win64-gcc --enable-static --disable-multithread --disable-install-docs --disable-unit-tests
$ make 
$ make install 

$ export GOPROXY="https://goproxy.cn"
$ export GO111MODULE=on
$ export CGO_ENABLED=1
$ export GOOS=windows
$ export GOARCH=amd64
$ export CC=x86_64-w64-mingw32-gcc
$ export CXX=x86_64-w64-mingw32-g++
$ export PKG_CONFIG_PATH=$PWD/third_party/libvpx/lib/pkgconfig
$ pushd $GA/
$ go get github.com/wadahana/ga
$ make

````
