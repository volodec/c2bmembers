package main

import (
	"encoding/json"
	"flag"
	"github.com/volodec/c2bmembers/pkg/models"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

const jsonUrl = "https://qr.nspk.ru/proxyapp/c2bmembers.json"

var host *string

func main() {

	host = flag.String("host", "http://localhost", "Значение замены хоста в json")
	flag.Parse()

	prepareDirs()

	ticker(func() error {
		fileData, err := saveJson()
		if err != nil {
			return err
		}

		err = handleFileData(fileData)
		if err != nil {
			return err
		}

		return nil
	})
}

func handleFileData(fileData []byte) error {
	var data models.Data

	err := json.Unmarshal(fileData, &data)
	if err != nil {
		log.Fatalln("Ошибка парсинга файла source.json")
		return err
	}

	for _, dictionary := range data.Dictionary {
		err := saveImage(dictionary)
		if err != nil {
			return err
		}
	}

	return nil
}

func saveImage(dictionary models.Dictionary) error {
	fileName := path.Base(dictionary.LogoURL)
	filePath := "files/public/proxyapp/logo/" + fileName

	if _, err := os.Stat(filePath); !os.IsNotExist(err) { // если файл уже существует, то пропустим его
		return nil
	}

	out, err := os.Create(filePath)
	if err != nil {
		log.Fatalln("Ошибка создания файла " + fileName)
		return err
	}

	resp, err := http.Get(dictionary.LogoURL)
	if err != nil {
		log.Fatalln("Ошибка получения файла " + fileName)
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Fatalln("Ошибка сохранения файла " + fileName)
		return err
	}
	return nil
}

func saveJson() ([]byte, error) {
	fileName := "files/app/source.json"

	if err, condition := checkUpdate(); condition == true {

		out, err := os.Create(fileName)
		if err != nil {
			log.Fatalln("Ошибка создания файла source.json")
			return nil, err
		}

		resp, err := http.Get(jsonUrl)
		if err != nil {
			log.Fatalln("Ошибка получения файла source.json")
			return nil, err
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln("Не удалось прочитать тело запроса source.json")
			return nil, err
		}

		_, err = out.Write(body)
		if err != nil {
			log.Fatalln("Не удалось записать данные в файл source.json")
			return nil, err
		}

		saveResultJson(body)

		return body, nil

	} else if err != nil {
		log.Fatalln("Ошибка получения файла source.json")
		return nil, err
	}

	fileData, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatalln("Ошибка чтения файла source.json")
		return nil, err
	}

	return fileData, nil
}

func saveResultJson(data []byte) {
	preparedBody := prepareData(data)
	filePath := "files/public/proxyapp/c2bmembers.json"

	out, err := os.Create(filePath)
	if err != nil {
		log.Fatalln("Ошибка создания файла c2bmembers.json")
		return
	}

	_, err = out.Write(preparedBody)
	if err != nil {
		log.Fatalln("Не удалось записать данные в файл c2bmembers.json")
		return
	}
}

func prepareData(body []byte) []byte {
	bodyString := string(body)

	result := strings.Replace(bodyString, "https://qr.nspk.ru", *host, -1)

	return []byte(result)
}

func checkUpdate() (error, bool) {
	lastFilePath := "files/app/last"

	head, err := http.Head(jsonUrl)
	if err != nil {
		return err, false
	}

	valueFromHeader := head.Header.Get("Last-Modified")

	if _, err = os.Stat(lastFilePath); os.IsNotExist(err) {
		lastFile, _ := os.Create(lastFilePath)
		lastFile.WriteString(valueFromHeader)

		return nil, true
	}

	valueFromFile, _ := os.ReadFile(lastFilePath)

	result := string(valueFromFile) != valueFromHeader

	if !result {
		lastFile, _ := os.Open(lastFilePath)
		lastFile.WriteString(valueFromHeader)
	}

	return nil, result
}

func prepareDirs() {
	checkDir := func(dirPath string) {
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			err := os.Mkdir(dirPath, 0777)
			if err != nil {
				log.Fatalln("Не удалось создать директорию " + dirPath)
				return
			}
		}
	}

	checkDir("files")
	checkDir("files/app")
	checkDir("files/public")
	checkDir("files/public/proxyapp")
	checkDir("files/public/proxyapp/logo")
}
