package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"uniq/uniqpkg"
)

func openFile(filename string) (*os.File, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("ошибка открытия файла %s: %w", filename, err)
	}
	return file, nil
}

func createFile(filename string) (*os.File, error) {
	file, err := os.Create(filename)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания файла %s: %w", filename, err)
	}
	return file, nil
}

func main() {
	opts := uniqpkg.Options{}
	flag.BoolVar(&opts.C, "c", false, "Подсчитать количество встречаний строки.")
	flag.BoolVar(&opts.D, "d", false, "Вывести только те строки, которые повторились.")
	flag.BoolVar(&opts.U, "u", false, "Вывести только те строки, которые не повторились.")
	flag.BoolVar(&opts.I, "i", false, "Не учитывать регистр букв.")
	flag.IntVar(&opts.F, "f", 0, "Не учитывать первые N полей.")
	flag.IntVar(&opts.S, "s", 0, "Не учитывать первые N символов.")

	flag.Parse()

	if (opts.C && opts.D) || (opts.C && opts.U) || (opts.D && opts.U) {
		fmt.Fprintln(os.Stderr, "ошибка использования: флаги -c, -d, и -u являются взаимоисключающими")
		os.Exit(1)
	}

	var reader io.Reader = os.Stdin
	var writer io.Writer = os.Stdout

	args := flag.Args()

	if len(args) > 0 {
		inputFile := args[0]
		file, err := openFile(inputFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		reader = file
		defer func() {
			if err := file.Close(); err != nil {
				fmt.Fprintf(os.Stderr, "ошибка закрытия input-файла: %v\n", err)
			}
		}()
	}

	if len(args) > 1 {
		outputFile := args[1]
		file, err := createFile(outputFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		writer = file
		defer func() {
			if err := file.Close(); err != nil {
				fmt.Fprintf(os.Stderr, "ошибка закрытия output-файла: %v\n", err)
			}
		}()
	}

	if err := uniqpkg.Uniq(reader, writer, opts); err != nil {
		fmt.Fprintln(os.Stderr, "ошибка при выполнении: ", err)
		os.Exit(1)
	}
}
