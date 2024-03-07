package Database

import (
	"database/sql"
	// "fmt"
	"log"
	"math/rand"
	"time"

	// "github.com/rdv-dev/pyrosaurus-server/ContestServer/util"

	_ "github.com/mattn/go-sqlite3"
)

const ( 
	dbFile = "Database/pyrosaurus-server.db"
)

var db *sql.DB

// Player represents a player in the game.
type Player struct {
	PlayerID     uint32  // Primary key (auto-incremented)
	CheckID      uint16    // Not null
	Arena        uint16    // Nullable
	Rating       uint16    // Nullable
	EmailAddr    string // Nullable
	MFAEnabled   int    // Nullable
	MFAMethod    int    // Nullable
	ContestEntry uint64  // Foreign key referencing CONTEST_ENTRY(ENTRY_ID)
	TeamName     string // Nullable
	Location     string // Nullable
	PlayerName   string // Nullable
}


// InitializeDatabase creates the necessary tables if they don't exist.
func InitializeDatabase() {

    var createErr error

	db, createErr = sql.Open("sqlite3", dbFile)
	if createErr != nil {
		log.Fatal("Error opening database:", createErr)
	}
	// defer db.Close()

	log.Print("Database opened")

	createPlayersTableSQL := `
		CREATE TABLE IF NOT EXISTS PLAYERS (
			PLAYER_ID INTEGER PRIMARY KEY AUTOINCREMENT,
			CHECK_ID INTEGER NOT NULL,
			ARENA INTEGER,
			RATING INTEGER,
			EMAIL_ADDR TEXT,
			MFA_ENABLED INTEGER,
			MFA_METHOD INTEGER,
			CONTEST_ENTRY INTEGER,
			TEAM_NAME TEXT,
			LOCATION TEXT,
			PLAYER_NAME TEXT,

			FOREIGN KEY (CONTEST_ENTRY) REFERENCES CONTEST_ENTRY(ENTRY_ID)
		);
	`

	createPlayerPyroIdTableSQL := `
		CREATE TABLE IF NOT EXISTS PLAYER_PYRO_ID (
			PYRO_ID INTEGER PRIMARY KEY,
			PLAYER_ID INTEGER,

			FOREIGN KEY (PLAYER_ID) REFERENCES PLAYERS(PLAYER_ID)
		);
	`

	createContestEntryTableSQL := `
		CREATE TABLE IF NOT EXISTS CONTEST_ENTRY (
			ENTRY_ID INTEGER PRIMARY KEY AUTOINCREMENT,
			UPLOADED DATETIME,
			ENTRY_DATA BLOB
		);
	`

	createPlayerEntryTableSQL := `
		CREATE TABLE IF NOT EXISTS PLAYER_ENTRY (
			PLAYER_ID INTEGER,
			ENTRY_ID INTEGER,

			FOREIGN KEY (PLAYER_ID) REFERENCES PLAYERS(PLAYER_ID),
			FOREIGN KEY (ENTRY_ID) REFERENCES CONTEST_ENTRY(ENTRY_ID)
		);
	`

	createContestTableSQL := `
		CREATE TABLE IF NOT EXISTS CONTEST_DATA (
			CONTEST_ID INTEGER PRIMARY KEY AUTOINCREMENT,
			ENTRY_ID_0 INTEGER NOT NULL,
			ENTRY_ID_1 INTEGER NOT NULL,
			RETRIEVED INTEGER, -- 0 or 1 (1 indicates retrieved)
			WINNER INTEGER, -- 0 or 1 corresponding to ENTRY_ID_0 or ENTRY_ID_1
			CONTEST_DATE DATETIME,
			CONT_DATA BLOB,

			FOREIGN KEY (ENTRY_ID_0) REFERENCES CONTEST_ENTRY(ENTRY_ID),
			FOREIGN KEY (ENTRY_ID_1) REFERENCES CONTEST_ENTRY(ENTRY_ID)

		);
	`

	createMessageTableSQL := `
		CREATE TABLE IF NOT EXISTS MESSAGE (
			MESSAGE_ID INTEGER PRIMARY KEY AUTOINCREMENT,
			PLAYER_ID_0 INTEGER NOT NULL, -- PLAYER_ID of the player receiving the message
			PLAYER_ID_1 INTEGER NOT NULL, -- PLAYER_ID of the player who sent the message
			RETRIEVED INTEGER, -- 0 or 1 (1 indicates retrieved)
			MESSAGE_DATA BLOB,

			FOREIGN KEY (PLAYER_ID_0) REFERENCES PLAYERS(PLAYER_ID),
			FOREIGN KEY (PLAYER_ID_1) REFERENCES PLAYERS(PLAYER_ID)
		);

	`

	// prepSQLCreatePlayer = "INSERT INTO PLAYERS (CHECK_ID) VALUES (?)"
	// prepSQLCreatePlayerPyroId = "INSERT INTO PLAYER_PYRO_ID (PYRO_ID, PLAYER_ID) VALUES (?, ?)"
	// prepSQLGetPlayer = "SELECT PLAYER_ID FROM PLAYER_PYRO_ID WHERE PYRO_ID = ?"


	_, err := db.Exec(createPlayersTableSQL)
	if err != nil {
		log.Fatal("Error creating PLAYERS table:", err)
	}

	_, err = db.Exec(createPlayerPyroIdTableSQL)
	if err != nil {
		log.Fatal("Error creating PLAYER_PYRO_ID table:", err)
	}

	_, err = db.Exec(createContestEntryTableSQL)
	if err != nil {
		log.Fatal("Error creating CONTEST_ENTRY table:", err)
	}

	_, err = db.Exec(createPlayerEntryTableSQL)
	if err != nil {
		log.Fatal("Error creating PLAYER_ENTRY table:", err)
	}

	_, err = db.Exec(createContestTableSQL)
	if err != nil {
		log.Fatal("Error creating CONTEST_DATA table:", err)
	}

	_, err = db.Exec(createMessageTableSQL)
	if err != nil {
		log.Fatal("Error creating MESSAGE table:", err)
	}
}

