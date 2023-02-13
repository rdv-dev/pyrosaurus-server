# Build Instructions
So here are the build instructions for pyrosaurus-server! This project is entirely written in golang so you'll need to install the go environment and set it up.

These instructions are Linux/WSL specific.

## WSL GUI Dev Environment Setup
Run the following apt install commands:
```
apt install xfce4 xrdp net-tools chromium-browser git firefox golang ghex xfce4-power-manager gedit
```

Use the power manager to disable any screensaver.

For installing Sublime Text, run the below commands first.
```
apt-get install apt-transport-https
apt install ca-certificates
```

Next, follow [these instructions](https://www.sublimetext.com/docs/linux_repositories.html) from Sublime:

Finally, install Sublime:
```
apt install sublime-text
```

## Setup Go Environment
If you are unfamiliar, read up on setting up the GOPATH [here](https://go.dev/doc/gopath_code).

Run the following commands to setup the directory for pyrosaurus-server:
```
mkdir -p ~/go/bin
mkdir -p ~/go/src/github.com/algae-disco
```

Clone the repository:
```
git clone https://github.com/algae-disco/pyrosaurus-server
```

Build the package:
```
cd ~/go/src/github.com/algae-disco/pyrosaurus-server
go build
```

# Setup DOSBox

Clone the repository for [DOSBox-X](https://github.com/joncampbell123/dosbox-x) and follow [the build instructions](https://github.com/joncampbell123/dosbox-x/blob/master/BUILD.md).

For Heavy Debug, you may need to search in the src directory for any occurrance of C_HEAVY_DEBUG and comment out the pre-compile checks. Next, follow the build instructions for the regular debug build.

Once DOSBox is installed, run DOSBox.



# Setup Pyrosaurus

## First Run

