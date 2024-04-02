# About
pyrosaurus-server is a custom Contest server for DOS game Pyrosaurus developed by Evryware and released in 1996.

## Reverse Engineering
Most reverse engineering work was completed with [IDA 5.0](https://www.scummvm.org/news/20180331/) and so those files are
provided here for any future RE work. Many of the basic functions are named for ease of reading, though there are massive
functions in there which are total guesses on what is going on. See [this document](Documentation/RE-function-map.md) which highlights functions to review
for ease of reading. Download the PYRO ida file [here](Documentation/PYRO.idb) and the MODEM ida file [here](Documentation/MODEM-mod-20211105.idb)

For reversing, debugging and testing of the game I recommend using DOSBox-X as this fork has a superior debug interface. For heavy debugging I recommend downloading DOSBox-X source and commenting out all instances of C_HEAVY_DEBUG so these will build regardless of configuration.

## Build Instructions
See [build instructions here](Documentation/Build-instructions.md).

## Current Functionality
The pyrosaurus-server project imitates the original Evryware network, communicating directly with the game's modem utility program. A modified modem utility is [available here](Mods/MODEM.EXE) which allows sending the team and downloading the contest all at once. This greatly increaces the pacing of games.

What is working:
* Sending teams
* Receiving contests
* Reading team entries
* Compiling contest headers
* Partial arena (contest) emulation

## Documentation
To see what has been documented so far from reverse engineering research and some of the functionality relevant to this server, see the [Documentation here](https://github.com/algae-disco/pyrosaurus-server/blob/main/Documentation/README.md).

## Planned Features
* Enable Playback functionality for Contests
  * See technical documentation [here](https://github.com/algae-disco/pyrosaurus-server/blob/main/Documentation/Contest%20File%20Format.md)!
* Interface with Pyrosaurus Modem functionality - Complete!
  * [See this page for documentation](https://github.com/algae-disco/pyrosaurus-server/blob/main/Documentation/Modem%20Functionality.md).

## How Do I Get This Game?
You can download it directly from this repository [here](https://github.com/algae-disco/pyrosaurus-server/blob/main/Mods/pyro20.exe). Optionally you can download the game from archive.org . Note that the following link at [this location](https://web.archive.org/web/20010516030044/http://www.evryware.com/pyrosaurus/pyro20.exe) points directly to a EXE file captured on May 16, 2001. 

Which ever source you download from, use the hashes below to validate you have the correct file.

### SHA-256
```ba13558001701b881e90fde3a0387a4547d3e97f98096fc42acb81d4c4cfeeb6  pyro20.exe```
### MD5
```3cb84976f4d2ded210e7f49a1e0f2f5f  pyro20.exe```

The game runs via DOSBox/DOSBox-x or virtual machine guests DOS, Windows 3.x, or Windows 95/98 .
Run pyro20.exe and it will self-extract the game and all of its files.

```
      Copy PYRO20.EXE into its own directory on your hard drive (we
      suggest that you name the directory PYRO) and type PYRO20.  This
      creates all the files that you need to use Pyrosaurus
```

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
