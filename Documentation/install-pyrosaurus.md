# Install Pyrosaurus

This guide is for installing Pyrosaurus on any machine, physical hardware, virtual or emulated.

The binary to install Pyrosaurus is available here. Setup a directory on your local machine for your DOSBox C drive, download Pyrosaurus and place it in the directory you created for the DOSBox C drive. 


## DOSBox Installation

For simply running the game, [DOSBox] is recommended as there is more readily available information related specificially to DOSBox.

Download DOSBox, install it, and run it.

Mount your DOSBox C Drive, and change to the C drive.

```
mount C <Your DOSBox C Drive>
cd c:
```

You can follow the instructions from the Manual here:

```
Copy PYRO20.EXE into its own directory on your hard drive - we
suggest that you name the directory PYRO - and type PYRO20.  This
creates all the files that you need to use Pyrosaurus.
```

If everything is working correctly, then you should see messages about files "inflating".

For more esoteric information, see the Manual, Section 2.

### DOSBox Configuration 

Run the following commands in DOSBox to configure DOSBox:

```
config -set "cpu cycles=30000"
config -set "serial serial1=modem listenport 2323"
```

You can save this configuration to file by running:

```
config -writeconf
```

You can also specify an argument to writeconf to specify a configuration specific file.


## First Run

To play:

```
DOS:

1. At the DOS Prompt, go to the Pyrosaurus directory.
2. Type: pyro
```

When you first run Pyrosaurus, you will be presented with the Sound Setup screen. The default settings should work but test them in case it doesn't. Adjust settings as needed.
 
Select DONE

Next you'll see the Game Settings screen. I usually set SWAMP and MIST settings to the lowest to get better performance.
  
Select DONE
  
Next you'll see the Modem Setup screen. For now, change nothing and select DONE.
  



