# Build Instructions
So here are the build instructions for pyrosaurus-server! This project is entirely written in golang so you'll need to install the go environment and set it up.

These instructions are Linux/WSL specific.

## WSL GUI Dev Environment Setup
Run the following apt install commands:
```
apt install xfce4 xrdp net-tools chromium-browser git firefox golang ghex xfce4-power-manager gedit
```

Use the power manager to disable any screensaver.

## Linux Dev Environment Setup
```
apt install net-tools git golang ghex gedit
```

## Setup Go Environment
If you are unfamiliar, read up on setting up the GOPATH [here](https://go.dev/doc/gopath_code).

Run the following commands to setup the directory for pyrosaurus-server:
```
mkdir -p ~/go/bin
mkdir -p ~/go/src/github.com/rdv-dev
cd ~/go/src/github.com/rdv-dev
```

Clone the repository:
```
git clone https://github.com/rdv-dev/pyrosaurus-server
```

Build the package:
```
cd ~/go/src/github.com/rdv-dev/pyrosaurus-server
go mod init
go mod tidy
go build
```

# Setup DOSBox
## Install

Clone the repository for [DOSBox-X](https://github.com/joncampbell123/dosbox-x) and follow [the build instructions](https://github.com/joncampbell123/dosbox-x/blob/master/BUILD.md).

For Heavy Debug - recommended for any reverse engineering or analysis - you may need to search in the src directory for any occurrance of C_HEAVY_DEBUG and comment out the pre-compile checks. Next, follow the build instructions for the regular debug build.

Once DOSBox is installed, run it.

## Install Pyrosaurus
[See these instructions](install-pyrosaurus.md) for installing Pyrosaurus.

In the Menu, select "Main" then "Configuration Tool".

Select CPU, set "cycles" to 30000. Anything between 20000 and 30000 is recommended.

Select OK.

Select Serial Ports.

In the "serial1" field enter `modem listenport 2323`.

Select OK.

Select AUTOEXEC.BAT .

An example of the Autoexec, where <pyrosaurus directory> is the location where pyrosaurus is installed, example ~/pyro-c/pyro-44022/ :
```
mount c: <pyrosaurus directory>
c:
pyro
```
Select OK.
  
Select "Save"

Enter the following directory:
```
~/.config/dosbox-x/dosbox-x-pyro1.conf
```
  
You can setup multiple Pyrosaurus installations using this method.

# Setup Pyrosaurus
  
As of the writing of this documentation, you need to download the modded MODEM file or make the modification yourself.
  
Make a backup of the original MODEM.EXE file.
 ```
 cp MODEM.EXE MODEM-bkup.EXE
 ```
  
Download the file [here](../Mods/MODEM.EXE).
  
Documentation on the modification [here](/Modem%20Functionality.md#challenge-procedure) and [here](Modem%20Functionality.md#mode-1-details).

## First Run
When you first run Pyrosaurus, you will be presented with the Sound Setup screen. The default settings should work but test them in case it doesn't. Adjust settings as needed.
 
Select DONE

Next you'll see the Game Settings screen. I usually set SWAMP and MIST settings to the lowest to get better performance.
  
Select DONE
  
Next you'll see the Modem Setup screen.
  
Set Initialization to:
```
ATNET0
```

Under Port, select COM1.
  
Under IRQ, select 4.
  
Under Baud Rate, select 57600.

Set Dialing Prefix: (Note: this may vary depending on where the server exists. This guide assumes the server is listening on the loopback (127.0.0.1).
```
127000
```

Select DONE
  
Configuration within the game is mostly done and you'll be brought to the title screen!
  
Quit Pyrosaurus.
  
Next, using a hex editor, open file PYRO.USR

The connect string is the first 12 characters of the file.

We will be entering the following instead. Enter the characters:
```
0000018888
```

Or the Hexadecimal version:
```
303030303138383838
```

This will fill the first 11 characters, and the 12th needs to be zeroed out, so set the last character to 0.
  
Here is the final result.
# Debugging

