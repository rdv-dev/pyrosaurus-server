# Levels

Level data files describe the properties and rules of an Arena. Game options can be configured here that affect the cost of dinos, 
the size of the arena, how many legs dinos can have, etc. These are fixed options where each byte corresponds to one option.

## Level Data Structure
Field Name|Size (bytes)|Position|Description
Arena Width X|2|0|||
Arena Width Y|2|2|||
Max Food Items|2|4|||
Max Team Score|2|6|||
Max Fire Range|2|8|||
Max Sight Range|2|10|||
Max Dinos on Team|2|12|||
Neck Cost|2|14|||
Sight Range Cost|2|16|||
Hearing Range Cost|2|18|||
Smell Cost|2|20|||
Leg Cost Multiplier|2|22|||
Leg/Foot Cost Multiplier|2|24|||
Leg Cost Multiplier?|2|26|||
Heart Size Cost/Multiplier|2|28|||
Base Endurance?|2|30|||
Base Neck Size Cost|2|32|||
Fire/Health?|2|34|||
Fire/Health?|2|36|||
Fire/Health?|2|38|||
Fire/Health?|2|40|||
Unknown|2|42|||
Unknown|2|44|||
Base Cost?|2|46|||
Unused/Zero value|2|48|||
Unused/Zero value|2|50|||
Unused/Zero value|2|52|||
Unused/Zero value|2|54|||
Max Time Limit (seconds)|2|56|||
Required Queens|2|58|||
Enable 4 Leg Dinos|2|60|||
Enable Pack Decision|2|62|||
Enable Increase Sight/Hearing?|2|64|||
Mini Map Size X|2|66|||
Mini Map Size Y|2|68|||
