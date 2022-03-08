# The Modem Driver
The Modem Driver is the gateway to the Pyrosaurus game servers. It performs the following functions:
 * Manages the [Pyro User File](https://github.com/algae-disco/pyrosaurus-server/blob/main/Documentation/PYRO.USR-File-Specification.md)
 * Sends [Team Entry Files](https://github.com/algae-disco/pyrosaurus-server/blob/main/Documentation/Team-Entry-File-Spec.md)
 * Receives and manages [Contest Files](https://github.com/algae-disco/pyrosaurus-server/blob/main/Documentation/Contest%20File%20Format.md)
 * Sends and receives Messages to Admins (Send message to Evryware) and to other Players (Next/Previous opponent)

When a player first starts up the game they are presented with the option to test the Modem connectivity. This test process is described below in the [Modem Test Procedure](https://github.com/algae-disco/pyrosaurus-server/blob/main/Documentation/Modem%20Functionality.md#modem-test-procedure) section.

## Challenge Procedure
One of the first operations of the Modem Driver is to use a Challenge Procedure. The Modem Server must send the following bytes in sequence before any functionality is invoked:
* 0x32
* 0x3C
* 0x46

### Note
The challenge procedure may need a modification to a jump instruction in order for it to send the identity data.
More details on the mod will be placed here.

If this test is passed, then the Modem Driver sends Identity Information (see [table](https://github.com/algae-disco/pyrosaurus-server/blob/main/Documentation/Modem%20Functionality.md#identity-information-table) below) to the Modem Server.

The Server can respond with either of the following responses to this Identity Information.
  * 0x27 - means the identity data was ok
  * 0x63 - means the identity data was bad

### Identity Information Table
Here is a table of the Identity bytes that Modem Driver will send:
Field|Size (bytes)|Value|Description
---|---|---|---
Check byte 1|1|0x32| |
Check byte 2|1|0xCD|Check byte 1 + Check byte 2 should equate to 0xff (255)
Pyro String|6|PYROB0| |
Pyro User ID|4| |The value of the Pyro User ID from PYRO.USR
Pyro Check|2| |This is the check value from PYRO.USR
Version Number|1|0x2|This is a static value
Data size|2|0x8|Not sure why this is provided since this block appears to be fixed in size

Modem Driver expects the Server to inspect this information and respond with the appropriate value.

## Modes
The Modem Driver uses a Mode Validation procedure to tell the Modem Server what the Player is doing. 
The Driver will try to send out the current Mode five times before it will error out. 
The Modem Server must respond with the same byte to "accept" the Mode.

Mode Code|Description
---|---
0x1|Get Messages/Send Backup [User file](https://github.com/algae-disco/pyrosaurus-server/blob/main/Documentation/PYRO.USR-File-Specification.md)/Get Backup [User file ](https://github.com/algae-disco/pyrosaurus-server/blob/main/Documentation/PYRO.USR-File-Specification.md)
0x2|Send [Team Entry File](https://github.com/algae-disco/pyrosaurus-server/blob/main/Documentation/Team-Entry-File-Spec.md)
0x3|Retrieve [Contest File](https://github.com/algae-disco/pyrosaurus-server/blob/main/Documentation/Contest%20File%20Format.md)
0x4|Send Messages
0x7|[Test modem](https://github.com/algae-disco/pyrosaurus-server/blob/main/Documentation/Modem%20Functionality.md#modem-test-procedures)
0x21|This byte will tell the Modem Driver to loop up to 5 times; similar to a "wait" signal

### Receiving Team Entry Files
When Modem Driver sends data, it sends in 1024 byte blocks which start with 2 check bits and end with some indeterminate data.
Read those 1024 bytes, then send the "continue" response.
When less than 1024 bytes, then send "end" response.

### Sending Contests
stub

### Sending Messages etc
stub

### Receiving Messages
stub

## Modem Test Procedure
Much like the normal procedure, the [Challenge Procedure](https://github.com/algae-disco/pyrosaurus-server/blob/main/Documentation/Modem%20Functionality.md#challendge-procedure) and [Mode](https://github.com/algae-disco/pyrosaurus-server/blob/main/Documentation/Modem%20Functionality.md#modes) is set accordingly. Next an error checking procedure, and finally a "phone number" update procedure is invoked. 

### Error Checking Procedures
The Error Checking procedure will send all bytes between 0x00 and 0xFF inclusive.

Server must respond with the following byte to confirm no errors:
* 0x04

Anything else will result in a "Bad Response" error

Next, Modem Driver expects Server to send all bytes between 0x00 and 0xFF inclusive. The Modem Driver allows for up to 3 errors before it reports "Too many errors" and aborts the test.

If everything passed, the Modem Driver will report "Test successful".

### Phone Number Update Procedure

As a final step to the Test procedure the Modem Driver expects Server to send an updated "phone number" which it will save for the next connection. This appears to be a method for load balancing where the initial phone number shipped with the game will hopefully only be used to test the Modem connectivity. Once tested, then every subsequent call into the Evryware servers could be spread among a set of phone numbers. 

If a phone number is not sent, then the Modem Driver polling will time out and hang up the phone normally.

## Technical Background
Pyrosaurus ships with a Modem Driver as a secondary executable file called MODEM.EXE . This is a packed but uncompressed binary which is loaded to memory segment 0x47EC. The unpacking procedure simply obfuscates the entry point for the program. It is called by using DOS-style fork() where processing is entirely handed off to the MODEM program until the process completes, then processing falls back to the main game process. It is self contained so it does have its own procedures for loading fonts and displaying text to the screen.

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
