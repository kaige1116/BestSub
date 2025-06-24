#!/bin/bash

# Exit on error
set -e

# 设置时区
export TZ='Asia/Shanghai'

# 设置主函数目录
MAIN_DIR="./cmd/bestsub"

# 设置输出目录
OUTPUT_DIR="dist"

# App information
APP_NAME="bestsub"
BUILT_DATA="$(TZ='Asia/Shanghai' date +'%F %T %z')"
GIT_AUTHOR="bestrui"
GIT_VERSION=$(git describe --tags --abbrev=0)
COMMIT_ID=$(git rev-parse --short HEAD)

# Build flags
LDFLAGS="-X 'main.Version=${GIT_VERSION}' \
         -X 'main.BuildTime=${BUILT_DATA}' \
         -X 'main.Author=${GIT_AUTHOR}' \
         -X 'main.Commit=${COMMIT_ID}' \
         -s -w"

# Build tools
MUSL_BASE="https://musl.cc/"
ANDROID_NDK_BASE="https://dl.google.com/android/repository/"
ANDROID_NDK_VERSION="r27c"
TOOLCHAIN_DIR="$HOME/.bestsub/toolchains"


# Prepare build environment
prepare_build() {
    echo "Preparing build environment..."
    mkdir -p "${OUTPUT_DIR}"

    # install swag
    go install github.com/swaggo/swag/cmd/swag@latest
    # make swagger
    swag init   -g ${MAIN_DIR}/main.go -o ./api

    # tidy go mod
    go mod tidy
}


setup_android_toolchain() {
    if [ ! -d "$TOOLCHAIN_DIR/android-ndk" ]; then
        echo "Downloading Android NDK ${ANDROID_NDK_VERSION}..."
        curl -L -o /tmp/android-ndk-${ANDROID_NDK_VERSION}.zip "${ANDROID_NDK_BASE}android-ndk-${ANDROID_NDK_VERSION}-linux.zip" > /dev/null 2>&1
        mkdir -p "$TOOLCHAIN_DIR/android-ndk"
        echo "Extracting android-ndk-${ANDROID_NDK_VERSION}.zip..."
        sudo unzip -q /tmp/android-ndk-${ANDROID_NDK_VERSION}.zip -d "$TOOLCHAIN_DIR/android-ndk"
        rm /tmp/android-ndk-${ANDROID_NDK_VERSION}.zip
        echo "Android NDK ${ANDROID_NDK_VERSION} setup completed"
    else
        echo "Android NDK ${ANDROID_NDK_VERSION} already exists"
    fi
}


# Build for standard Linux
build() {
    local os=$1
    local arch=$2


    case $arch in
        "amd64")
            cgo_arch=amd64
            ;;
        "arm64")
            cgo_arch=arm64
            ;;
        "x86")
            cgo_arch=386
            ;;
        "arm")
            cgo_arch=arm
            ;;
        *)
            echo "Unsupported architecture: $arch"
            return 1
            ;;
    esac

    echo "Building ${os} ${arch}..."
    GOOS=${os} GOARCH=${cgo_arch} CGO_ENABLED=0 \
        go build -o "${OUTPUT_DIR}/${APP_NAME}-${os}-${cgo_arch}" -ldflags="${LDFLAGS}" -tags=jsoniter ${MAIN_DIR}
    echo "${os} ${arch} build completed"
}


build_android() {
    bin_path="${TOOLCHAIN_DIR}/android-ndk/android-ndk-${ANDROID_NDK_VERSION}/toolchains/llvm/prebuilt/linux-x86_64/bin"
    local arch=$1
    local cgo_cc
    local cgo_arch

    case $arch in
        "amd64")
            cgo_cc="x86_64-linux-android24-clang"
            cgo_arch="amd64"
            ;;
        "arm64")
            cgo_cc="aarch64-linux-android24-clang"
            cgo_arch="arm64"
            ;;
        "arm-7")
            cgo_cc="armv7a-linux-androideabi24-clang"
            cgo_arch="arm"
            ;;
        "386")
            cgo_cc="i686-linux-android24-clang"
            cgo_arch="386"
            ;;
        *)
            echo "Unsupported architecture: $1"
            return 1
            ;;
    esac
    echo "Building android ${cgo_arch}..."
    GOOS=android GOARCH=${cgo_arch} CC=${bin_path}/${cgo_cc} CGO_ENABLED=1 \
        go build -o "${OUTPUT_DIR}/${APP_NAME}-android-${arch}" -ldflags="${LDFLAGS}" -tags=jsoniter ${MAIN_DIR}
    ${bin_path}/llvm-strip "${OUTPUT_DIR}/${APP_NAME}-android-${arch}"
    echo "Android ${arch} build completed"
}



