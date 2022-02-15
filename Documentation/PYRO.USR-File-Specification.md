# PYRO.USR File Specification




Field|Size|Position|Possible Values|Notes
---|---|---|---|---
Phone Number|12|0| | ||
ID Number|4|C|Number in file + 7| ||
ID check|2|10|First two bytes of ID number?| ||
Games Available|1|12| | ||
Arena|1|13|A-0| ||
Rating|1|14|1-5, 1-25| ||
Send status indicator|1|15|0 - send team, 1 - Download contest| ||
Unknown|1| | | ||
Unknown|2| | | ||
Next/Last opponent ID|4| |For sending messages to last opponent WTFFF| ||
Message/Call button state|1|1D|F-A|A - call only, retrieve contest & get messages B - call only, no special message C - message only D - message only E - msg and call, retrieve contest, send messages, gives option to resend team? F - message and call, doesn’t really work well||
Last Date Contacted|4| | |Some unsigned int representing time?||
