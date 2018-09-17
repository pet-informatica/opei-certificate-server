package main

import (
	"os"
    "encoding/json"
	"database/sql"
	"log"
	"fmt"
	"net/http"
	"github.com/jung-kurt/gofpdf"
	_ "github.com/mattn/go-sqlite3"
)

// Dados do estudante que serão utilizados no PDF
type Student struct {
	Cpf string
	Name string
}

// Conexão é mantida aberta
var db *sql.DB

// Abre a conexão com o banco de dados
func setupDB() {
	localDb, err := sql.Open("sqlite3", "./opei.db")
	if err != nil {
		log.Fatal(err)
	}

	db = localDb
}

// Dado o CPF do estudante, preenche os dados da estrutura
func injectInfo(st *Student) {
	stmt, err := db.Prepare("select name from tbl where cpf = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	err = stmt.QueryRow(st.Cpf).Scan(&st.Name)
	if err != nil {
		log.Fatal(err)
	}
}

// Estrutura do PDF
func createCertificate(st Student) {
	pdf := gofpdf.New(gofpdf.OrientationPortrait, "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, st.Name)
	pdf.OutputFileAndClose(st.Cpf + ".pdf")
}

// Enviar um certificado
func serveCertificate(st Student, rw *http.ResponseWriter, req *http.Request) {
	// Preenche as informaçoes do estudante
    injectInfo(&st)

	// Cria o arquivo do PDF do certificado
    createCertificate(st)

	// Envia o arquivo para download
    http.ServeFile(*rw, req, st.Cpf + ".pdf")

	// Deleta o arquivo do server
    os.Remove(st.Cpf + ".pdf")
}

// Receber as requisições GET e POST
func webserver(rw http.ResponseWriter, req *http.Request) {
    if req.URL.Path != "/" {
        http.Error(rw, "404 not found.", http.StatusNotFound)
        return
    }

    switch req.Method {
    case "GET":
    	var st Student

	    // Extrai o CPF do parâmetro da URL
        st.Cpf = req.URL.Query()["cpf"][0]

        // Envia o certificado
        serveCertificate(st, &rw, req)
    case "POST":
    	var st Student

    	// Extrai o CPF do JSON do POST
    	decoder := json.NewDecoder(req.Body)
        err := decoder.Decode(&st)
        if err != nil {
        	log.Fatal(err)
        }

        // Envia o certificado
        serveCertificate(st, &rw, req)
    default:
        fmt.Fprintf(rw, "Sorry, only GET and POST methods are supported.")
    }
}

func main() {
	setupDB()

    http.HandleFunc("/", webserver)

    fmt.Printf("Starting server for testing HTTP GET and POST...\n")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal(err)
    }

    db.Close()
}
