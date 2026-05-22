[![Sketch fonts](https://see.fontimg.com/api/rf5/BWWo5/YTZkZTMxNDlhNDEwNDZhZmFiZThhODFhNjA5N2U3NTgub3Rm/d2hvaWFt/typo-draft-demo.png?r=fs&h=250&w=2000&fg=5A3922&bg=FFFFFF&tb=1&s=125)](https://www.fontspace.com/typo-draft-font-f41179)

`whoiam` is a CLI tool that prevents fat-finger deployments to the wrong AWS account — because nothing ruins your morning quite like realising you just ran `terraform apply` against production.

<div align="center"><a href="https://imgflip.com/i/ashxye"><img src="https://i.imgflip.com/ashxye.jpg" title="made at imgflip.com"/></a></div>

You know that sinking feeling. You get dizzy and the walls start closing in on you. Was your session pointed to dev... or prod?
You know your team shouldn't have local production credentials, but hey... startups. We've all been there.

This has happened to me, and teams I have worked on more times than I care to admit. And that's why I built `whoiam`
A CLI tool that prevents accidental deployments to the wrong AWS account. Before running a command, it verifies that
your current credentials match the account you expect — protecting you from "fat finger" mistakes when working across multiple environments.

## Features

- Run commands scoped to a specific AWS account
- Retrieve AWS IAM Role information
- Supports multiple AWS accounts

## Documentation
For more information and usage examples, please refer to the [documentation](https://cliwright.github.io/whoiam/).

## Installation

### Using Homebrew

```sh
brew tap cliwright/homebrew-awstools
brew install whoiam
```

### Download Binary

You can download the pre-compiled binaries from the [releases page](https://github.com/cliwright/whoiam/releases).

To download the pre-compiled binaries from the releases page using `curl`, you can use the following commands. Replace `VERSION`, `OS`, and `ARCH` with the appropriate values for the version, operating system, and architecture you need.

For example, to download the Linux binary for version `v1.0.0`:

```sh
curl -L -o whoiam_linux_x86_64.tar.gz https://github.com/cliwright/whoiam/releases/download/v1.0.0/whoiam_Linux_x86_64.tar.gz
```

For the Windows binary:

```sh
curl -L -o whoiam_windows_x86_64.zip https://github.com/cliwright/whoiam/releases/download/v1.0.0/whoiam_Windows_x86_64.zip
```

For the macOS binary:

```sh
curl -L -o whoiam_darwin_x86_64.tar.gz https://github.com/cliwright/whoiam/releases/download/v1.0.0/whoiam_Darwin_x86_64.tar.gz
```

Make sure to replace `v1.0.0` with the actual version number you want to download.

### Build from Source

```sh
git clone https://github.com/cliwright/whoiam.git
cd whoiam
go build -o whoiam
```

## Usage

```sh
whoiam --help
```

## Initialisation
A config file can be generated at the default location `~/.whoiam/whoiam.yaml` by running the following command:

```sh
whoiam config init
```

## Configuration

`whoiam` uses the AWS SDK for Go, so it will look for credentials and configuration in the default locations used by the AWS CLI and SDKs.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request on GitHub.

## License

This project is licensed under the Apache License, Version 2.0. See the [LICENSE](LICENSE) file for details.

## Author

Jesse Maitland - [jesse@cliwright.com](mailto:jesse@cliwright.com)
```
