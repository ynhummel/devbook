package repositorios

import (
	"api/src/models"
	"database/sql"
	"fmt"
)

type usuarios struct {
	db *sql.DB
}

// NovoRepositorioDeUsuarios cria um repositório de usuarios
func NovoRepositorioDeUsuarios(db *sql.DB) *usuarios {
	return &usuarios{db}
}

// Criar insere um usuário no banco de dados
func (u usuarios) Criar(usuario models.Usuario) (uint64, error) {
	statement, err := u.db.Prepare(
		"insert into usuarios (nome, nick, email, senha) values(?, ?, ?, ?)",
	)
	if err != nil {
		return 0, err
	}
	defer statement.Close()

	resultado, err := statement.Exec(usuario.Nome, usuario.Nick, usuario.Email, usuario.Senha)
	if err != nil {
		return 0, err
	}

	ultimoIDInserido, err := resultado.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint64(ultimoIDInserido), nil
}

// Buscar traz todos os usuarios que possuam nomeOuNick
func (u usuarios) Buscar(nomeOuNick string) ([]models.Usuario, error) {
	nomeOuNick = fmt.Sprintf("%%%s%%", nomeOuNick)

	linhas, err := u.db.Query(
		"select id, nome, nick, email, criadoEm from usuarios where nome LIKE ? or nick LIKE ?",
		nomeOuNick, nomeOuNick,
	)
	if err != nil {
		return nil, err
	}
	defer linhas.Close()

	var usuarios []models.Usuario

	for linhas.Next() {
		var usuario models.Usuario

		if err = linhas.Scan(
			&usuario.ID,
			&usuario.Nome,
			&usuario.Nick,
			&usuario.Email,
			&usuario.CriadoEm,
		); err != nil {
			return nil, err
		}

		usuarios = append(usuarios, usuario)
	}
	return usuarios, nil
}

// BuscarPorID traz um usuário do banco de dados
func (u usuarios) BuscarPorID(ID uint64) (models.Usuario, error) {
	linhas, err := u.db.Query(
		"select id, nome, nick, email, criadoEm from usuarios where id = ?", ID,
	)
	if err != nil {
		return models.Usuario{}, err
	}
	defer linhas.Close()

	var usuario models.Usuario

	if linhas.Next() {
		if err = linhas.Scan(
			&usuario.ID,
			&usuario.Nome,
			&usuario.Nick,
			&usuario.Email,
			&usuario.CriadoEm,
		); err != nil {
			return models.Usuario{}, err
		}
	}
	return usuario, nil
}

// BuscarPorEmail traz um usuario do banco de dados que possua o email passado
func (u usuarios) BuscarPorEmail(email string) (models.Usuario, error) {
	linhas, err := u.db.Query(
		"select id, senha from usuarios where email = ?", email,
	)
	if err != nil {
		return models.Usuario{}, err
	}
	defer linhas.Close()

	var usuario models.Usuario

	if linhas.Next() {
		if err = linhas.Scan(&usuario.ID, &usuario.Senha); err != nil {
			return models.Usuario{}, err
		}
	}

	return usuario, nil
}

// Altualizar altera as informações de um usuario no banco de dados
func (u usuarios) Atualizar(ID uint64, usuario models.Usuario) error {
	statement, err := u.db.Prepare(
		"update usuarios set nome = ?, nick = ?, email = ? where id = ?",
	)
	if err != nil {
		return err
	}
	defer statement.Close()

	if _, err = statement.Exec(usuario.Nome, usuario.Nick, usuario.Email, ID); err != nil {
		return err
	}
	return nil
}

// Exclui as informações de um usuário no banco de dados
func (u usuarios) Deletar(ID uint64) error {
	statement, err := u.db.Prepare("delete from usuarios where id = ?")
	if err != nil {
		return err
	}
	defer statement.Close()

	if _, err = statement.Exec(ID); err != nil {
		return err
	}
	return nil
}
