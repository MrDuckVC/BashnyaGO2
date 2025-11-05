package uniqpkg

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

func getComparisonKey(line string, opts Options) string {
	key := line

	if opts.I {
		key = strings.ToLower(key)
	}

	if opts.F > 0 {
		fields := strings.Fields(key)
		if len(fields) > opts.F {
			key = strings.Join(fields[opts.F:], " ")
		} else {
			key = ""
		}
	}

	if opts.S > 0 {
		if len(key) > opts.S {
			key = key[opts.S:]
		} else {
			key = ""
		}
	}

	return key
}

func Uniq(reader io.Reader, writer io.Writer, opts Options) error {
	scanner := bufio.NewScanner(reader)

	var prevLine *string
	count := 0

	flush := func() error {
		if prevLine == nil {
			return nil
		}

		line := *prevLine
		var err error

		switch {
		case opts.C:
			_, err = fmt.Fprintf(writer, "%d %s\n", count, line)
		case opts.D:
			if count > 1 {
				_, err = fmt.Fprintln(writer, line)
			}
		case opts.U:
			if count == 1 {
				_, err = fmt.Fprintln(writer, line)
			}
		default:
			_, err = fmt.Fprintln(writer, line)
		}

		return err
	}

	for scanner.Scan() {
		currentLine := scanner.Text()

		if prevLine == nil {
			prevLine = &currentLine
			count = 1
			continue
		}

		prevKey := getComparisonKey(*prevLine, opts)
		currentKey := getComparisonKey(currentLine, opts)

		if prevKey == currentKey {
			count++
		} else {
			if err := flush(); err != nil {
				return err
			}
			prevLine = &currentLine
			count = 1
		}
	}

	if err := flush(); err != nil {
		return err
	}

	return scanner.Err()
}
