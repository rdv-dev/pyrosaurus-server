# PYRO.USR File Specification

## Background
The PYRO.USR file, hereon referenced as the User File, holds local information about the player, as well as a "phone number" to call the Pyrosaurus servers.
The game does not ship with a User File out of the box, however anytime the Game is run and it finds that the "PYRO.USR" file does not exist, it generates a basic User File with enough information to call the Pyrosaurus servers and gather critical information such as a User ID and new "phone number".

After the template User File is generated, the Game runs the Modem Driver. One of the Modem Driver's functions is sending and receiving the User FIle.
The first time this happens, the Pyrosaurus servers receive a file with a zeroed User ID field. 

This is a signal to the server to generate a new ID, perform any setup server side, and fill in these values:
 * Phone Number - after successful test only
 * ID Number
 * ID Check
 * Games Available
 * Arena
 * Rating
 * Last Opponent ID

## File Map
Field|Size|Position|Possible Values|Notes
---|---|---|---|---
Phone Number|12|0| |"Phone Number" means that for us with internet, this will need to be an IP address or domain name. DOSBox and other variants are able to handle this configuration. More information on this in the DOSBox Setup Instructions. This is updated after a successful Test.||
ID Number|4|C| |Maintained by Server||
ID check|2|10| |Maintained by ServerSee enhancement details below||
Games Available|1|12| |Maintained by Server||
Arena|1|13|A-0|Maintained by Server||
Rating|1|14|1-5, 1-25|Maintained by Server, 1-5 is relevant for levels 1 - 10, 1-25 is for level 11 (final level)||
Send/Retrieve Indicator|1|15|0 - Send Team, 1 - Retrieve Contests/Messages| ||
Unused data|3|16| |This isn't maintained by server and doesn't appear to be maintained by Pyro game||
Last opponent ID|4|19| |For sending messages to last opponent ||
Message/Call button state|1|1D|2, 4, 6|Displays or Hides the Call and Message buttons: 2 - Call only, 4 - Message only, 6 - Call & Message||
Last Date Contacted|4|1E| |Some unsigned int representing time?||

### Pyro User Check Enhancement
The Pyro User controls the players information for contests.
The second field Pyro User Check is never referenced or used by either the Game or the Modem driver except to write or read whatever value it contains.
The Modem Driver handles the data by placing whatever value it receives from the Server into this field.
We the Community have the opportunity to define what this data means.

The Pyro User Check field should be used as a key in addition to the Pyro User ID. 
If this field in a Player's PYRO.USR file is not what the Server is expecting, then the Server should fail the authentication process.
In addition to failing the authentication process, it may be valuable to enact a ban of the source IP if subsequent authentications fail as well.

The risk that comes with this process is that losing the PYRO.USR file would require a re-install of the game.
The Game's original Manual states this as well, so this was always a known issue.

Additionally, the algorithm used to determine the Pyro User Check has the following requirements:
 * Cannot be derived from known data; since we only have 2 bytes to store a hash, computation time is relatively small and the algorithm will be publicly available in this repository
 * Must reduce bias or be difficult to predict

Algorithm Design (for now):

 * The Pyro User Check field will store a pseudo-random number which changes based on a variable period stored by the Server.

Using an existing pseudo-random number generation process is ok because many other modern cryptographic libraries use this, so this is as good as we can get.
Changing this number periodically increases risk of an issue with the PYRO.USR file not authenticating, however it does increase security dramatically as this number is now no longer static.
Changing the period of when the pseudo-random number is changed increases security further as it is now unclear at which point-in-time the number could change.
However, it may be too risky to change this field every time a Player connects to the Server.
