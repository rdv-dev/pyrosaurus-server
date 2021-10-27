import java.io.IOException;
import java.nio.ByteBuffer;
import java.nio.ByteOrder;
import java.io.FileOutputStream;
import java.util.Arrays;
//import java.nio.BufferOverflowException;

class PyroContests {

	static final short CONTEST_HEADER_RECORD_LEN = 16;

	static final short PYRO_USER_ID_LEN = 4;
	static final short TEAM_COLORS_LEN = 12;
	static final short TEAM_FIRE_COLORS_LEN = 6;
	static final short TEAM_COLORS_RECORD_LEN = PYRO_USER_ID_LEN 
							+ TEAM_COLORS_LEN 
							+ TEAM_FIRE_COLORS_LEN;

	static final short NUM_SPECIES_LEN = 1;
	static final short SPECIES_LEN = 0x20;

	static final short TOTAL_DINOS_LEN = 1;
	static final short TEAM_QUEEN_ARRAY_LEN = 1;
	static final short TEAM_SPECIES_LEG_NUM_LEN = 1;
	static final short TEAM_DINO_RESIZE = 3;
	static final short TEAM_X_POS_LEN = 2;
	static final short TEAM_Y_POS_LEN = 2;
	static final short TEAM_ROT_LEN = 2;
	static final short TEAM_SOME_DATA_CONT_LEN = 2;
	static final short ENTRY_DINO_RECORD_LEN = TEAM_QUEEN_ARRAY_LEN
									+ TEAM_SPECIES_LEG_NUM_LEN
									+ TEAM_X_POS_LEN
									+ TEAM_Y_POS_LEN
									+ TEAM_ROT_LEN
									+ TEAM_SOME_DATA_CONT_LEN;
	static final short CONTEST_DINO_RECORD_LEN = TEAM_QUEEN_ARRAY_LEN
									+ TEAM_SPECIES_LEG_NUM_LEN
									+ TEAM_DINO_RESIZE
									+ TEAM_X_POS_LEN
									+ TEAM_Y_POS_LEN
									+ TEAM_ROT_LEN
									+ TEAM_SOME_DATA_CONT_LEN;

	static final short LEVEL_DATA_SIZE = 70;

