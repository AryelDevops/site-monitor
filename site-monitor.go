package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const monitoramento = 3
const delay = 5

var fileName, logsPath = pegaArgs()

func main() {

	exibeIntroducao()

	for {
		exibeMenu()
		comando := leComando()

		switch comando {
		case 1:
			iniciarMonitoramento()

		case 2:
			fmt.Println("Exibindo logs...")
			imprimeLogs()
		case 0:
			fmt.Println("Saindo...")
			os.Exit(0)
		default:
			fmt.Println("Não conheço esse comando")
			os.Exit(-1)
		}
	}
}

func exibeIntroducao() {

	nome := "Arielson"
	versao := 1.1

	fmt.Println("Olá, sr.", nome)
	fmt.Println("Versão:", versao)
}

func exibeMenu() {

	fmt.Println("1 - iniciar monitoramento")
	fmt.Println("2 - ver logs")
	fmt.Println("0 - sair")
}

func leComando() int {

	var comandoLido int

	fmt.Scan(&comandoLido)

	return comandoLido
}

func iniciarMonitoramento() {

	fmt.Println("Monitorando...")

	sites := leSitesDoArquivo(fileName)

	for i := 0; i < monitoramento; i++ {

		for _, site := range sites {
			testaSite(site)
		}
		time.Sleep(time.Second * delay)
		fmt.Println("")
	}
}

func testaSite(site string) {

	resp, err := http.Get(site)

	if err != nil {
		fmt.Println("Erro ao acessar o site", site, "-", err)
	} else if resp.StatusCode == 200 {
		fmt.Println("Site", site, "foi carregado com sucesso")
		fmt.Println("Status code:", resp.StatusCode)
		registraLogs(site, true)
	} else {
		fmt.Println("Site", site, "está com problemas. Status code:", resp.StatusCode)
		registraLogs(site, false)
	}
}

func leSitesDoArquivo(fileName string) []string {
	var sites []string

	arquivo, err := os.Open(fileName)

	if err != nil {
		fmt.Println("Ocorreu um erro ao abrir o arquivo:", err)
	}

	leitor := bufio.NewReader(arquivo)

	for {
		linha, err := leitor.ReadString('\n')
		linha = strings.TrimSpace(linha)

		sites = append(sites, linha)

		if err == io.EOF {
			break
		}
	}

	arquivo.Close()

	return sites
}

func registraLogs(site string, status bool) {

	arquivo, err := os.OpenFile(logsPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		msgError(err)
	}

	arquivo.WriteString(time.Now().Format("01/02/2006 15:04:05") + " - " + site + " - online: " + strconv.FormatBool(status) + "\n")
	arquivo.Close()
}

func imprimeLogs() {

	arquivo, err := ioutil.ReadFile(logsPath)

	if err != nil {
		msgError(err)
	}

	fmt.Println(string(arquivo))
}

func msgError(err error) {

	fmt.Println("Ocorreu um erro ao abrir o arquivo:", err)
	os.Exit(-1)
}

func pegaArgs() (string, string) {

	fileName := flag.String("file", "sites.txt", "Arquivo de sites")
	logsPath := flag.String("logs", "logs.txt", "Arquivo de logs")

	flag.Parse()

	return *fileName, *logsPath
}
