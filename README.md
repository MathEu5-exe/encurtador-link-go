# 🔗 Encurtador de Links em Go (Golang)

Um projeto *Full-Stack* focado nos fundamentos do desenvolvimento backend, construído inteiramente com as bibliotecas padrão do Go, sem depender de frameworks externas para o roteamento.

## Sobre o Projeto

O objetivo principal deste projeto foi construir uma API REST funcional do zero, compreendendo o fluxo completo de uma aplicação web: desde a receção de um pedido HTTP, processamento lógico, base de dados e devolução da resposta para uma interface de utilizador.

## Principais tecnologias Usadas

*   **Backend:** Go (Golang)
*   **Base de Dados:** SQLite (via `modernc.org/sqlite` para suporte multiplataforma sem dependência de CGO).
*   **Frontend:** HTML5, JavaScript e Tailwind CSS.
*   **Arquitetura:** API REST (JSON).

## O Que aprendi com este projeto?

Este projeto foi um excelente laboratório para consolidar conceitos essenciais da engenharia de software:

*   **Fundamentos de Go:** Criação de servidores HTTP, manipulação de rotas nativas (`http.HandleFunc`) e concorrência básica da linguagem.
*   **SQL e CRUD na Prática:** Integração com SQLite, compreendendo como os comandos `INSERT` guardam a informação gerada e os `SELECT` recuperam a URL original.
*   **Integração de Sistemas:** Como o frontend (HTML/JS) comunica dinamicamente com o backend (Go) através de requisições `POST` (para criar o link) e `GET` (para redirecionar).
*   **Papel do Desenvolvedor:** Compreendi que a engenharia de software vai muito além de "decorar e digitar código do zero". Trata-se de arquitetar soluções, integrar ferramentas existentes (como o Tailwind para a UI) e focar na lógica de negócios que realmente resolve o problema.

## Como executar localmente?

1. Clone o repositório:
   ```bash
   git clone [https://github.com/SEU-USUARIO/encurtador-links-go.git](https://github.com/SEU-USUARIO/encurtador-links-go.git)


