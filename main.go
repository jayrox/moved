package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

var (
	flagdebug      = flag.Bool("d", true, "show debug output")
	flagdir        = flag.String("dir", "cwd", "directory to scan. default is current working directory (cwd)")
	flagtarget     = flag.String("tar", "cwd", "directory to extract to. default is current working directory (cwd)")
	flagminsize    = flag.Int64("min", 200000000, "minimum file size to include in scan. default is 200MB") // 3MB
	flagmovesample = flag.Bool("ms", false, "move sample files")
	rf             mvdFlags
	moveext        = []string{
		".mkv", ".mp4", ".avi", ".m4v", ".divx",
	}
)

func main() {
	flag.Parse()
	// Print the logo :P
	printLogo()

	rf.Min = flagInt(flagminsize)

	// Root folder to scan
	fpSAbs, _ := filepath.Abs(flagString(flagdir))
	rf.Dir = fpSAbs
	if flagString(flagdir) == "cwd" {
		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			log.Fatal(err)
		}
		rf.Dir = dir
	}

	// Root folder to move to
	fpTAbs, _ := filepath.Abs(flagString(flagtarget))
	rf.Target = fpTAbs
	if flagString(flagtarget) == "cwd" {
		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			log.Fatal(err)
		}
		rf.Target = dir
	}

	fmt.Printf("Scanning directory: %s\n", rf.Dir)
	fmt.Println("_____________________")

	i := folderWalk(rf.Dir)
	if i < 1 {
		fmt.Println("No movable files found.")
	}
}

func flagString(fs *string) string {
	return fmt.Sprint(*fs)
}

func flagInt(fi *int64) int64 {
	return int64(*fi)
}

func flagBool(fb *bool) bool {
	return bool(*fb)
}

func folderWalk(file string) (i int64) {
	i = 0
	var err = filepath.Walk(file, func(file string, _ os.FileInfo, _ error) error {
		for _, x := range moveext {
			if !flagBool(flagmovesample) && strings.Contains(strings.ToLower(file), "sample") {
				//fmt.Println("Skipping sample.")
				continue
			}
			if filepath.Ext(file) == x {
				var ok bool

				ok = moveable(file)

				if ok == true {
					ok = move(file, rf.Target)

					if ok == false {
						printDebug("Move failed %s\n", "")
					}
				}
				fmt.Println("_________")
			}
		}
		return nil
	})
	if err != nil {
		printDebug("Error: %+v\n", err)
	}
	return
}

func moveable(file string) bool {
	fmt.Printf("Checking file size: %s\n", file)

	sizeOne := fileSize(file)
	fmt.Println(sizeOne)

	time.Sleep(200 * time.Millisecond)
	sizeTwo := fileSize(file)
	fmt.Println(sizeTwo)

	time.Sleep(200 * time.Millisecond)
	sizeThree := fileSize(file)
	fmt.Println(sizeThree)

	if sizeOne == sizeTwo && sizeOne == sizeThree {
		return true
	}
	return false
}

func fileSize(path string) int64 {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
		return -1
	}
	fi, err := file.Stat()
	if err != nil {
		log.Fatal(err)
		return -1
	}
	file.Close()
	return fi.Size()
}

func move(file, destpath string) (ok bool) {
	ok = true

	destpath = strings.Replace(destpath, "\\", "/", -1)
	d := path.Dir(destpath)

	file = strings.Replace(file, "\\", "/", -1)
	target := d + "/" + path.Base(file)

	printDebug("Moving: %s\nDestination: %s\n", file, target)

	err := os.Rename(file, target)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}

// Check err
func check(e error) {
	if e != nil {
		panic(e)
	}
}

// Only print debug output if the debug flag is true
func printDebug(format string, vars ...interface{}) {
	if *flagdebug {
		if vars[0] == nil {
			fmt.Println(format)
			return
		}
		fmt.Printf(format, vars...)
	}
}

// Hold flag data
type mvdFlags struct {
	Dir    string
	Target string
	Debug  bool
	Min    int64
}

// Print the logo, obviously
func printLogo() {
	fmt.Println("███╗   ███╗██╗   ██╗██████╗")
	fmt.Println("████╗ ████║██║   ██║██╔══██╗")
	fmt.Println("██╔████╔██║██║   ██║██║  ██║")
	fmt.Println("██║╚██╔╝██║╚██╗ ██╔╝██║  ██║")
	fmt.Println("██║ ╚═╝ ██║ ╚████╔╝ ██████╔╝")
	fmt.Println("╚═╝     ╚═╝  ╚═══╝  ╚═════╝ moved")
	fmt.Println("")
}
