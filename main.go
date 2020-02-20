package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

type library struct {
	numBooks     int
	signUp       int
	books        []int //Value is id
	scansPerDay  int
	booksScanned []int
	score        float64

	onebook int

	libidx int
}

type meta struct {
	bookCount    int
	libCount     int
	dayCount     int
	scores       []int
	libraries    []library
	libraryOrder []*library
}

var bookAmt map[int]int

func printOutput() {
	fmt.Println(len(metadata.libraryOrder))

	for _, lib := range metadata.libraryOrder {
		if len(lib.booksScanned) == 0 {
			fmt.Println(lib.libidx, 1)
			fmt.Println(lib.onebook)
		} else {
			fmt.Println(lib.libidx, len(lib.booksScanned))
			printed := false
			for _, book := range lib.booksScanned {
				if printed {
					fmt.Printf(" ")
				}
				fmt.Printf("%d", book)
				printed = true
			}
			fmt.Println()
		}
	}
}

func readLibrary(reader *bufio.Reader) (library, error) {
	lib := library{}
	linea, err := reader.ReadString('\n')
	if err != nil {
		return lib, err
	}

	libData := strings.Fields(linea)
	if len(libData) == 0 {
		return lib, fmt.Errorf("EMPTY LINE BSAD")
	}
	lib.numBooks = atoi(libData[0])
	lib.signUp = atoi(libData[1])
	lib.scansPerDay = atoi(libData[2])
	lib.books = make([]int, lib.numBooks)
	lineb, err := reader.ReadString('\n')
	if err != nil {
		return lib, err
	}

	books := strings.Fields(lineb)
	i := 0
	for _, book := range books {
		name := atoi(book)
		lib.onebook = name
		// lib.books[i] = metadata.scores[name]
		amt, ok := bookAmt[name]
		if !ok {
			bookAmt[name] = 1
		} else {
			bookAmt[name] = amt + 1
		}
		lib.books[i] = name
		i++
	}

	return lib, nil
}

var metadata *meta

func atoi(s string) int {
	s = strings.TrimSpace(s)
	i, _ := strconv.Atoi(s)

	return i
}

func setMetaData(line string) {
	metaFields := strings.Split(line, " ")

	metadata = &meta{
		bookCount: atoi(metaFields[0]),
		libCount:  atoi(metaFields[1]),
		dayCount:  atoi(metaFields[2]),
	}
}

func setBookScore(line string) {
	fields := strings.Split(line, " ")

	metadata.scores = make([]int, metadata.bookCount)

	for i, v := range fields {
		metadata.scores[i] = atoi(v)
	}
}

func readFile(path string) {
	file, err := os.Open(path)

	if err != nil {
		log.Fatal(err)
	}

	buf := bufio.NewReader(file)

	metaLine, err := buf.ReadString('\n')

	if err != nil {
		log.Fatal(err)
	}

	setMetaData(metaLine)

	scoreLine, err := buf.ReadString('\n')

	if err != nil {
		log.Fatal(err)
	}

	setBookScore(scoreLine)

	bookAmt = make(map[int]int)

	metadata.libraries = make([]library, metadata.libCount)
	i := 0
	for lib, err := readLibrary(buf); err == nil; lib, err = readLibrary(buf) {
		//fmt.Println(lib)
		lib.libidx = i
		metadata.libraries[i] = lib
		i++
	}

}