# Compress built files
compress_files() {
    echo "Compressing built files..."
    cd "${OUTPUT_DIR}"
    
    cp ../README.md ../LICENSE ./
    
    echo "Compressing Linux and Darwin builds..."
    for file in ${APP_NAME}-linux-* ${APP_NAME}-darwin-* ${APP_NAME}-android-*; do
        if [ -f "$file" ]; then
            cp "$file" "${APP_NAME}"
            tar -czf "${file}.tar.gz" "${APP_NAME}" README.md LICENSE
            rm "${APP_NAME}"   
            rm "$file"
        fi
    done
    
    echo "Compressing Windows builds..."
    for file in ${APP_NAME}-windows-*; do
        if [ -f "$file" ]; then
            cp "$file" "${APP_NAME}.exe"
            zip "${file%.*}.zip" "${APP_NAME}.exe" README.md LICENSE
            rm "${APP_NAME}.exe"   
            rm "$file"
        fi
    done
    
    rm README.md LICENSE
    
    cd ..
}

generate_checksums() {
    echo "Generating MD5 checksums..."
    cd "${OUTPUT_DIR}"
    find . -type f -print0 | xargs -0 md5sum > md5.txt
    cat md5.txt
    cd ..
}

rename_files() {
    echo "Renaming built files..."
    cd "${OUTPUT_DIR}"
    
    local renames=(
        "amd64:x86_64"
        "386:x86"
        "arm-7:armv7"
    )
    
    local max_length=0
    for file in ${APP_NAME}-*; do
        [ ${#file} -gt $max_length ] && max_length=${#file}
    done
    
    for file in ${APP_NAME}-*; do
        if [ -f "$file" ]; then
            local new_name="$file"
            for rename in "${renames[@]}"; do
                local old_arch="${rename%%:*}"
                local new_arch="${rename#*:}"
                new_name="${new_name//-$old_arch/-$new_arch}"
            done
            
            if [ "$file" != "$new_name" ]; then
                printf "%-${max_length}s -> %s\n" "$file" "$new_name"
                mv "$file" "$new_name"
            fi
        fi
    done
    
    cd ..
}

copy_docker_bin() {
    echo "Copying binaries for Docker build..."
    mkdir -p docker/{glibc,musl}
    
    local platforms=(
        "amd64:linux/amd64"
        "386:linux/386"
        "arm-7:linux/arm/v7"
        "arm64:linux/arm64"
    )
    
    local max_length=0
    for platform in "${platforms[@]}"; do
        local arch="${platform%%:*}"
        local name_glibc="${APP_NAME}-linux-${arch}"
        local name_musl="${APP_NAME}-linux-musl-${arch}"
        
        [ ${#name_glibc} -gt $max_length ] && max_length=${#name_glibc}
        [ ${#name_musl} -gt $max_length ] && max_length=${#name_musl}
    done
    
    for platform in "${platforms[@]}"; do
        local arch="${platform%%:*}"
        local docker_platform="${platform#*:}"
        
        mkdir -p "docker/glibc/${docker_platform}"
        mkdir -p "docker/musl/${docker_platform}"
        
        local name_glibc="${APP_NAME}-linux-${arch}"
        if [ -f "${OUTPUT_DIR}/${name_glibc}" ]; then
            printf "%-${max_length}s -> docker/glibc/%s/%s\n" \
                "${name_glibc}" "${docker_platform}" "${APP_NAME}"
            cp "${OUTPUT_DIR}/${name_glibc}" "docker/glibc/${docker_platform}/${APP_NAME}"
        elif [ -f "${OUTPUT_DIR}/${name_glibc}.tar.gz" ]; then
            printf "%-${max_length}s -> docker/glibc/%s/%s (from tar)\n" \
                "${name_glibc}" "${docker_platform}" "${APP_NAME}"
            tar xzf "${OUTPUT_DIR}/${name_glibc}.tar.gz" -C "docker/glibc/${docker_platform}" "${APP_NAME}"
        else
            printf "%-${max_length}s -> not found\n" "${name_glibc}"
        fi
        
        local name_musl="${APP_NAME}-linux-musl-${arch}"
        if [ -f "${OUTPUT_DIR}/${name_musl}" ]; then
            printf "%-${max_length}s -> docker/musl/%s/%s\n" \
                "${name_musl}" "${docker_platform}" "${APP_NAME}"
            cp "${OUTPUT_DIR}/${name_musl}" "docker/musl/${docker_platform}/${APP_NAME}"
        elif [ -f "${OUTPUT_DIR}/${name_musl}.tar.gz" ]; then
            printf "%-${max_length}s -> docker/musl/%s/%s (from tar)\n" \
                "${name_musl}" "${docker_platform}" "${APP_NAME}"
            tar xzf "${OUTPUT_DIR}/${name_musl}.tar.gz" -C "docker/musl/${docker_platform}" "${APP_NAME}"
        else
            printf "%-${max_length}s -> not found\n" "${name_musl}"
        fi
    done
    
    echo "Docker binaries prepared successfully"
}


if [ "$1" == "release" ]; then

    prepare_build
    setup_android_toolchain  # Only needed for Android builds

    # Android builds (requires CGO)
    build_android amd64
    build_android 386
    build_android arm-7
    build_android arm64

    # Standard Linux builds (CGO disabled for pure Go)
    build linux amd64
    build linux arm64
    build linux arm-7
    build linux 386

    # Windows builds
    build windows amd64
    build windows arm64
    build windows arm-7
    build windows 386

    # macOS builds
    build darwin amd64
    build darwin arm64
    build darwin arm-7
    build darwin 386

    copy_docker_bin
    rename_files
    generate_checksums
    compress_files
fi
