# Installation

## Using Homebrew

To install `whoiam` using Homebrew, run the following commands:

```sh
brew tap pytoolbelt/awstools && \
brew install whoiam
```

## Download Binary

You can download the pre-compiled binaries from the [releases page](https://github.com/pytoolbelt/whoiam/releases).

To download the pre-compiled binaries using `curl`, use the following commands. Replace `VERSION`, `OS`, and `ARCH` with the appropriate values for the version, operating system, and architecture you need.

For example, to download the Linux binary for version `v1.0.0`:

```sh
curl -L -o whoiam_linux_x86_64.tar.gz https://github.com/pytoolbelt/whoiam/releases/download/v1.0.0/whoiam_Linux_x86_64.tar.gz && \
tar -xzf whoiam_linux_x86_64.tar.gz && \
chmod +x whoiam && \
sudo mv whoiam /usr/local/bin/
```

For the Windows binary:

```sh
curl -L -o whoiam_windows_x86_64.zip https://github.com/pytoolbelt/whoiam/releases/download/v1.0.0/whoiam_Windows_x86_64.zip
```

For the macOS binary:

```sh
curl -L -o whoiam_darwin_x86_64.tar.gz https://github.com/pytoolbelt/whoiam/releases/download/v1.0.0/whoiam_Darwin_x86_64.tar.gz && \
tar -xzf whoiam_darwin_x86_64.tar.gz && \
chmod +x whoiam && \
sudo mv whoiam /usr/local/bin/
```

Make sure to replace `v1.0.0` with the actual version number you want to download.

## Build from Source

To build `whoiam` from source, run the following commands:

```sh
git clone https://github.com/pytoolbelt/whoiam.git && \
cd whoiam && \
go build -o whoiam
```