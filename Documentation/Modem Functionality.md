# Modem Functionality

Pyrosaurus ships with a Modem Driver as a secondary executable file called MODEM.EXE . This is a packed but uncompressed binary which is loaded to memory segment 0x47EC. The unpacking procedure simply obfuscates the entry point for the program. It is called by using DOS-style fork() where processing is entirely handed off to the MODEM program until the process completes, then processing falls back to the main game process. It is self contained so it does have its own procedures for loading fonts and displaying text to the screen.

When a player first starts up the game they are presented with the option to test the Modem functionality. This test process is described below in the [Modem Test Functionality](https://github.com/algae-disco/pyrosaurus-server/blob/main/Documentation/Modem%20Functionality.md#modem-test-procedures) section.

The Modem Driver is the gateway to the Pyrosaurus game servers. It performs the following functions:
 * Manages the [Pyro User File](https://github.com/algae-disco/pyrosaurus-server/blob/main/Documentation/PYRO.USR-File-Specification.md)
 * Sends [Team Entry Files](https://github.com/algae-disco/pyrosaurus-server/blob/main/Documentation/Team-Entry-File-Spec.md)
 * Receives and manages [Contest Files](https://github.com/algae-disco/pyrosaurus-server/blob/main/Documentation/Contest%20File%20Format.md)
 * Sends and receives Messages to Admins (Send message to Evryware) and to other Players (Next/Previous opponent)

## Modem Test Procedures
There is one challenge procedure, a mode validation procedure, then an error checking procedure, and finally a "phone number" update procedure.

### Challenge Procedure
The challenge procedure requires the following bytes be sent by the server in sequence:
* 0x32
* 0x3C
* 0x46
* At this point, the Modem Driver will send identity data to the server. Subsequently the server must respond with
  * 0x27 - means the identity data was ok
  * 0x63 - means the identity data was bad

### Modifications 
The challenge procedure may need a modification to a jump instruction in order for it to send the identity data.
More details on the mod will be placed here.

### Mode Validation

The mode validation procedure of the Modem Driver will try to read 1 byte five times before it will error out. 
This procedure will accept various input bytes which the server must provide. 
The Modem Driver knows what data it will accept from the server based on an argument passed to the procedure through the AX register.
* 0x7 - this byte means "test modem"
* 0x21 - regardless of input passed to this procedure, receiving this byte from the server will cause the procedure to loop without erroring up to 5 times
* Other bytes to be listed here in the future 

### Error Checking Procedures
Next comes the error checking procedures.
First, Modem Driver will send all bytes between 0x00 and 0xFF inclusive.

Server must respond with the following byte:
* 0x04

Anything else will result in a "Bad Response" error

Next, Modem Driver expects server to send all bytes between 0x00 and 0xFF inclusive. The Modem Driver allows for up to 3 errors before it reports "Too many errors" and aborts the test.

Otherwise, the Modem Driver will report "Test successful".

### Phone Number Update Procedure

Finally, the Modem Driver expects server to send an updated "phone number" which it will save for the next connection. This appears to be a method for load balancing where the initial phone number shipped with the game will hopefully only be used to test the modem functionality. Once tested, then every subsequent call into the Evryware servers could be spread among a set of phone numbers. 

If a phone number is not sent, then the Modem Driver polling will time out and hang up the phone normally.

## Community Modem Driver 
Since the Modem Driver is a secondary executable, it is possible to develop a community Modem Driver as well. 
A community Modem Driver would be executed in the same way by the Game - it doesn't care about what MODEM.EXE does. 
The primary requirement is MODEM.EXE must be built for a 16-bit Real-mode environment.

Use cases for this:
 * Overhaul communication protocol to increased security
 * Expand integration and capabilities with a modern server architecture
 * Advanced integration with VM

This is entirely optional and the original Modem Driver has enough features and advanced UART functionality. 
Some modern features can still be attained, like sending an Admin message for a link to a discord server.
