package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/elman23/articleapi/pkg/mocks"
	_ "github.com/lib/pq"
)

const (
	host = "localhost"
	port = 5432
)

func Connect(user string, password string, dbname string) *sql.DB {

	connInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", connInfo)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected to the database!")
	return db
}

func CloseConnection(db *sql.DB) {
	defer db.Close()
}

func CreateTable(db *sql.DB) {
	var exists bool
	if err := db.QueryRow("SELECT EXISTS (SELECT FROM pg_tables WHERE  schemaname = 'public' AND tablename = 'articles' );").Scan(&exists); err != nil {
		fmt.Println("Failed to execute query!", err)
		return
	}
	if !exists {
		results, err := db.Query("CREATE TABLE articles (id VARCHAR(36) PRIMARY KEY, title VARCHAR(100) NOT NULL, description VARCHAR(50) NOT NULL, content VARCHAR(50) NOT NULL);")
		if err != nil {
			fmt.Println("Failed to execute query!", err)
			return
		}
		fmt.Println("Table created successfully!", results)

		for _, article := range mocks.Articles {
			queryStmt := `INSERT INTO articles (id,title,description,content) VALUES ($1, $2, $3, $4) RETURNING id;`

			err := db.QueryRow(queryStmt, &article.Id, &article.Title, &article.Desc, &article.Content).Scan(&article.Id)
			if err != nil {
				log.Println("Failed to execute query!", err)
				return
			}
		}
		fmt.Println("Mock Articles included in Table!", results)
	} else {
		fmt.Println("Table 'articles' already exists!")
	}
}
