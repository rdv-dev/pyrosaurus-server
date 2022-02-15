# PYRO.USR File Specification

## Background
The PYRO.USR file, hereon referenced as the user file, holds local information about the player, as well as a "phone number" to call the Pyrosaurus servers.
The game does not ship with a user file out of the box, however anytime the Game is run and it finds that the "PYRO.USR" file does not exist, it generates a basic user file with enough information to call the Pyrosaurus servers and gather critical information such as a User ID and new "phone number".

After the template user file is generated, the Game runs the Modem Driver. One of the Modem Driver's functions is sending and receiving the user file.
The first time this happens, the Pyrosaurus servers receive a file with a zeroed User ID field. 

This is a signal to the server to generate a new ID, perform any setup server side, and fill in these values:
 * Phone Number
 * ID Number
 * ID Check
 * Games Available
 * Send Status Indicator
 * Message/Call Button State
 * Last Date Contacted

## File Map
Field|Size|Position|Possible Values|Notes
---|---|---|---|---
Phone Number|12|0| |"Phone Number" means that for us with internet, this will need to be an IP address or domain name. DOSBox and other variants are able to handle this configuration. More information on this in the DOSBox Setup Instructions||
ID Number|4|C|Number in file + 7| ||
ID check|2|10| |First two bytes of ID number?||
Games Available|1|12| | ||
Arena|1|13|A-0| ||
Rating|1|14|1-5, 1-25|1-5 is relevant for levels 1 - 10, 1-25 is for level 11 (final level)||
Send status indicator|1|15|0 - Send Team, 1 - Download Contest| ||
Random data|3|16| |Storing opponent data in a file like this is a clear security risk, so I'm guessing the field is random data which is meant to throw off anyone with a hex editor who might want to impersonate their opponent. This combined with the ID check earlier in the file would hopefully mitigate this.||
Next/Last opponent ID|4|19| |For sending messages to last opponent WTFFF||
Message/Call button state|1|1D|F-A|A - call only, retrieve contest & get messages B - call only, no special message C - message only D - message only E - msg and call, retrieve contest, send messages, gives option to resend team? F - message and call, doesn’t really work well??||
Last Date Contacted|4|1E| |Some unsigned int representing time?||