func solve() {
	libs := make([]*library, 0, 1)
	for i := range metadata.libraries {
		libs = append(libs, &metadata.libraries[i])
	}
	daysLeft := metadata.dayCount
	score := 0
	var signingLibrary *library
	signingLeft := 0

	activeLibrary := make([]*library, 0, 1) //Currently scanning libraries
	scannedBooks := make(map[int]bool)      //Scanned books
	for daysLeft > 0 {

		// process signing up library
		if signingLeft == 0 {
			if len(libs)%1 == 0 {
				for idx, lib := range libs {
					for i := len(lib.books) - 1; i >= 0; i-- {
						id := lib.books[i]
						_, ok := scannedBooks[id]
						if ok {
							libs[idx].books = append(lib.books[:i], lib.books[i+1:]...)
						}
					}
					libs[idx].score = score_ptr(lib, daysLeft)
				}

				sort.Slice(libs,
					func(i, j int) bool {
						scorei := libs[i].score
						scorej := libs[j].score
						return scorei > scorej
					})

			}
			if signingLibrary != nil {
				activeLibrary = append(activeLibrary, signingLibrary)
				metadata.libraryOrder = append(metadata.libraryOrder, signingLibrary)
			}
			for len(libs) > 0 {
				for len(libs[0].books) > 0 {
					id := libs[0].books[0]
					_, ok := scannedBooks[id]
					if ok {
						libs[0].books = libs[0].books[1:]
					} else {
						break
					}
				}
				if len(libs[0].books) > 0 {
					break
				}
				libs = libs[1:]
			}
			if len(libs) > 0 {
				signingLibrary = libs[0]
				signingLeft = signingLibrary.signUp
				libs = libs[1:]
			}
		}
		signingLeft--

		// Process all Scans
		for i := len(activeLibrary) - 1; i >= 0; i-- {
			ScanLeft := activeLibrary[i].scansPerDay
			for ScanLeft > 0 {
				if len(activeLibrary[i].books) == 0 {
					activeLibrary = append(activeLibrary[:i], activeLibrary[i+1:]...)
					break
				}
				popped := false
				//Get arbitary books and scan
				for !popped {
					if len(activeLibrary[i].books) == 0 {
						break
					}
					id := activeLibrary[i].books[0]
					_, ok := scannedBooks[id]
					if !ok {
						scannedBooks[id] = true
						score += metadata.scores[id]
						activeLibrary[i].booksScanned = append(activeLibrary[i].booksScanned, id)
						popped = true
					}
					activeLibrary[i].books = activeLibrary[i].books[1:]
				}
				ScanLeft--
			}
		}
		daysLeft--
	}
}

func score(lib library, daysLeft int) float64 {
	speeed := float64(lib.scansPerDay) / float64(lib.signUp)

	multi := []float64{8.0, 4.0, 2.0, 1.5, 1.0}

	//numBooks := lib.numBooks
	perspScore := 0.0
	booksCount := (daysLeft - lib.signUp) * lib.scansPerDay
	for _, id := range lib.books {
		if booksCount == 0 {
			break
		}
		perspScore += float64(metadata.scores[id]) * multi[int(math.Min(float64(5.0), float64(bookAmt[id])))-1]
		booksCount--
	}
	return float64(perspScore) / float64(1+lib.scansPerDay) * speeed
}

func score_ptr(lib *library, daysLeft int) float64 {
	return score(*lib, daysLeft)
}

// sort books based on score
func sortBooks(books []int) []int {
	sort.Slice(books,
		func(i, j int) bool {
			namei := books[i]
			namej := books[j]
			return metadata.scores[namei] > metadata.scores[namej]
		})
	return books
}

func main() {
	// files := []string{
	// 	"a_example.txt",
	// 	"b_read_on.txt",
	// 	"c_incunabula.txt",
	// 	"d_tough_choices.txt",
	// 	"e_so_many_books.txt",
	// 	"f_libraries_of_the_world.txt",
	// }

	flag.Parse()
	args := flag.Args()

	if len(args) < 1 {
		fmt.Println("GIVE ARGUMENTS FUCKER!!!")
		os.Exit(69)
	}

	// for _, v := range files {
	// 	readFile("input/" + v)

	// 	fmt.Printf("%+v\n", metadata)
	// }

	readFile(args[0])

	for idx, lib := range metadata.libraries {
		metadata.libraries[idx].score = score(lib, metadata.dayCount)
		sortBooks(lib.books)
	}

	sort.Slice(metadata.libraries,
		func(i, j int) bool {
			scorei := metadata.libraries[i].score
			scorej := metadata.libraries[j].score
			return scorei > scorej
		})

	solve()
	printOutput()
	// fmt.Printf("%+v\n"``, metadata)
}
