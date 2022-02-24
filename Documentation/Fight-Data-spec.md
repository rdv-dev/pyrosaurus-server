# Fight Data Specification

Fight Data consists of three data structures
 * Array of Movement Points: Point Rototation/Heading, and Movement Type
 * Array of Movement Point X Coordinates
 * Array of Movement Point Y Coordinates

Each of the Coordinates must rotate around the enemy dino's origin and heading. See Movement Origin section for more information.

## Fight Dataset 1

### Header
Field|Size (bytes)|Description
---|---|---
Mirror Flag|2|0 - Not Mirrored, 1 - Mirrored. This flag mirrors the fight points across the Y axis, yet only storing one set of them; mirrored points are computed later. This also doubles the number of points possible for the training.
Number of Bytes|2|This is an absolute number of bytes contained in the Array of Movement Points
Movement Point Array|See table|See table

This dataset has a maximum length of 0xFA (250) bytes. Each Movement Point structure is 5 bytes long. 
This means there are a maximum of 50 distinct Movement Points available.

### Movement Points
Field|Size (bytes)|Description
---|---|---
Unknown|1| 
Normal Vector in Degrees|1| 
Current Heading Offset from Normal in Degrees|1| 
Movement Type|1|See table
New Heading in Degrees|1|

#### Movement Types
Code|Movement Type
---|---
1|Jump Forward
2|Jump Backward
3|Jump Left
4|Jump Right
5|Run

## Fight Dataset 2 and 3
Both datasets have a maximum length of 0x64 or 100 bytes since each Coordinate is expressed with 2 bytes, one for each of the possible 50 Movement Points.

Fight Dataset 2 is for X Coordinates. Dataset 3 is for Y Coordinates.

They both follow the same convention where index 0 maps to the Movement Point at index 0 .


