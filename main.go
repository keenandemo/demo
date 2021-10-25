package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/sqlite"
)

var dbcon db.Session // database object

/* database settings. includes multiple adapters for
   postgresql, mysql, mongodb etc https://upper.io/v4/getting-started/
*/
var dbsettings = sqlite.ConnectionURL{
	Database: "./data/database",
}

var filepaths = map[string]string{
	"ingest": "./data/ingest",
	"movies": "./data/movies",
	"shows":  "./data/shows",
	//"galleries": "./data/galleries",
	//"music": "./data/music",
}

// connect to the database. in SQLite's case, this creates the db file if it doesn't already exist.
func initalizeDatabase() db.Session {
	connection, err := sqlite.Open(dbsettings)
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected to database...")

	return connection
}

func parseArgs() *Flags {
	flags := new(Flags)

	flag.BoolVar(&flags.DebugMode, "debug", false, "Launch the webserver in debug mode")
	flag.BoolVar(&flags.CreateMissingDirs, "create-missing-dirs", true, "Create the sorting and storage directories on startup if not already present.")
	flag.StringVar(&flags.DefaultPassword, "set-password", "", "Set the default access password in the database.")
	flag.Parse()

	return flags
}

// check for the storage directories on startup, defaults
func checkDirectories(createMissing bool) {

	for name, path := range filepaths {
		// fmt.Printf("skip: %+v | name: %+v | path: %+v\n", createMissing, name, path)

		_, err := os.Stat(path) // read the directory metadata from the host filesystem

		if os.IsNotExist(err) {
			fmt.Printf("Directory %+v - %+v doesn't exist\n", name, path)
		} else {
			fmt.Printf("Found directory %+v\n", path)
		}

		if createMissing && os.IsNotExist(err) {
			err := os.Mkdir(path, 0755)
			if err != nil {
				fmt.Printf("ERROR while trying to create %+v\n", path)
				panic(err)
			}

			fmt.Printf("Created directory %+v\n", path)
		}
	}
}

// stores the filepath & assigns an ID for files in the ingest folder
func parseWaitingRoom() {
	waitingFiles, err := os.ReadDir(filepaths["ingest"])
	if err != nil {
		panic(err)
	}

	for _, file := range waitingFiles {
		fmt.Println(file.Name())
		//fmt.Printf("%+v\n", file) //
	}
}

func main() {

	// connect to the database. if using SQLite, this creates the database
	dbcon = initalizeDatabase()
	defer dbcon.Close()

	args := parseArgs()

	// set the default content password
	setDefaultPassword(args.DefaultPassword)

	// check that content directories exist.
	checkDirectories(args.CreateMissingDirs)

	parseWaitingRoom() // import waiting files/folders into the db for sorting by the user

	go launchWebserver(args.DebugMode)

	select {} // stop main thread from exiting
}

type Flags struct {
	CreateMissingDirs bool
	DebugMode         bool
	DefaultPassword   string
}
