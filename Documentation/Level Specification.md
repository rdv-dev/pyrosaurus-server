# Levels

Level data files describe the properties and rules of each Arena and the Division at large. Game options can be configured here that affect the cost of dinos, 
the size of the arena, how many legs dinos can have, etc. The position of each option is fixed.

Pyrosaurus comes delivered with 11 Arenas or Levels which make up the default Division. As it is coded into the game, each 
Arena has 5 rankings, except for the last Arena which has 25 rankings. Once the highest ranking is reached for that Arena, the
player is advanced to the next Arena at the bottom ranking.

Advancing to the next ranking is handled by overwriting file LEVEL.000 with the next available LEVL.00X file where X is the level the player has reached, counting down from 10 to 0.

## Level Data Structure
Field Name|Size (bytes)|Position|Description
---|---|---|---
Arena Size X|2|0|The total Arena horizontal width||
Arena Size Y|2|2|The total Arena vertical height||
Max Food Items|2|4|Change to 99 causes graphical issues! tests ok up to 15, full crash at 100h.||
Max Team Score|2|6|The maximum number of points allowed for the team||
Max Fire Range|2|8|More testing needed||
Max Sight Range|2|10|More testing needed||
Max Dinos on Team|2|12|Always set to 10 as delivered with the game but could be configured for less, unclear if the game can support more||
Neck Cost|2|14|More testing needed||
Sight Range Cost|2|16|More testing needed||
Hearing Range Cost|2|18|More testing needed||
Smell Cost|2|20|More testing needed||
Leg Cost Multiplier|2|22|More testing needed||
Leg/Foot Cost Multiplier|2|24|More testing needed||
Leg Cost Multiplier?|2|26|More testing needed||
Heart Size Cost/Multiplier|2|28|More testing needed||
Base Endurance?|2|30|More testing needed||
Base Neck Size Cost|2|32|More testing needed||
Fire/Health?|2|34|More testing needed||
Fire/Health?|2|36|More testing needed||
Fire/Health?|2|38|More testing needed||
Fire/Health?|2|40|More testing needed||
Unknown|2|42|More testing needed||
Unknown|2|44|More testing needed||
Base Cost?|2|46|More testing needed||
Unused/Zero value|2|48|All arenas delivered with the game have this set to 0, no idea if it's used by the game||
Unused/Zero value|2|50|All arenas delivered with the game have this set to 0, no idea if it's used by the game||
Unused/Zero value|2|52|All arenas delivered with the game have this set to 0, no idea if it's used by the game||
Unused/Zero value|2|54|All arenas delivered with the game have this set to 0, no idea if it's used by the game||
Max Time Limit (seconds)|2|56|The maxiumum time length for the arena. Tested successfully up to 25 minutes||
Required Queens|2|58|The main lose condition of the game is if one of your queens die, you lose. Where arenas are configured for 0 required queens, then the team with the last standing dino wins||
Enable 4 Leg Dinos|2|60|Enables quadrupeds in the Species - Leg screen ||
Enable Pack Decision|2|62|Display the Pack decision in the Train - Decisions screen||
Enable Increase Sight/Hearing?|2|64|More testing needed||
Mini Map Size X|2|66|The width of the minimap in pixels||
Mini Map Size Y|2|68|The height of the minimap in pixels||

## Original Arena Descriptions
Arena Number|File Name|Arena Name|Description
---|---|---|---
11|LEVEL.000|Trio Bellum: "Beginner’s Battle"|New players start in the Trio Bellum arena. Trio Bellum means      "Beginner's Battle" but don't be deceived, you are in for some tough fights. This arena is small and you only have 3000 points to build your team with.  Use your points where they will do the most good and train your dinos to take advantage of the attributes that you give them.||
10|LEVL.009|Quadriennium: "Going on Four Feet"|Now you can make strong and fast Quadrupeds for your team! Do you have that penned in feeling? It is time to stretch out: the arena is over 5 times larger! With all that extra room, you will need extra points to build a stronger team to cover all that territory: you now have 4,000 points to build your team!||
9|LEVL.008|Augustus Harena: "Narrow Arena"|Food is a new addition to this arena. Small bushes are growing in the arena. When bushes are eaten, health and strength are restored. You also get an extra 1,000 points to spend on your team!||
8|LEVL.007|Tribus Bellatrix: "Three Female Warriors"|Welcome to the land of three Queens. You can now teach your team to run in packs and you have three Queens that need to be protected! The contests are longer, the arena is bigger and you have 6,000 points to spend!||
7|LEVL.006|Laxus Brevis: "Wide and Short"|Good offensive and good defensive skills are both critical. Danger is never very far away because there is nowhere to hide! You have 7,000 points to build a powerful team! This is truly survival of the fittest!||
6|LEVL.005|Regina Omnis: "All Queens"|This time, every member of your team is a Queen! The first team to eliminate an opponent wins. Use caution. The arena is big with plenty of opportunity to get into trouble. You will want to spend every one of your 8,000 points!||
5|LEVL.004|Magnus Contendo: "Big Struggle"|This is a return to the long and narrow arena you saw in Augustus Harena only now it is much longer and not quite so narrow. Your dinos have longer fire range: just the thing to clear a path to that enemy Queen. You also get more points and a longer time limit.||
4|LEVL.003|Nulla Regina: "No Queens"|Here is a twist, this arena has no Queens! The team with the last dino standing wins.||
3|LEVL.002|Vastus Campus: "A Vast Plain"|This is a huge arena with plenty of room to roam. Large hearts cost less points (to increase endurance) and the time limit is extended to 9 minutes.||
2|LEVL.001|Congregari Venator: "Pack Hunters"|Five Queens to a side, 12,000 points, a 9 1/2 minute time limit and vicious opponents. This is no place for cowards!||
1|LEVL.000|Tyrannicus Supremus: "Supreme Tyrant"|Does your team have what it takes to become the Supreme Tyrant? This is the last Arena. When you make it this far, you are among the best of the best. Two significant rules change now. The other Arenas had 5 ratings each, this one has 25! You now will go down a rating if you lose (but you can't drop down to a lower arena). This arena is not for the fainthearted! You have 13,000 points, a 10 minute limit and tough competition! If you reach the top level, your name and your team name will be added to The Official Pyrosaurus Home Page for the whole world to see!||

## Level Data Analysis
Below is an image of the analysis of Level Data. When a value for the arena increases, it is highlighted in green. When it decreases, it is highlighted in red.
![Level Data Analysis](Pyrosaurus%20Level%20Analysis.PNG?raw=true)
