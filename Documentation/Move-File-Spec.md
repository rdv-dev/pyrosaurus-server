# Move Data Specification
The Move data is one of the structures which hold Species training for moving around the arena. 
It is included in Team Entry Files when a Team is being sent to the Pyrosaurus servers, which encompasses the scope of this document.

The Move data has a limited size of 194 bytes (0xC2). All unused bytes in this block are marked by 0xFF .

## Source Codes
Here are the different types of Movements and their Source Codes
Source|Code
---|---
Self-Rotate|0||
Self-Fixed|2||
Self-Mobile|1||
Other-Rotate|10||
Other-Fixed|12||
Other-Mobile|11||
Map|20||
