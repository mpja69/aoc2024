package main

import (
	"aoc2024/day03/lexer"
	"aoc2024/day03/token"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const (
	SKIP_STATE = iota
	LP_STATE
	A_STATE
	COMMA_STATE
	B_STATE
	RP_STATE
)

func part1_with_parser(str string) (sum int) { //, idx [][]int) {
	var err error
	a, b := 0, 0
	// i := 0
	l := lexer.New(str)
	state := SKIP_STATE
	for t := l.NextToken(); t.Type != token.EOF; t = l.NextToken() {

		// Always check if a new MUL starts
		if t.Type == token.MUL {
			state = LP_STATE
			continue
		}

		switch state {
		case LP_STATE:
			// Check if we can continue the MUL...otherwise go back to START
			if t.Type == token.LPAREN {
				state = A_STATE
			} else {
				state = SKIP_STATE
			}
		case A_STATE:
			// Check if we can continue the MUL...otherwise go back to START
			if t.Type == token.INT {
				a, err = strconv.Atoi(t.Literal)
				if err != nil {
					log.Fatalf("%v", err)
				}
				state = COMMA_STATE
			} else {
				state = SKIP_STATE
			}
		case COMMA_STATE:
			// Check if we can continue the MUL...otherwise go back to START
			if t.Type == token.COMMA {
				state = B_STATE
			} else {
				state = SKIP_STATE
			}
		case B_STATE:
			// Check if we can continue the MUL...otherwise go back to START
			if t.Type == token.INT {
				b, err = strconv.Atoi(t.Literal)
				if err != nil {
					log.Fatalf("%v", err)
				}
				state = RP_STATE
			} else {
				state = SKIP_STATE
			}
		case RP_STATE:
			// Check if we can continue the MUL...AND actually finish off the MUL...otherwise go back to START
			if t.Type == token.RPAREN {
				sum += a * b
				state = SKIP_STATE
			} else {
				state = SKIP_STATE
			}
		}
	}
	return sum //, idx
}

func part2_with_parser(str string) (sum int) {
	var err error
	a, b := 0, 0
	enable := true
	l := lexer.New(str)
	state := SKIP_STATE
	for t := l.NextToken(); t.Type != token.EOF; t = l.NextToken() {

		if t.Type == token.IDENT {
			fmt.Println(t.Literal)
		}

		// Always check if a new MUL starts
		if t.Type == token.MUL {
			state = LP_STATE
			continue
		}

		if t.Type == token.DO {
			enable = true
			continue
		}

		if t.Type == token.DONT {
			enable = false
			continue
		}

		switch state {
		case LP_STATE:
			// Check if we can continue the MUL...otherwise go back to START
			if t.Type == token.LPAREN {
				state = A_STATE
			} else {
				state = SKIP_STATE
			}
		case A_STATE:
			// Check if we can continue the MUL...otherwise go back to START
			if t.Type == token.INT {
				a, err = strconv.Atoi(t.Literal)
				if err != nil {
					log.Fatalf("%v", err)
				}
				state = COMMA_STATE
			} else {
				state = SKIP_STATE
			}
		case COMMA_STATE:
			// Check if we can continue the MUL...otherwise go back to START
			if t.Type == token.COMMA {
				state = B_STATE
			} else {
				state = SKIP_STATE
			}
		case B_STATE:
			// Check if we can continue the MUL...otherwise go back to START
			if t.Type == token.INT {
				b, err = strconv.Atoi(t.Literal)
				if err != nil {
					log.Fatalf("%v", err)
				}
				state = RP_STATE
			} else {
				state = SKIP_STATE
			}
		case RP_STATE:
			// Check if we can continue the MUL...AND actually finish off the MUL...otherwise go back to START
			if t.Type == token.RPAREN {
				if enable {
					sum += a * b
				}
				state = SKIP_STATE
			} else {
				state = SKIP_STATE
			}
		}
	}
	return sum //, idx
}

func part1_with_regex(str string) int {
	re := regexp.MustCompile(`mul\((\d+),(\d+)\)`)
	match := re.FindAllStringSubmatch(str, -1)
	res := 0
	for _, sl := range match {
		a, _ := strconv.Atoi(sl[1])
		b, _ := strconv.Atoi(sl[2])
		res += a * b
	}
	return res
}

func main() {
	// Get the data into strings
	data, err := os.ReadFile("d.txt")
	if err != nil {
		log.Fatalf("%v", err)
	}
	str := string(data)
	str = strings.Trim(str, " \n")
	// lines := strings.Split(str, "\n")

	res := part1_with_parser(str)
	fmt.Println("Part 1 - with parser:", res)

	res = part1_with_regex(str)
	fmt.Println("Part 1 - with regex: ", res)

	res = part2_with_parser(str)
	fmt.Println("Part 2 - with parser:", res)

}
