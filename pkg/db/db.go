package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/elman23/articleapi/pkg/mocks"
	_ "github.com/lib/pq"
)

func Connect(host string, port string, user string, password string, dbname string) *sql.DB {

	connInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

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

func CreateTables(db *sql.DB) {
	createUsersTable(db)
	createArticlesTable(db)
}

func createArticlesTable(db *sql.DB) {
	var exists bool
	if err := db.QueryRow("SELECT EXISTS (SELECT FROM pg_tables WHERE  schemaname = 'public' AND tablename = 'articles' );").Scan(&exists); err != nil {
		fmt.Println("Failed to execute query!", err)
		return
	}
	if !exists {
		// Create the 'articles' table.
		results, err := db.Query("CREATE TABLE articles (id VARCHAR(36) PRIMARY KEY, title VARCHAR(100) NOT NULL, description VARCHAR(50) NOT NULL, content VARCHAR(50) NOT NULL);")
		if err != nil {
			fmt.Println("Failed to create 'articles' table!", err)
			return
		}
		fmt.Println("Table 'articles' created successfully!", results)

		for _, article := range mocks.Articles {
			queryStmt := `INSERT INTO articles (id,title,description,content) VALUES ($1, $2, $3, $4) RETURNING id;`

			err := db.QueryRow(queryStmt, &article.Id, &article.Title, &article.Desc, &article.Content).Scan(&article.Id)
			if err != nil {
				log.Println("Failed insert mock articles in table!", err)
				return
			}
		}
		fmt.Println("Mock articles included in table!", results)
	} else {
		fmt.Println("Table 'articles' already exists!")
	}
}

func createUsersTable(db *sql.DB) {
	var exists bool
	if err := db.QueryRow("SELECT EXISTS (SELECT FROM pg_tables WHERE  schemaname = 'public' AND tablename = 'users' );").Scan(&exists); err != nil {
		fmt.Println("Failed to execute query!", err)
		return
	}
	if !exists {
		// Create the 'users' table.
		results, err := db.Query("CREATE TABLE users (id VARCHAR(36) PRIMARY KEY, username VARCHAR(50) NOT NULL, password VARCHAR(50) NOT NULL);")
		if err != nil {
			fmt.Println("Failed to create 'users' table!", err)
			return
		}
		fmt.Println("Table 'users' created successfully!", results)

		for _, user := range mocks.Users {
			queryStmt := `INSERT INTO users (id,username,password) VALUES ($1, $2, $3) RETURNING id;`

			err := db.QueryRow(queryStmt, &user.Id, &user.Username, &user.Password).Scan(&user.Id)
			if err != nil {
				log.Println("Failed to insert default user in table!", err)
				return
			}
		}
		fmt.Println("Default users included in table!", results)
	} else {
		fmt.Println("Table 'users' already exists!")
	}
}
