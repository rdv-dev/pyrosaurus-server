# Species Data Specification
The Species data describes all information about the basic properties of a Dino from that Species. All training is defined at the Species level and is used by all Dinos of that Species. This specification's scope is limited to the data structure found in the Team Entry File.

The Species data is 32 bytes long. There is a discrepancy between the documentation here and what ends up in the Team Entry File. 

Field Name|Size (bytes)|Position|Description|Sample
---|---|---|---|---
Head Size|1|0|0-99d|||
Fire Range|1|1|0-99d|||
Fire Speed|1|2|0-99d|||
Fire Density|1|3|0-99d|||
Fire Pattern|1|4|0-99d|||
Fire Variation|1|5|0-99d|||
???|1|6|20-158d (14-9Eh)|Some number based on head size/predator/prey/etc
Leg Type + Straight/Sprawling|1|8|0=2 legs straight, 1=2 legs sprawl, 2=no legs, 3=4 legs straight, 4=4 legs sprawl| |
Leg Size|1|9|0-99d||||
Foot Type|1|A|0 - Hoof, 2 - Webbed, 1 - Claw|||
Foot Size|1|B|0-99d|||
Fire Risk|1|C|99-0h 0 high risk|||
Skin Armor|1|D|0 - Thin, 1 - Medium, 2 - Thick, 3 - Armor|||
Heart Size|1|E|0-99d|||
Tail Size|1|F|0-99d|||
Neck Size|1|10|0-99d|||
Predator/Prey|1|11|0 - Pred, 1 - Prey|||
sight range|1|12|0-99d|||
hearing range|1|13|0-99d|||
smell range|1|14|0-99d|||
sight field|1|15|0-99d|||
???|1|16|always 03?|||
???|1|18|E0, CE changes with pred/prey?? weird bug and resets?|||
Neck Speed|1|19|0-99d|||
Neck Variety|1|1A|0-99d|||
Fire Head/Body Target|1|1B|0-99d|||
Fire Head Movement|1|1C|0-99d|||
Leg code|1|1D|01 - no legs, 02 - legs|||
Leg code|1|1E|01 - legs, 02 - no legs|||
Fire Resolve|1|1F|0-99d|||

