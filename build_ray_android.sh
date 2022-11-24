#!/usr/bin/env bash

ANDROID_NDK_HOME=/Users/six519/Library/Android/sdk/ndk/18.1.5063045
PATH=${ANDROID_NDK_HOME}/toolchains/aarch64-linux-android-4.9/prebuilt/darwin-x86_64/bin:${PATH}
ANDROID_SYSROOT=${ANDROID_NDK_HOME}/platforms/android-28/arch-arm64

if [ -z $1 ]; then
        echo "Titik filename is required..."
        exit 0
else
	echo "Setting up..."
	mv main.go main.go.bk
	mv main_android main.go

	touch tcode.go
	echo -e "package main\n\nvar tcode string = \`" > tcode.go
	cat $1 >> tcode.go
	echo -e "\n\`" >> tcode.go

	echo "Building..."
	CC=aarch64-linux-android-gcc \
	CGO_CFLAGS="-D__ANDROID_API__=28 -I${ANDROID_NDK_HOME}/sysroot/usr/include -I${ANDROID_NDK_HOME}/sysroot/usr/include/aarch64-linux-android --sysroot=${ANDROID_SYSROOT}" \
	CGO_LDFLAGS="-L${ANDROID_NDK_HOME}/sysroot/usr/lib -L${ANDROID_NDK_HOME}/toolchains/aarch64-linux-android-4.9/prebuilt/darwin-x86_64/lib/gcc/aarch64-linux-android/4.9.x/ --sysroot=${ANDROID_SYSROOT}" \
	CGO_ENABLED=1 GOOS=android GOARCH=arm64 \
	go build -tags ray -buildmode=c-shared -ldflags="-s -w -extldflags=-Wl,-soname,libexample.so" -o=android_build/libs/arm64-v8a/libexample.so

	echo "Cleaning up..."
	rm tcode.go
	mv main.go main_android
	mv main.go.bk main.go
fi