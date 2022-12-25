package main

import (
	"advent-of-code/aoc"
	"fmt"
)

func main() {
	filename := aoc.GetFilename()
	lines := aoc.GetInputLines(filename)

	fmt.Println(part1(lines))
}

func part1(lines []string) string {
	total := 0
	for _, line := range lines {
		total += decodeSnafu(line)
	}
	return encodeSnafu(total)
}

func decodeSnafu(digits string) int {
	value := 0

	for _, digit := range digits {

		value *= 5;

		switch digit {

		case '=': value -= 2
		case '-': value -= 1
		case '0': value += 0
		case '1': value += 1
		case '2': value += 2
		default: panic(digit)

		}
	}
	return value
}

/*
fun to_pentary_n($n) {
    return 0 if $n == 0;

    my @digits;
    while ($n > 0) {
        my $carry = 0;
        my $digit = $n % 5;
        if ($digit == 4) {
            $digit = '-';
            $carry = 1;
        }
        elsif ($digit == 3) {
            $digit = '=';
            $carry = 1;
        }
        push(@digits, $digit);
        $n = int($n / 5) + $carry;
    }

    return join('', reverse @digits);
}
*/

func encodeSnafu(n int) string {
	if n == 0 {
		return "0"
	}

	digits := ""
	for n > 0 {
		carry := 0
		var digit rune
		switch n % 5 {

		case 0: digit = '0'
		case 1: digit = '1'
		case 2: digit = '2'
		case 3: digit = '='; carry = 1
		case 4: digit = '-'; carry = 1

		}
		digits = fmt.Sprintf("%c%s", digit, digits)
		n = n / 5 + carry
	}
	return digits
}
