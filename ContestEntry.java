import java.io.IOException;
import java.io.FileInputStream;
import java.nio.ByteBuffer;
import java.nio.ByteOrder;
import java.nio.channels.FileChannel;
import java.lang.String;
import java.util.Arrays;

class ContestEntry {
	static final short ENTRY_HEADER_LEN = 15;
	static final short MAX_DINOS_ON_TEAM = 10;
	static final short SPECIES_LEN = 0x20;
	static final short MOVE_DATA_LEN = 0x64 * 2;
	static final short FITE_DATA1_LEN = 0xFE;
	static final short FITE_DATA2_LEN = 0x32 * 2;
	static final short FITE_DATA3_LEN = 0x32 * 2;
	static final short DECISIONS_LEN = 0x17D;

	static final short NUM_DINOS_ON_TEAM_LEN = 1;
	static final short TEAM_QUEEN_ARRAY_LEN = 1;
	static final short TEAM_SPECIES_LEG_NUM_LEN = 1;
	static final short TEAM_X_POS_LEN = 2;
	static final short TEAM_Y_POS_LEN = 2;
	static final short TEAM_ROT_LEN = 2;
	static final short TEAM_SOME_DATA_CONT_LEN = 2;
	static final short DINO_INIT_DATA_LEN = 25;
	static final short TEAM_ENTRY_RECORD_LEN = TEAM_QUEEN_ARRAY_LEN
									+ TEAM_SPECIES_LEG_NUM_LEN
									+ TEAM_X_POS_LEN
									+ TEAM_Y_POS_LEN
									+ TEAM_ROT_LEN
									+ TEAM_SOME_DATA_CONT_LEN
									+ DINO_INIT_DATA_LEN;
	static final short TEAM_CONTEST_RECORD_LEN = TEAM_QUEEN_ARRAY_LEN
									+ TEAM_SPECIES_LEG_NUM_LEN
									+ TEAM_X_POS_LEN
									+ TEAM_Y_POS_LEN
									+ TEAM_ROT_LEN
									+ TEAM_SOME_DATA_CONT_LEN;

	static final short TEAM_COLORS_LEN = 12;
	static final short TEAM_FIRE_COLORS_LEN = 6;
	static final short TEAM_COLORS_RECORD_LEN = TEAM_COLORS_LEN + TEAM_FIRE_COLORS_LEN;

	static final short TEAM_NAMES_MAX_LEN = 22;
	static final short TEAM_NAMES_COUNT = 3;

	int pyroUserId;
	short numSpecies;
	byte[][] speciesData;
	byte[][] moveData; // max length 64 * 2 * numSpecies; need to convert byte to word
	byte[][] fiteData1, fiteData2, fiteData3;
	byte[][] decisions;
	byte[] numDinos;
	byte[] dinoData;
	byte[] dinoNames;
	byte[] dinoColors;
	String[] playerData;
	int playerNamesSize;
	int dinoNamesSize;

	public void ContestEntry() {
	
	}

