package leaks

import "errors"

var (
	VirusDetectedErr = errors.New("Знайдений вірус в файлі")
	FileNotFound     = errors.New("Файл не знайдено")
	FileCheckErr     = errors.New("Сталась помилка, на стадії перевірки файлу на віруси")
)
