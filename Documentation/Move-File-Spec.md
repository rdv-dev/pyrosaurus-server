# Move Data Specification
The Move data is one of the structures which hold Species training for moving around the arena. 
It is included in Team Entry Files when a Team is being sent to the Pyrosaurus servers, which encompasses the scope of this document.

The Move data has a limited size of 194 bytes (0xC2).

The Move data consists of a single data structure, repeated for how many Moves a Dino's Species has.

## Move Data Structure
The following data structure is repeated for every Move a Dino's Species has.
Field|Size (bytes)|Description
---|---|---
Number of points|1|Maxiumum of 10 points, referred to as n||
Source-Goal|1|See [Source-Goal Codes](https://github.com/algae-disco/pyrosaurus-server/blob/main/Documentation/Move-File-Spec.md#source-goal-codes) table||
X Coordinate 1|2|See note on [Coordinates](https://github.com/algae-disco/pyrosaurus-server/edit/main/Documentation/Move-File-Spec.md#source-goal-coordinate-information)||
Y Coordinate 1|2|See note on [Coordinates](https://github.com/algae-disco/pyrosaurus-server/edit/main/Documentation/Move-File-Spec.md#source-goal-coordinate-information)||
Goal Size 1|1| ||
...|...| ||
X Coordinate n|2| ||
Y Coordinate n|2| ||
Goal Size n|1| ||

## Source-Goal Codes
Here are the different types of Movements and their Source Codes
Source-Goal|Code|Description
---|---|---
Self-Rotate|0| ||
Self-Fixed|2| ||
Self-Mobile|1| ||
Other-Rotate|10| ||
Other-Fixed|12| ||
Other-Mobile|11|This is enabled with Flocking, defined in the Level Data||
Map|20| ||

## Source-Goal Coordinate Information
For all of the Source-Goals except Map, the Origin (0,0) is either the Dino's self, or some other Dino. All Movement points are relative to this Origin and ultimately need to be mapped to the Arena.

 * For all Rotate goals, the Movement points rotate with the Dino, as do the axies; the point does not move with the Source Dino
 * For all Fixed goals, the Movement points and axies points do not rotate or move
 * For all Mobile goals, the Movement point (there can only be one!) moves and rotates with the Source Dino; the point will never be reached (normally)

Put graphic here!

* The X axis is the horizontal center line of the dino. Positive X is right, negative X is left
* The Y axis is the vertical center line of the dino. Positive Y is behind, negative Y is ahead
* The Origin is at the very center of the Dino.

### Map Source-Goal Coordinates
The Level data's first two fields are Arena X & Y bounds expressed as a positive number.
This number, for example X / Y bounds of 2000 / 2000, is an absolute width of the level, however coordinates are mapped to an arena with 4 distinct quadrants.

Put graphic here!

* The X axis aligns along the horizontal boundary between the two teams. From the perspective of "home team" positive X is right, negative X is left
* The Y axis bifurcates the arena into equally spaced left and right areas. Positive Y is "home team" space, negative Y is "enemy team" space
* The Origin is at the very center of the Arena

This means no coordinate ever exceeds Arena X / 2 or Arena Y / 2 since each axis is, using our example from earlier, -1000 to 1000 coordinates, for an arena sized 2000 units for each axis.

The Community server doesn't have to worry about this too much because the Game handles the extreme boundaries. For example, a player cannot place a dino at the exact edge, the game has a boundary around the dino which disallows placing the dinos like that.
