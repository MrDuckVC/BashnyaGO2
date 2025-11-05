package uniqpkg

type Options struct {
	C bool // -c: подсчитать количество
	D bool // -d: вывести только дубликаты
	U bool // -u: вывести только уникальные
	I bool // -i: игнорировать регистр
	F int  // -f: пропустить N полей
	S int  // -s: пропустить N символов
}
