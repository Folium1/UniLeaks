package leaks

import "errors"

var (
	ErrVirusDetected = errors.New("Знайдений вірус в файлі")
	ErrFileNotFound     = errors.New("Файл не знайдено")
	ErrFileCheck     = errors.New("Сталась помилка, на стадії перевірки файлу на віруси")
)
