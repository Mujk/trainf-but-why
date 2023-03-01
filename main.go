package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
)

type train struct {
	dirc, futrDirc [2]int8
	pos            [2]int
	cells          map[int]int
	cellPos        int
}

func errHand(err error) {
	if err != nil {
		panic(err)
	}
}

func checkFileGiven(args []string) error {
	if len(os.Args) == 2 {
		return nil
	}
	return errors.New("Expected 1 file, got " + strconv.Itoa(len(args)-1))
}

func checkExt(file string) error {
	findExt := strings.Split(file, ".")
	if findExt[len(findExt)-1] == "trainf" {
		return nil
	}
	return errors.New("not a .trainf file")
}

func findInRow(row string, symbol string) []int {
	symbols := []int{}
	for x := range row {
		if string(row[x]) == symbol {
			symbols = append(symbols, x)
		}
	}
	return symbols
}

func findChar(code []string, char string) ([]string, error) {
	stations := []string{}
	for yPos, row := range code {
		xPosArr := findInRow(row, "+")
		for _, xPos := range xPosArr {
			stations = append(stations, fmt.Sprintf("%d%s%d", yPos, ",", xPos))
		}
	}
	if len(stations) == 0 {
		return stations, errors.New("Now trainstation + found")
	}
	return stations, nil
}

func posToInt(pos string) [2]int {
	splitPos := strings.Split(pos, ",")
	posInt := [2]int{}
	var err error
	for i := range splitPos {
		posInt[i], err = strconv.Atoi(splitPos[i])
		errHand(err)
	}
	return posInt
}

func move(pos [2]int, dirc [2]int8) [2]int { // breaks code, find bug
	for i, direction := range dirc {
		pos[i] += int(direction)
	}
	return pos
}

func userInput() (int, error) {
	input := ""
	fmt.Scanln(&input)
	if len(input) <= 0 {
		errHand(errors.New("Expected 1 input character got " + strconv.Itoa(len(input))))
	}
	givenLett := int(input[0])
	if len(input) >= 2 {
		return givenLett, errors.New("Expected 1 input character got " + strconv.Itoa(len(input)) + "\nWill continue with the first character: " + string(input[0]))
	}
	return givenLett, nil
}

func commands(trainInfo train, cmd rune) train { // switch breaks the code
	switch cmd {
	case 'h':
		fmt.Println("Test: Hi")
	case '^':
		trainInfo.futrDirc[1] = 1
	case 'v':
		trainInfo.futrDirc[1] = -1
	case '<':
		trainInfo.futrDirc[0] = -1
	case '>':
		trainInfo.futrDirc[0] = 1
	case '+':
		if trainInfo.dirc[0] != 0 {
			trainInfo.cellPos += int(trainInfo.dirc[0])
		} else {
			trainInfo.cells[trainInfo.cellPos] += int(trainInfo.dirc[1])
		}
	case '.':
		// output
	case ',':
		userIn, err := userInput()
		if err != nil {
			fmt.Println(err)
		}
		trainInfo.cells[trainInfo.cellPos] = userIn
	case 'o':
		trainInfo.dirc = trainInfo.futrDirc
	}
	return trainInfo
}

func runTrain(code []string, trainInfo train, crashPos chan [2]int, wgRound *sync.WaitGroup, wgControl *sync.WaitGroup, wgCrashed *sync.WaitGroup, pos chan [2]int) {
	for i := 0; i < 2; i++ { // make until crash later
		fmt.Println(trainInfo.pos)
		//fmt.Printf("Test: %c\n", code[trainInfo.pos[0]][trainInfo.pos[1]])
		trainInfo = commands(trainInfo, rune(code[trainInfo.pos[0]][trainInfo.pos[1]]))
		pos <- trainInfo.pos
		defer wgControl.Wait()
		defer wgRound.Done()
	}
	wgCrashed.Done()
}

func trainControl(wgRound *sync.WaitGroup, wgControl *sync.WaitGroup, trainsNum int, pos chan [2]int) {
	for trainsNum < 0 {
		defer wgRound.Wait()
		for posi := range <-pos {

		}
		wgControl.Add(1)
		wgRound.Add(trainsNum)
		defer wgControl.Done()
	}
}

func trainMath() {

}

func spawnTrains(code []string, stations []string) {
	trainsNum := len(stations) * 4
	wgRound := sync.WaitGroup{}
	trains := 0
	crashPos := make(chan [2]int)
	wgCrashed := sync.WaitGroup{}
	wgControl := sync.WaitGroup{}
	pos := make(chan [2]int)
	//numPos := make(chan int)
	//numToAdd := make(chan int)
	dircs := [4][2]int8{{0, 1}, {0, -1}, {1, 0}, {-1, 0}} // up, down, right, left
	go trainControl(&wgRound, &wgControl, trainsNum, pos)
	for i := 0; trains < trainsNum; i++ {
		for j := 0; j < 4; j++ {
			trainInfo := train{
				pos:  posToInt(stations[i]),
				dirc: dircs[j],
			}
			go runTrain(code, trainInfo, crashPos, &wgRound, &wgControl, &wgCrashed, pos)
			go trainMath()
			trains++
			wgCrashed.Add(1)
		}
	}
	wgCrashed.Wait()
	fmt.Println("test: ", "bye")
}

func main() {
	errHand(checkFileGiven(os.Args))
	errHand(checkExt(os.Args[1]))
	file, err := os.Open(os.Args[1])
	errHand(err)
	var code []string
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		code = append(code, scanner.Text())
	}
	stations, err := findChar(code, "+")
	errHand(err)
	defer file.Close()
	spawnTrains(code, stations)
}
