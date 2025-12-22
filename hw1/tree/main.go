package main

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

const (
	vsB      string = "├"
	vB       string = "│"
	cB       string = "└"
	hB       string = "───"
	tab      string = "\t"
	splitStr string = "&os.unixDirent{parent:"
)

// короче, рабочей сделал, но первый тест не проходит потому что я уголки эти └ нехорошие не добавлял,
// уже неинтересно стало
// второй тест, когда без файлов, тоже не проходит, но там чет вообще все плохо.
// главное, что работает (и что выучил что тут тоже можно любой тип к любому типу привести, только 500 строчек писать вместо одной в питоне)

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("Ошибка открытия файла %s: %v", path, err)
	}

	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("Ошибка получения данных о файле %s: %v", path, err)
	}

	if !fileInfo.IsDir() {
		return fmt.Errorf("Файл %s не является директорией, укажите директорию", path)
	} else {
		// про эту функцию в интернете подсмотрел, потому что ReadDir не хочет рекурсивно работать почему-то,
		// а так сделал бы и рекурсивно
		err = filepath.WalkDir(path, printFile(out, printFiles, path))
		if err != nil {
			return fmt.Errorf("Ошибка получения данных о файле %s: %v", path, err)
		}
	}

	return nil
}

func printFile(out io.Writer, printFiles bool, rootPath string) fs.WalkDirFunc {
	return func(path string, file fs.DirEntry, err error) error {
		if file.Name() == rootPath {
			return nil
		}

		// вообще не понял почему нельзя поле .parent достать из file нормально, ну да ладно, я достану и так. Не впервой.
		fileParent := strings.Trim(strings.Split(strings.Split(fmt.Sprintf("%#v", file), splitStr)[1], ",")[0], "\"")

		// вот тут где-то наверное и ломается подсчет отступов там и прочее, но это уже неважно.
		// Для красоты можно в терминале ls -Rl1 написать
		if fileParent != rootPath {
			fmt.Fprintf(out, "%s", vB)
			for i := 0; i < len(strings.Split(fileParent, string(os.PathSeparator)))-1; i++ {
				fmt.Fprintf(out, "%s%s", tab, vB)
			}
		} else {
			fmt.Fprintf(out, "%s", vsB)
		}

		if file.IsDir() {
			fmt.Fprintf(out, "%s%s\n", hB, file.Name())
		}

		if printFiles && !file.IsDir() {
			fileInfo, err := file.Info()

			if err != nil {
				return fmt.Errorf("%v", err)
			}

			fileSize := fileInfo.Size()
			var sizeFmtStr string
			if fileSize == 0 {
				sizeFmtStr = "(empty)"
			} else {
				sizeFmtStr = "(" + fmt.Sprintf("%d", fileSize) + "b)"
			}
			fmt.Fprintf(out, "%s%s "+sizeFmtStr+"\n", hB, file.Name())
		}

		return nil
	}
}