// CreatePlayerS adds a new players to the PLAYERS table and returns the generated ID.
func CreatePlayer() (uint32, error) {
	var newId uint32
	var newCheckId int
	rand.Seed(time.Now().UnixNano())
	
	newCheckId = rand.Intn(65534) + 1
	result, err := db.Exec("INSERT INTO PLAYERS (CHECK_ID) VALUES (?)", newCheckId)
	if err != nil {
		return 0, err
	}

	internalPlayerId, _ := result.LastInsertId()

	newId = uint32(rand.Intn(4294967270) + 1) // Not full 4,294,967,295 because game adds 7 to this number when displaying
	result, err = db.Exec("INSERT INTO PLAYER_PYRO_ID (PYRO_ID, PLAYER_ID) VALUES (?, ?)", newId, internalPlayerId)

	for err != nil {
		newId = uint32(rand.Intn(4294967270) + 1) // Not full 4,294,967,295 because game adds 7 to this number when displaying
		result, err = db.Exec("INSERT INTO PLAYER_PYRO_ID (PYRO_ID, PLAYER_ID) VALUES (?, ?)", newId, internalPlayerId)
	}
	// if err != nil {
	// 	return 0, err
	// }

	return newId, nil
}

// GetPlayerSByID retrieves a players by their ID.
func GetPlayerByID(pyroId uint32) (uint64, error) {
	var playerId uint64
	err := db.QueryRow("SELECT PLAYER_ID FROM PLAYER_PYRO_ID WHERE PYRO_ID = ?", pyroId).Scan(&playerId)
	if err != nil {
		return 0, err
	}
	return playerId, nil
}

// // InsertContEST_ENTRY adds a new contest_entry to the CONTEST_ENTRY table.
func CreateContestEntry(entryData []byte, playerId uint64) (int64) {
	result, err := db.Exec("INSERT INTO CONTEST_ENTRY (UPLOADED, ENTRY_DATA) VALUES (datetime('now'), ?)", entryData)
	if err != nil {
		log.Fatal("Insert into CONTEST_ENTRY failed", err)
		return 0
	}

	entryId, _ := result.LastInsertId()

	_, err = db.Exec("INSERT INTO PLAYER_ENTRY (PLAYER_ID, ENTRY_ID) VALUES (?, ?)", playerId, entryId)
	if err != nil {
		log.Fatal("Insert into PLAYER_ENTRY failed", err)
	}

	return entryId
}

func FindOpponentEntry(notPlayerId uint64) (uint64, []byte, error) {
	var entryId uint64
	var entry []byte
	err := db.QueryRow("SELECT ENTRY_ID, ENTRY_DATA FROM CONTEST_ENTRY A WHERE ENTRY_ID = (SELECT MAX(ENTRY_ID) FROM PLAYER_ENTRY B WHERE B.ENTRY_ID = A.ENTRY_ID AND B.PLAYER_ID <> ?)", notPlayerId).Scan(&entryId, &entry)
	if err != nil {
		return 0, nil, err
	}

	return entryId, entry, nil
}

// InsertContest_DATA adds a new contest_data to the CONTEST_DATA table.
func InsertContest(myInternalId uint64, opponentEntryId uint64, contest_data []byte) error {
	var myEntryId uint64
	err := db.QueryRow("SELECT ENTRY_ID FROM PLAYER_ENTRY WHERE PLAYER_ID = ?", myInternalId).Scan(&myEntryId)
	if err != nil {
		log.Fatal("Failed to insert contest data", err)
	}
	_, err = db.Exec("INSERT INTO CONTEST_DATA (ENTRY_ID_0, ENTRY_ID_1, CONT_DATA) VALUES (?, ?, ?)", myInternalId, opponentEntryId, contest_data)
	return err
}

// // GetContest_DATAByID retrieves a contest_data by its ID.
// func GetContestEntryByID(db *sql.DB, contest_dataID int64) (string, error) {
// 	var contest_dataName string
// 	err := db.QueryRow("SELECT Name FROM CONTEST_DATA WHERE ID = ?", contest_dataID).Scan(&contest_dataName)
// 	if err != nil {
// 		return "", err
// 	}
// 	return contest_dataName, nil
// }

// func main() {
// 	db, err := sql.Open("sqlite3", dbFile)
// 	if err != nil {
// 		log.Fatal("Error opening database:", err)
// 	}
// 	defer db.Close()

// 	InitializeDatabase(db)

// 	// Example usage:
// 	playersID, _ := CreatePlayer(db, "John Doe")
// 	fmt.Println("New players ID:", playersID)

// 	playersName, _ := GetPlayerByID(db, playersID)
// 	fmt.Println("Player name:", playersName)

// 	_ = CreateContestEntry(db, "ContEST_ENTRY A")

// 	_ = CreateContest(db, "Summer Contest_DATA")
// 	contest_dataName, _ := GetContestByid(db, 1)
// 	fmt.Println("Contest_DATA name:", contest_dataName)
// }
