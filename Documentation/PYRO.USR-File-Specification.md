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


Field|Size|Position|Possible Values|Notes
---|---|---|---|---
Phone Number|12|0| | ||
ID Number|4|C|Number in file + 7| ||
ID check|2|10|First two bytes of ID number?| ||
Games Available|1|12| | ||
Arena|1|13|A-0| ||
Rating|1|14|1-5, 1-25| ||
Send status indicator|1|15|0 - send team, 1 - Download contest| ||
Random data|3| | |Storing opponent data in a file like this is a clear security risk, so I'm guessing the field is random data which is meant to throw off anyone with a hex editor who might want to impersonate their opponent. This combined with the ID check earlier in the file would hopefully mitigate this.||
Next/Last opponent ID|4| |For sending messages to last opponent WTFFF| ||
Message/Call button state|1|1D|F-A|A - call only, retrieve contest & get messages B - call only, no special message C - message only D - message only E - msg and call, retrieve contest, send messages, gives option to resend team? F - message and call, doesn’t really work well||
Last Date Contacted|4| | |Some unsigned int representing time?||
