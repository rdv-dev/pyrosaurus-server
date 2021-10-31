# How To Create Contest Entry

## 1. Modify PYRO.USR with Hex editor

Note all entered values must be little endian.

* Enter user ID number at position 0x0C within 4 bytes (32-bit Long)
* At position 0x12 enter a non zero number, FF will give you 255d games
* At position 0x1D enter hex number 0E (14d)

Note: If the game is running while you modify this file, you'll need to restart the game to see these changes.

On the main menu, press the 'd' key to see your changes.

## 2. Create Your Team
The first steps of creating a Contest Entry is creating your team. Start by building species, where basic properties like head size, heart size, eye location, smell, hearing and vision distance are defined. Once you have the species defined, then start the training. Training is tied to the species level, so the same training can be shared between separate Dino's of the same species. Here movements, decisions for movements and fighting training is defined. Next create the actual Dino's. Finally you will be able to add Dino's to your team.

## 3. Creating The Entry
* Select Modem from the Main Menu
* Select Call
* A screen should display giving you the option to send your team. Press Y for yes 
* Let the system calculate
* Navigate to your PYRO folder outside of DOSBox
* Find the file “T.TMP”, rename and copy it to another folder
