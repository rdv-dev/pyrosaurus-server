# About
pyrosaurus-server is a custom Contest server for DOS game Pyrosaurus developed by Evryware in 1996.

## Current Functionality
The pyrosaurus-server project can consume two Contest Entry files and produce a Contest file stub. All this means is that when the Contest is loaded, it will only display the dinos on the screen, and no actual actions will be performed.

## Planned Features
* Enable Playback functionality for Contests
  * See technical documentation [here](https://github.com/algae-disco/pyrosaurus-server/blob/main/Documentation/Contest%20File%20Format.md)!
* Interface with Pyrosaurus Modem functionality
  * Progress is being made here! [See this page for documentation](https://github.com/algae-disco/pyrosaurus-server/blob/main/Documentation/Modem%20Functionality.md).

## How Do I Get This Game?
You can download the game from archive.org . Note that the following link at [this location](https://web.archive.org/web/20010516030044/http://www.evryware.com/pyrosaurus/pyro20.exe) points directly to a EXE file captured on May 16, 2001. Use the hashes below to validate you have the correct file.

### SHA-256
```ba13558001701b881e90fde3a0387a4547d3e97f98096fc42acb81d4c4cfeeb6  pyro20.exe```
### MD5
```3cb84976f4d2ded210e7f49a1e0f2f5f  pyro20.exe```

The game runs via DOSBox/DOSBox-x or virtual machine guests DOS, Windows 3.x, or Windows 95/98 .
Run pyro20.exe and it will self-extract the game and all of its files.

## How Do I Create a Test Contest Entry?
Follow [these steps](https://github.com/algae-disco/pyrosaurus-server/blob/main/Documentation/Create%20Contest%20Entry.md) to create a Contest Entry.

## What is Pyrosaurus?
From the manual:
```      
      Welcome to Pyrosaurus.

      Now that you have Pyrosaurus, you can create and train a team of
      fire breathing dinos.  Make 'em lean and make 'em mean because
      they are going up against some heavy competition with bad
      intentions.

      Playing Pyrosaurus is like having your own professional sports
      team (only these guys battle to the death):

          You're the Scout: Find the players with the most potential.
          You're the General Manager: Organize your team, hire or fire
                                      players, etc.
          You're the Coach: Train your players to create an efficient
                            and professional team.

      We are the League Commissioners. We match up your team with
      another player's team. If your team wins, you move up in the
      standings. The best teams compete in the playoffs, Tyrannicus
      Supremus! After watching each "game", you can change players,
      retrain your team, or anything else you want to do before they
      compete again.
```
