package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"

	// O sublinhado (_) diz ao Go: "Importa isto só para iniciar o driver,
	// mesmo que eu não chame nenhuma função dele diretamente"
	_ "modernc.org/sqlite"
)

type Requisicao struct {
	URL string `json:"url"` // essa variável mapeia a chave "url" do JSON
}

// db é a nossa ligação global à base de dados.
// Usar uma variável global ajuda a evitar abrir e fechar a base de dados a cada requisição.
var db *sql.DB

func initDB() {
	var err error
	// Abre a ligação à base de dados, e se não tiver, o próprio Go cria uma base dados vazia
	db, err = sql.Open("sqlite", "./encurtador.db")
	if err != nil {
		log.Fatal("Erro ao abrir base de dados: ", err)
	}

	// Código SQL puro para criar a nossa tabela, caso ela não exista
	query := `
	CREATE TABLE IF NOT EXISTS urls (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		url_original TEXT NOT NULL,
		codigo_curto TEXT NOT NULL UNIQUE
	);`

	// Executa a query
	_, err = db.Exec(query)
	if err != nil {
		log.Fatal("Erro ao criar a tabela: ", err)
	}

	fmt.Println("Base de dados iniciada com sucesso!")
}

const caracteres = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func gerarCodigo() string {
	// Cria um "array" vazio de 6 espaços para guardar nossas letras
	codigo := make([]byte, 6)

	// Faz um loop passando por cada um dos 6 espaços
	for i := range codigo {
		// Sorteia uma posição aleatória da nossa string 'caracteres'
		posicaoSorteada := rand.Intn(len(caracteres))
		// Guarda a letra sorteada no espaço atual do array
		codigo[i] = caracteres[posicaoSorteada]
	}

	// Converte o array de letras para uma String normal e devolve
	return string(codigo)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Servidor do Encurtador a correr com sucesso!")
}

func handleEncurtar(w http.ResponseWriter, r *http.Request) {
	// 1. Barreira de segurança: Só aceitamos o método POST (envio de dados)
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido. Use POST.", http.StatusMethodNotAllowed)
		return
	}

	// 2. Lemos o JSON que o usuário enviou e colocamos na nossa Struct
	var req Requisicao
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.URL == "" {
		http.Error(w, "URL inválida ou não fornecida", http.StatusBadRequest)
		return
	}

	// 3. Geramos nosso código de 6 letras
	codigo := gerarCodigo()

	// 4. Inserimos no banco de dados SQLite
	// Os sinais de interrogação (?) previnem SQL Injection (ataques de injeção de SQL)
	_, err = db.Exec("INSERT INTO urls (url_original, codigo_curto) VALUES (?, ?)", req.URL, codigo)
	if err != nil {
		http.Error(w, "Erro ao salvar no banco de dados", http.StatusInternalServerError)
		return
	}

	// 5. Montamos o link encurtado e devolvemos ao usuário
	linkCurto := "http://localhost:8080/" + codigo

	w.Header().Set("Content-Type", "application/json") // Avisa que a resposta é um JSON
	w.WriteHeader(http.StatusCreated)                  // Retorna HTTP 201 (Criado)

	// Escreve a resposta final na tela
	fmt.Fprintf(w, `{"url_encurtada": "%s"}`, linkCurto)
}

func handleRedirect(w http.ResponseWriter, r *http.Request) {
	// 1. Extrair o código do caminho (se a visita for em /G1pmGY, guardamos apenas "G1pmGY")
	// r.URL.Path[1:] ignora o primeiro caractere, que é a barra (/) usando aritmética de strings
	codigo := r.URL.Path[1:]

	// Se o utilizador aceder apenas a localhost:8080/, mostramos uma mensagem simples
	if codigo == "" {
		http.ServeFile(w, r, "index.html")
		return
	}

	// 2. Fazer o SELECT na base de dados
	var urlOriginal string
	// Pedimos ao Go para executar a query e "injetar" (Scan) o resultado na variável urlOriginal
	err := db.QueryRow("SELECT url_original FROM urls WHERE codigo_curto = ?", codigo).Scan(&urlOriginal)

	// Lidar com o caso do link não existir na base de dados ou se houver outro erro
	if err == sql.ErrNoRows {
		http.Error(w, "Link não encontrado na base de dados", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Erro interno ao procurar o link", http.StatusInternalServerError)
		return
	}

	// Se encontrou a URL, pedimos para ele redirecionar a ela.
	http.Redirect(w, r, urlOriginal, http.StatusFound)
}

func main() {
	// 1. Prepara a base de dados antes de iniciar o servidor
	initDB()

	// 2. Garante que a base de dados é fechada quando o servidor for desligado
	defer db.Close()

	// Usada apenas para testar se o servidor está online
	http.HandleFunc("/ping", handleHealth)
	// Recebe o link via POST (apenas recebe) e gera o código curto
	http.HandleFunc("/encurtar", handleEncurtar)
	// Rota "Curinga" redireciona pra URL original caso nenhuma das outras funcione.
	http.HandleFunc("/", handleRedirect)

	fmt.Println("Servidor iniciado na porta 8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Erro crítico no servidor: ", err)
	}
}
