package uniqpkg

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetComparisonKey(t *testing.T) {
	testCases := []struct {
		name     string
		line     string
		opts     Options
		expected string
	}{
		{"Без опций", "Hello World", Options{}, "Hello World"},
		{"Флаг -i", "Hello World", Options{I: true}, "hello world"},
		{"Флаг -f 1", "first second third", Options{F: 1}, "second third"},
		{"Флаг -f 3 (больше или равно)", "first second third", Options{F: 3}, ""},
		{"Флаг -s 2", "Hello World", Options{S: 2}, "llo World"},
		{"Флаг -s 100 (больше длины)", "Hello World", Options{S: 100}, ""},
		{"Флаг -f 1 и -s 3", "skip_me and_this now_compare", Options{F: 1, S: 3}, "_this now_compare"},
		{"Флаг -i, -f 1, -s 3", "skip_me AND_THIS Now_Compare", Options{I: true, F: 1, S: 3}, "_this now_compare"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := getComparisonKey(tc.line, tc.opts)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestUniq(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		opts     Options
		expected string
	}{
		{
			name:     "Простой тест (без флагов)",
			input:    "a\na\nb\nc\nc\nc\nd",
			opts:     Options{},
			expected: "a\nb\nc\nd\n",
		},
		{
			name:     "Пустой ввод",
			input:    "",
			opts:     Options{},
			expected: "",
		},
		{
			name:     "Ввод с пустыми строками",
			input:    "a\na\n\n\nb\n",
			opts:     Options{},
			expected: "a\n\nb\n",
		},
		{
			name:     "Флаг -c (Count)",
			input:    "a\na\nb\nc\nc\nc\nd\n",
			opts:     Options{C: true},
			expected: "2 a\n1 b\n3 c\n1 d\n",
		},
		{
			name:     "Флаг -d (Duplicates only)",
			input:    "a\na\nb\nc\nc\nc\nd",
			opts:     Options{D: true},
			expected: "a\nc\n",
		},
		{
			name:     "Флаг -u (Uniques only)",
			input:    "a\na\nb\nc\nc\nc\nd",
			opts:     Options{U: true},
			expected: "b\nd\n",
		},
		{
			name:     "Флаг -i (Ignore case)",
			input:    "a\nA\nb\nC\nc\nc\nd",
			opts:     Options{I: true},
			expected: "a\nb\nC\nd\n",
		},
		{
			name:     "Флаг -i и -c",
			input:    "a\nA\nb\nC\nc\nc\nd",
			opts:     Options{I: true, C: true},
			expected: "2 a\n1 b\n3 C\n1 d\n",
		},
		{
			name:     "Флаг -f 1 (Skip 1 field)",
			input:    "1 a\n2 a\n3 b\n4 b\n5 c",
			opts:     Options{F: 1},
			expected: "1 a\n3 b\n5 c\n",
		},
		{
			name:     "Флаг -s 2 (Skip 2 chars)",
			input:    "xxa\nxxa\nyyb\nyyc",
			opts:     Options{S: 2},
			expected: "xxa\nyyb\nyyc\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			inputReader := strings.NewReader(tc.input)
			var outputWriter bytes.Buffer

			err := Uniq(inputReader, &outputWriter, tc.opts)

			require.NoError(t, err)
			require.Equal(t, tc.expected, outputWriter.String())
		})
	}
}
