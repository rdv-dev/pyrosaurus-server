# Species Data Specification
The Species data describes all information about the basic properties of a Dino from that Species. All training is defined at the Species level and is used by all Dinos of that Species. This specification's scope is limited to the data structure found in the Team Entry File.

The Species data is 32 bytes long. The Team Entry file switches around some data between the Species Data and Dino data with the Leg Type - this needs to be more clearly documented.

Field Name|Size (bytes)|Position|Sample|Description
---|---|---|---|---
Head Size|1|0|0-99d|||
Fire Range|1|1|0-99d|||
Fire Speed|1|2|0-99d|||
Fire Density|1|3|0-99d|||
Fire Pattern|1|4|0-99d|Higher score here means more scatter across X axis||
Fire Variation|1|5|0-99d|Higher score here means more scatter across Y axis||
Error Checking 1|2|6|14h or 9Eh|Some derived number based on head size/predator/prey/etc|
Leg Type + Straight/Sprawling|1|8|0=2 legs straight, 1=2 legs sprawl, 2=no legs, 3=4 legs straight, 4=4 legs sprawl| |
Leg Size|1|9|0-99d||||
Foot Type|1|A|0 - Hoof, 2 - Webbed, 1 - Claw|||
Foot Size|1|B|0-99d|||
Fire Risk|1|C|99-0d, 0 = high risk|This is an invered scale, lower score here means higher risk tolerance during firing (chance of hitting friendly)||
Skin Armor|1|D|0 - Thin, 1 - Medium, 2 - Thick, 3 - Armor|||
Heart Size|1|E|0-99d|||
Tail Size|1|F|0-99d|||
Neck Size|1|10|0-99d|||
Predator/Prey|1|11|0 - Pred, 1 - Prey|Prey dinos have a 10% (? confirm this) bump in all sensing ranges: sight, hearing, smell. The selected score does not change but it results in a difference in-game. Prey has two sight cones whereas predator has one forward sight cone||
Sight Range|1|12|0-99d|||
Hearing Range|1|13|0-99d|||
Smell Range|1|14|0-99d|||
Sight Field|1|15|0-99d|||
Error Checking 2|2|16|always 03h|||
Error Checking 3|1|18|E0h or CEh|Some derived number changes with pred/prey||
Neck Speed|1|19|0-99d|||
Neck Variety|1|1A|0-99d|||
Fire Head/Body Target|1|1B|0-99d|Head/Body targeting chance. Lower score means more chance for body targeting, mid-range score means 50/50 chance between body/head, highest score always targets head||
Fire Head Movement|1|1C|0-99d|Score for head shake while firing. This is limited by neck length/mobility||
Leg code|1|1D|01 - no legs, 02 - legs|||
Leg code|1|1E|01 - legs, 02 - no legs|||
Fire Resolve|1|1F|0-99d|||