	public boolean load(String contestEntryFile)  throws IOException {
		int levelCheckSum;
		short teamDataOffset, dinoNamesOffset, teamColorStrings;
		short speciesCount;

		byte[] header = new byte[ENTRY_HEADER_LEN];
		this.speciesData = new byte[MAX_DINOS_ON_TEAM][SPECIES_LEN];
		this.moveData = new byte[MAX_DINOS_ON_TEAM][MOVE_DATA_LEN];
		this.fiteData1 = new byte[MAX_DINOS_ON_TEAM][FITE_DATA1_LEN];
		this.fiteData2 = new byte[MAX_DINOS_ON_TEAM][FITE_DATA2_LEN];
		this.fiteData3 = new byte[MAX_DINOS_ON_TEAM][FITE_DATA3_LEN];
		this.decisions = new byte[MAX_DINOS_ON_TEAM][DECISIONS_LEN];
		this.numDinos = new byte[1];
		this.dinoColors = new byte[TEAM_COLORS_RECORD_LEN];
		byte[] m_playerData = new byte[TEAM_NAMES_MAX_LEN * TEAM_NAMES_COUNT];

		FileInputStream entry = new FileInputStream(contestEntryFile);

		entry.read(header);
		
		ByteBuffer headerBuffer = ByteBuffer.wrap(header);
		headerBuffer.order(ByteOrder.LITTLE_ENDIAN);

		pyroUserId = headerBuffer.getInt(0);
		levelCheckSum = headerBuffer.getInt(4);
		teamDataOffset = headerBuffer.getShort(8);
		dinoNamesOffset = headerBuffer.getShort(10);
		teamColorStrings = headerBuffer.getShort(12);
		numSpecies = (short)headerBuffer.get(14);

		for(speciesCount=0; speciesCount<numSpecies; speciesCount++) {
			entry.read(this.speciesData[speciesCount]);
		}

		for(speciesCount=0; speciesCount<numSpecies; speciesCount++) {
			entry.read(this.moveData[speciesCount]);
		}

		for(speciesCount=0; speciesCount<numSpecies; speciesCount++) {
			entry.read(this.fiteData1[speciesCount]);
		}

		for(speciesCount=0; speciesCount<numSpecies; speciesCount++) {
			entry.read(this.fiteData2[speciesCount]);
		}

		for(speciesCount=0; speciesCount<numSpecies; speciesCount++) {
			entry.read(this.fiteData3[speciesCount]);
		}

		for(speciesCount=0; speciesCount<numSpecies; speciesCount++) {
			entry.read(this.decisions[speciesCount]);	
		}

		//*** Start Read Team Dinos Data 

		if (entry.getChannel().position() != teamDataOffset) {
			System.out.println("CONTEST ENTRY VALIDATION FAILED - Team Data wrong position");
			return false;
		}

		// read 1 byte number of dinos on team
		entry.read(this.numDinos);

		dinoData = new byte[TEAM_ENTRY_RECORD_LEN * this.numDinos[0]];

		entry.read(dinoData);
		
		ByteBuffer dinoBuffer = ByteBuffer.wrap(dinoData);
		dinoBuffer.order(ByteOrder.LITTLE_ENDIAN);
		// End Read Team Dinos Data

		//*** Start Read Team Dino Names
		if (entry.getChannel().position() != dinoNamesOffset) {
			System.out.println("CONTEST ENTRY VALIDATION FAILED - Dino Names wrong position");
			return false;
		}

		this.dinoNamesSize = teamColorStrings - dinoNamesOffset;

		dinoNames = new byte[dinoNamesSize];

		entry.read(dinoNames);

		ByteBuffer dinoNamesBuffer = ByteBuffer.wrap(dinoNames);
		dinoNamesBuffer.order(ByteOrder.LITTLE_ENDIAN);

		// End Read Team Dino Names

		//*** Start Read Team Colors
		if (entry.getChannel().position() != teamColorStrings) {
			System.out.println("CONTEST ENTRY VALIDATION FAILED - Team Colors wrong position");
			return false;
		}

		entry.read(this.dinoColors);

		// End Read Team Colors

		//*** Start Read Player Info strings
		entry.read(m_playerData);

		playerData = new String(m_playerData).split("\0");

		// adding TEAM_NAMES_COUNT here to account for 3 null (0) bytes to terminate strings
		this.playerNamesSize = 
				Arrays.stream(playerData).mapToInt(String::length).sum() 
				+ TEAM_NAMES_COUNT;
		// End Read Player Info strings

		entry.close();

		System.out.println("Pyro User ID: " + pyroUserId);
		System.out.println("Level Checksum: " + levelCheckSum);
		System.out.println("Player Name Size: " + this.playerNamesSize);
		System.out.println("Dino Names Size: " + this.dinoNamesSize);
		
		System.out.println("Species Count: " + numSpecies);
		System.out.println("Dino Count: " + this.numDinos[0]);

		return true;
	}
}