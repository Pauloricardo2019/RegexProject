package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
)

//regex const

const (
	pattern = "((?P<IPV4>(?:[0-9]{1,3}\\.){3}[0-9]{1,3})|(?P<IPV6>(?:[0-9a-z]+\\:){7}[0-9a-z]+))\\s\\-\\s\\-\\s(?P<date>\\[[0-9]+\\/[a-zA-Z]+\\/[0-9]+(?:\\:[0-9]+){3}\\s\\+[0-9]+\\])\\s\\\"(?P<Verb>[A-Z]+)*"
)
// regex var
var regex *regexp.Regexp

func main() { 
  table := &Table{}
	tables := []Table{}

	//Compila a regex criada na const pattern
	re := regexp.MustCompile(pattern)

	regex = re

	//Abre o arquivo txt de logs
	file, err := os.Open("logs")
	if err != nil {
		fmt.Println(err)
	}

	defer file.Close()

	//Passa o arquivo para recever um scanner
	scanner := bufio.NewScanner(file)
	//Divide a leitura do arquivo em linhas
	scanner.Split(bufio.ScanLines)

	//Percorre por todo o arquivo
	for scanner.Scan() {

		//Atribui o texto da linha para a variável teste
		teste := scanner.Text()

		//Procura por combinações entre minha regex e a linha; retorna os matches
		matches := regex.FindStringSubmatch(teste)
		if matches == nil {
			continue
		}

		//Procura por todos os grupos que incrementei na minha regex; e retorna os grupos
		groups := regex.SubexpNames()
		fmt.Println(matches, groups)

		//Entra na função e pega qual valor de determinado grupo eu quero pegar
		ip := getMatchedValueByIdentifier("IPV4", matches, groups)
		if ip == "" {
			ip = getMatchedValueByIdentifier("IPV6", matches, groups)
		}

		//Entra na função e pega qual valor de determinado grupo eu quero pegar
		date := getMatchedValueByIdentifier("date", matches, groups)

		//Entra na função e pega qual valor de determinado grupo eu quero pegar
		verb := getMatchedValueByIdentifier("Verb", matches, groups)

		//Atribui os valores na struct
		table.IP = ip
		table.Date = date
		table.Verb = verb

		
		tables = append(tables, *table)
	}
	if err := persistDb(tables); err != nil {
		panic(err)
	}

	fmt.Println(tables)
 fmt.Println("finish process")

}

// get the results of the match groups
func getMatchedValueByIdentifier(id string, matches []string, groups []string) string {
	for _, v := range groups {
		if v == id {
			idx := regex.SubexpIndex(v)
			return matches[idx]
		}
	}
	return ""
}


func persistDb(tables []Table) error {
	db, err := GetGormDB()
	if err != nil {
		return err
	}
	if err = db.Create(tables).Error; err != nil {
		return err
	}

	return nil
}