	public static void main(String[] args) throws IOException {
		int[] entry1Offsets;
		int[] entry2Offsets;
		int levelDataOffset;
		int contestDataOffset;
		int count;
		ByteBuffer contestPlayer1;
		ByteBuffer contestPlayer2;
		ByteBuffer contestData;
		boolean success1, success2;
		ContestEntry entry1 = new ContestEntry();
		ContestEntry entry2 = new ContestEntry();


		success1 = entry1.load(
			"C:\\Users\\rpdel\\OneDrive\\Projects\\Pyrosaurus Reversing\\Special files\\T-3dino-5deci.TMP.txt");

		success2 = entry2.load(
			"C:\\Users\\rpdel\\OneDrive\\Projects\\Pyrosaurus Reversing\\Special files\\T-2 dinos.TMP.txt");

		if (success1 && success2) {
			int entry1Size = TEAM_COLORS_RECORD_LEN + entry1.playerNamesSize
							+ NUM_SPECIES_LEN + (entry1.numSpecies * SPECIES_LEN)
							+ TOTAL_DINOS_LEN + (entry1.numDinos[0] * CONTEST_DINO_RECORD_LEN)
							+ entry1.dinoNamesSize;

			System.out.println("TEAM_COLORS_RECORD_LEN: " + TEAM_COLORS_RECORD_LEN);
			System.out.println("entry1.playerNamesSize: " + entry1.playerNamesSize);
			System.out.println("NUM_SPECIES_LEN: " + NUM_SPECIES_LEN);
			System.out.println("entry1.numSpecies * SPECIES_LEN: " + (entry1.numSpecies * SPECIES_LEN));
			System.out.println("TOTAL_DINOS_LEN: " + TOTAL_DINOS_LEN);
			System.out.println("entry1.numDinos[0] * CONTEST_DINO_RECORD_LEN: " + (entry1.numDinos[0] * CONTEST_DINO_RECORD_LEN));
			System.out.println("entry1.dinoNamesSize: " + entry1.dinoNamesSize);

			int entry2Size = TEAM_COLORS_RECORD_LEN + entry2.playerNamesSize
							+ NUM_SPECIES_LEN + (entry2.numSpecies * SPECIES_LEN)
							+ TOTAL_DINOS_LEN + (entry2.numDinos[0] * CONTEST_DINO_RECORD_LEN)
							+ entry2.dinoNamesSize;

			int contestEntrySize = CONTEST_HEADER_RECORD_LEN 
							+ entry1Size 
							+ entry2Size 
							+ LEVEL_DATA_SIZE;

			System.out.println("Allocating space for contest file: " 
							+ CONTEST_HEADER_RECORD_LEN + " + " 
							+ entry1Size + " + " + entry2Size
							+ " + " + LEVEL_DATA_SIZE
							+ " = " + contestEntrySize);

			
			contestPlayer1 = ByteBuffer.allocate(contestEntrySize);
			contestPlayer1.order(ByteOrder.LITTLE_ENDIAN);
			contestPlayer2 = ByteBuffer.allocate(contestEntrySize);
			contestPlayer2.order(ByteOrder.LITTLE_ENDIAN);
			contestData = ByteBuffer.allocate(1024);
			contestData.order(ByteOrder.LITTLE_ENDIAN);

			//*** Begin constructing Player1 Contest Header
			contestPlayer1.put((byte)0);
			contestPlayer1.putShort((short)(CONTEST_HEADER_RECORD_LEN + 1));
			contestPlayer1.putShort((short)(CONTEST_HEADER_RECORD_LEN + 1 
							+ TEAM_COLORS_RECORD_LEN 
							+ entry1.playerNamesSize));
			contestPlayer1.putShort((short)0);
			contestPlayer1.putShort((short)(CONTEST_HEADER_RECORD_LEN + 1 + entry1Size));
			contestPlayer1.putShort((short)(CONTEST_HEADER_RECORD_LEN + 1 
							+ entry1Size
							+ TEAM_COLORS_RECORD_LEN 
							+ entry2.playerNamesSize));
			contestPlayer1.putShort((short)0);
			contestPlayer1.putShort((short)(CONTEST_HEADER_RECORD_LEN + 1 + entry1Size + entry2Size));
			contestPlayer1.putShort((short)(contestEntrySize + 1));

			//*** Begin constructing Player1 file Home Team Strings
			contestPlayer1.putInt(entry1.pyroUserId);
			contestPlayer1.put(entry1.dinoColors);

			for(count=0;count<entry1.TEAM_NAMES_COUNT;count++) {
				contestPlayer1.put(entry1.playerData[count].getBytes());
				contestPlayer1.put((byte)0);
			}

			contestPlayer1.put((byte)entry1.numSpecies);

			for(count=0; count<entry1.numSpecies;count++) {
				contestPlayer1.put(entry1.speciesData[count]);
			}

			contestPlayer1.put(entry1.numDinos);

			// write queen data and species leg type
			contestPlayer1.put(Arrays.copyOfRange(
				entry1.dinoData,0,entry1.numDinos[0] + entry1.numDinos[0]));

			// write the weird extra 3 bytes for dino data
			for(count=0;count<entry1.numDinos[0];count++) {
				contestPlayer1.put((byte)0);
				contestPlayer1.put((byte)entry1.speciesData[0][9]);
				contestPlayer1.put((byte)0);
			}

			// write the rest of the dino data
			contestPlayer1.put(Arrays.copyOfRange(
				entry1.dinoData
				,entry1.numDinos[0] + entry1.numDinos[0]
				,entry1.numDinos[0] * ENTRY_DINO_RECORD_LEN));

			contestPlayer1.put(entry1.dinoNames);

			contestPlayer1.putInt(0);
			contestPlayer1.put(entry2.dinoColors);

			for(count=0;count<entry2.TEAM_NAMES_COUNT;count++) {
				contestPlayer1.put(entry2.playerData[count].getBytes());
				contestPlayer1.put((byte)0);
			}

			contestPlayer1.put((byte)entry2.numSpecies);

			for(count=0; count<entry2.numSpecies;count++) {
					contestPlayer1.put(entry2.speciesData[count]);
			}

			contestPlayer1.put(entry2.numDinos);

			// write queen data and species leg type
			contestPlayer1.put(Arrays.copyOfRange(
				entry2.dinoData,0,entry2.numDinos[0] + entry2.numDinos[0]));

			// write the weird extra 3 bytes for dino data
			for(count=0;count<entry2.numDinos[0];count++) {
				contestPlayer1.put((byte)0);
				contestPlayer1.put((byte)entry2.speciesData[0][9]);
				contestPlayer1.put((byte)0);
			}
			// write the rest of the dino data
			contestPlayer1.put(Arrays.copyOfRange(
				entry2.dinoData
				,entry2.numDinos[0] + entry2.numDinos[0]
				,entry2.numDinos[0] * ENTRY_DINO_RECORD_LEN));

			contestPlayer1.put(entry2.dinoNames);

			FileOutputStream of = new FileOutputStream("CONT.000",false);
			of.write(contestPlayer1.array());
			of.write(contestData.array());
			of.close();



			/*contestPlayer2.put((byte)0);
			contestPlayer2.putShort((short)(CONTEST_HEADER_RECORD_LEN + 1));
			contestPlayer2.putShort((short)(CONTEST_HEADER_RECORD_LEN + 1 
							+ TEAM_COLORS_RECORD_LEN 
							+ entry2.playerNamesSize));
			contestPlayer2.putShort((short)0);
			contestPlayer2.putShort((short)(CONTEST_HEADER_RECORD_LEN + 1 + entry2Size));
			contestPlayer2.putShort((short)(CONTEST_HEADER_RECORD_LEN + 1 
							+ entry2Size
							+ TEAM_COLORS_RECORD_LEN 
							+ entry1.playerNamesSize));
			contestPlayer2.putShort((short)0);
			contestPlayer2.putShort((short)(CONTEST_HEADER_RECORD_LEN + entry1Size + entry2Size));
			contestPlayer2.putShort((short)(contestEntrySize + 1));

			contestPlayer2.putInt(0);
			contestPlayer2.put(entry1.dinoColors);
			
			for(count=0;count<entry2.TEAM_NAMES_COUNT;count++) {
				contestPlayer2.put(entry1.playerData[count].getBytes());
				contestPlayer2.put((byte)0);
			}

			contestPlayer2.put((byte)entry1.numSpecies);

			for(count=0; count<entry1.numSpecies;count++) {
				contestPlayer2.put(entry1.speciesData[count]);
			}

			contestPlayer2.put(entry1.numDinos);

			contestPlayer2.put(Arrays.copyOfRange(entry1.dinoData,0,
				entry1.numDinos[0] * ENTRY_DINO_RECORD_LEN));

			contestPlayer2.put(entry1.dinoNames);

			contestPlayer2.putInt(entry2.pyroUserId);
			contestPlayer2.put(entry2.dinoColors);

			for(count=0;count<entry2.TEAM_NAMES_COUNT;count++) {
				contestPlayer2.put(entry2.playerData[count].getBytes());
				contestPlayer2.put((byte)0);
			}

			contestPlayer2.put((byte)entry2.numSpecies);

			for(count=0; count<entry2.numSpecies;count++) {
				contestPlayer2.put(entry2.speciesData[count]);
			}

			contestPlayer2.put(entry2.numDinos);

			contestPlayer2.put(Arrays.copyOfRange(entry2.dinoData,0,
				entry2.numDinos[0] * ENTRY_DINO_RECORD_LEN));

			contestPlayer2.put(entry2.dinoNames);*/
		}

	}
}
