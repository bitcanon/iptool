# IP Tool

## Overview

IP Tool is a command-line utility designed to simplify common tasks related to IP and subnetting. It's written in Go and is available as a single executable file that can be downloaded and used on most operating systems, including Windows, macOS, and Linux.

![iptool-inspect-demo1](docs/img/iptool-inspect-demo1.gif)

IP Tool offers the following functionality:

- **Inspection**: IP Tool enables you to inspect IP addresses, providing you with information such as the network address, broadcast address, and the number of hosts in the subnet.
- **Subnetting**: It also allows you to perform subnetting operations, and can be a handy tool for network engineers and administrators.

## Usage

```bash
iptool [command]
```

## Available Commands

- `inspect`: Take a closer look at an IP address
- `subnet`: Subnetting tools for IP networks

## Flags

- `--config string` : Config file (default is $HOME/.iptool.yaml)
- `-h, --help` : Help for iptool
- `-v, --version` : Version for iptool

Use `iptool [command] --help` for more information about a specific command.

## Installation

Here's a short instruction on how to get started using the IP Tool application by downloading the executable from its GitHub releases page and placing the file in your PATH:

1. **Download the Executable**:

   Visit the IP Tool releases page on GitHub by following this link: [https://github.com/bitcanon/iptool/releases](https://github.com/bitcanon/iptool/releases).

2. **Choose Your Operating System**:

   On the releases page, you will find a list of available releases for various operating systems. Choose the appropriate release for your system. IP Tool supports most operating systems, so select the one that matches your setup, such as Windows, macOS, or Linux.

3. **Download and Extract the Executable**:

   When you click on the release version to download, you will typically receive a compressed archive file (e.g., a `.zip` or `.tgz` file) that contains the IP Tool executable. Download this compressed archive and extract the executable from it using your preferred archive utility.

4. **Place the Executable in Your PATH**:

   Once the executable is downloaded and extracted, move it to a directory that is included in your systems PATH. This ensures that you can run the `iptool` command from any location in your terminal.

5. **Verify Installation**:

   To verify that the installation was successful, open a terminal window and type:

   ```bash
   iptool --version
   ```
   If you see the version number of IP Tool displayed in the terminal, it means the tool is installed and ready to use.

With these steps, you should now have the IP Tool executable properly downloaded, extracted, and accessible from your terminal. Enjoy using IP Tool for your networking tasks!

## Getting Started

Let's explore some of the common use cases for IP Tool.

### Inspect IP Addresses

To inspect the details if an IP address, use the `inspect` command. For example:

```bash
iptool inspect 10.0.0.1 255.255.255.0
```

If the IP address entered is valid, you will see the following output:

```bash
Address Details:
 IPv4 address       : 10.0.0.1
 Network mask       : 255.255.255.0

Netmask Details:
 Network mask       : 255.255.255.0
 Network bits       : 24
 Wildcard mask      : 0.0.0.255

Network Details:
 CIDR notation      : 10.0.0.0/24 (256 addresses):
 Network address    : 10.0.0.0
 Broadcast address  : 10.0.0.255
 Usable hosts       : 10.0.0.1 - 10.0.0.254 (254 hosts)
 ```

![iptool-inspect-demo2](docs/img/iptool-inspect-demo2.gif)

For more details on the `inspect` command, please refer to [Inspect Command](https://github.com/bitcanon/iptool/wiki/Inspect-Command) documentation.

### Subnet Commands

IP Tool also provides a set of commands for subnetting operations. To see the list of available commands, type:

```bash
iptool subnet
```

#### Subnet List

You can display a simple subnetting cheat sheet using the `subnet list` command. For instance:

```bash
iptool subnet list
```
>The alias `iptool subnet ls` can also be used.

This will print the following output:
```bash
CIDR  Subnet Mask      Addresses   Wildcard Mask
--------------------------------------------------
 /32  255.255.255.255  1           0.0.0.0
 /31  255.255.255.254  2           0.0.0.1
 /30  255.255.255.252  4           0.0.0.3
 /29  255.255.255.248  8           0.0.0.7
 /28  255.255.255.240  16          0.0.0.15
 ...
```

![iptool-subnet-list-demo](docs/img/iptool-subnet-list-demo.gif)

For more details on the `subnet list` command, please refer to [Subnet Command](https://github.com/bitcanon/iptool/wiki/Subnet-Command) documentation.

## Configuration

You can customize IP Tool's behavior by using a configuration file. By default, the tool looks for a configuration file at `$HOME/.iptool.yaml`.

## License

IP Tool is open-source software licensed under the [MIT License](LICENSE).