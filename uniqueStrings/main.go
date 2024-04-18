package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"strconv"
	"strings"
)

type Flag struct {
	useCount,
	useRep,
	useNoRep,
	useReg bool
	skipField,
	skipChar int
}

func uniqueCheck(s *bufio.Scanner, FlagCheck Flag) ([]string, map[string]int, map[string]int) {
	var unStr []string                     // для правильного порядка вывода
	uniqueMap := make(map[string]int)      // считаем уникальные строки
	uniqueMapCheck := make(map[string]int) // считаем уникальные строки
	for s.Scan() {
		signifChars := s.Text()
		if FlagCheck.skipField > 0 { //не учитываем поля(до пробела)
			spaceInd := strings.Index(signifChars, " ")

			for i := 0; i < FlagCheck.skipField; i++ {
				signifChars = signifChars[(spaceInd + 1):]
				spaceInd = strings.Index(signifChars, " ")
			}
		}

		signifChars = signifChars[FlagCheck.skipChar:] //не учитываем первые n символов

		if FlagCheck.useReg {
			signifChars = strings.ToLower(signifChars)
		}

		_, inMap := uniqueMap[signifChars]

		//насчет флага первых символов, попробуй создать мапу(или срез) куда будешь пихать строки без первых n cимволов
		if s.Text() != "" { // не нужный элемент перевода строки
			if !inMap {
				uniqueMap[signifChars] = 1
				uniqueMapCheck[s.Text()] = 1
				unStr = append(unStr, s.Text())
			} else {
				uniqueMap[signifChars] += 1
				uniqueMapCheck[s.Text()] += 1
			}
			if s.Text() == "Thanks." {
				break
			} //мб не надо тк в os.Stdin также лимитное количество сообщений
		}
	}
	return unStr, uniqueMap, uniqueMapCheck
}

func ScannerS(arguments []string) *os.File {
	file := os.Stdin

	if len(arguments) >= 1 && (arguments[0] == "input_file.txt" || arguments[1] == "input_file.txt") {
		file, err := os.Open("input_file.txt")
		if err != nil {
			log.Fatal("Problems with opening file")
		}

		defer file.Close()
	}

	return file
}

func WriterS(arguments, unStr []string, uniqueMap, uniqueMapCheck map[string]int, FlagCheck Flag) *os.File {
	file := os.Stdout

	if len(arguments) >= 1 && (arguments[0] == "output_file.txt" || arguments[1] == "output_file.txt") {
		file, err := os.Open("output_file.txt")
		if err != nil {
			log.Fatal("Problems with opening file")
		}

		defer file.Close()
	}

	var data string

	if !FlagCheck.useRep && FlagCheck.useNoRep {
		if FlagCheck.useCount {
			for _, value := range unStr {
				for key, valMap := range uniqueMapCheck {
					if (key == value) && (valMap == 1) {
						data = "1 " + value
						file.WriteString(data)
					}
				}
			}
		} else {
			for _, value := range unStr {
				for key, valMap := range uniqueMapCheck {
					if (key == value) && (valMap == 1) {
						file.WriteString(value)
					}
				}
			}
		}
	} else if !FlagCheck.useNoRep && FlagCheck.useRep {
		if FlagCheck.useCount {
			for _, value := range unStr {
				for key, valMap := range uniqueMapCheck {
					if (key == value) && (valMap > 1) {
						data = strconv.Itoa(valMap) + value
						file.WriteString(data)
					}
				}
			}
		} else {
			for _, value := range unStr {
				for key, valMap := range uniqueMapCheck {
					if (key == value) && (valMap > 1) {
						file.WriteString(value)
					}
				}
			}
		}
	} else if !FlagCheck.useNoRep && !FlagCheck.useRep {
		return file
	} else {
		if FlagCheck.useCount {
			for _, value := range unStr {
				data = strconv.Itoa(uniqueMapCheck[value]) + value
				file.WriteString(data)
			}
		} else {
			for _, value := range unStr {
				file.WriteString(value)
			}
		}
	}

	return file

}

func main() {
	useCount := flag.Bool("c", false, "display count of unique strings")
	useRep := flag.Bool("d", false, "display repeated strings")
	useNoRep := flag.Bool("u", false, "display not repeated strings")
	useReg := flag.Bool("i", false, "Don't need your YEHFWJF letters")
	var skipField int
	var skipChar int
	flag.IntVar(&skipField, "f", 0, "number of started lines which we delete")
	flag.IntVar(&skipChar, "s", 0, "number of started letters which we delete")

	flag.Parse()

	FlagCheck := &Flag{*useCount, *useRep, *useNoRep, *useReg, skipField, skipChar} // подумай куда передавать

	arguments := flag.Args() // было os.Args, уточни насчет нумерации с 0 или 1

	f := bufio.NewScanner(ScannerS(arguments))
	unStr, uniqueMap, uniqueMapCheck := uniqueCheck(f, *FlagCheck)
	WriterS(arguments, unStr, uniqueMap, uniqueMapCheck, *FlagCheck)

}

/* с флагами
func Init(filename *string) *os.File{
	flag.Parse()
	if *filename == "" {
		*filename = "noInpF"
		return os.Stdin
	} else {
		inputF, err := os.Open(*filename)
		if err != nil {
			log.Fatalf("Error opening file: %v, ", err)
		}
		defer inputF.Close()
		return inputF
	}
}
*/

/* с флагами
filename := flag.String("filename", "", "Filename to read from")

f := bufio.NewScanner(Init(filename))
*/

// inputF, err := os.Open("input_file.txt")
// defer inputF.Close()

// if err != nil {
// 	s := bufio.NewScanner(os.Stdin)
// } else {
// 	s := bufio.NewScanner(inputF)
// }

// outputF, err := os.Open("output_file.txt")
// defer outputF.Close()

// if err != nil {
// 	n := bufio.NewWriter(outputF)
// } else {
// 	n := bufio.NewWriter(os.Stdout)
// }

// uniqueCheck(s, n)
