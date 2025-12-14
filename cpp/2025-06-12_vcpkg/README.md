# 2025-06-12_vcpkg

## Get started

Requirements: `cmake`

```bash
git submodule add --depth 1 https://github.com/microsoft/vcpkg.git
git submodule update --init
./vcpkg/bootstrap-vcpkg.sh

cmake -B build -S .
cmake --build build

./build/HelloWorld
```

## Log

Ref: [C++のためのvcpkgとCMake(Linux) #fmt - Qiita](https://qiita.com/king_dog_fun/items/bf2ff1fc961220d8c4bf)

```bash
sudo apt update
sudo apt install -y cmake

git submodule add --depth 1 https://github.com/microsoft/vcpkg.git
./vcpkg/bootstrap-vcpkg.sh

./vcpkg/vcpkg new --application
./vcpkg/vcpkg add port fmt

touch CMakeLists.txt main.cpp # and edit

cmake -B build -S .
cmake --build build

./build/HelloWorld # -> Hello World!
```
