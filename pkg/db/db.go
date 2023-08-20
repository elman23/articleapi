package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/elman23/articleapi/pkg/mocks"
	_ "github.com/lib/pq"
)

type DbServiceSingleton struct {
	dbService *DbService
}

func (s *DbServiceSingleton) GetService() *DbService {
	if s.dbService == nil {
		s.dbService = newDbService()
	}
	return s.dbService
}

type DbService struct {
	DB *sql.DB
}

func newDbService() *DbService {

	dbService := DbService{}
	// Retreive database connection information from enviroment variables
	url, port, user, password, dbname :=
		os.Getenv("POSTGRES_URL"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB")

	// Connect to the database
	dbService.DB = dbService.connect(url, port, user, password, dbname)

	// Create the necessary tables in the database
	dbService.createTables(dbService.DB)

	return &dbService
}

func (s *DbServiceSingleton) CloseConnection() {
	// Get the service
	singleton := s.GetService()
	// Close database connection
	defer singleton.DB.Close()
}

func (s *DbService) connect(host string, port string, user string, password string, dbname string) *sql.DB {

	connInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", connInfo)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	log.Println("Successfully connected to the database!")
	return db
}

func (s *DbService) createTables(db *sql.DB) {
	s.createUsersTable(db)
	s.createArticlesTable(db)
}

func (s *DbService) createArticlesTable(db *sql.DB) {
	var exists bool
	if err := db.QueryRow("SELECT EXISTS (SELECT FROM pg_tables WHERE  schemaname = 'public' AND tablename = 'articles' );").Scan(&exists); err != nil {
		log.Println("Failed to execute query!", err)
		return
	}
	if !exists {
		// Create the 'articles' table.
		results, err := db.Query("CREATE TABLE articles (id VARCHAR(36) PRIMARY KEY, title VARCHAR(100) NOT NULL, description VARCHAR(50) NOT NULL, content VARCHAR(50) NOT NULL);")
		if err != nil {
			log.Println("Failed to create 'articles' table!", err)
			return
		}
		log.Println("Table 'articles' created successfully!", results)

		for _, article := range mocks.Articles {
			queryStmt := `INSERT INTO articles (id,title,description,content) VALUES ($1, $2, $3, $4) RETURNING id;`

			err := db.QueryRow(queryStmt, &article.Id, &article.Title, &article.Desc, &article.Content).Scan(&article.Id)
			if err != nil {
				log.Println("Failed insert mock articles in table!", err)
				return
			}
		}
		log.Println("Mock articles included in table!", results)
	} else {
		log.Println("Table 'articles' already exists!")
	}
}

func (s *DbService) createUsersTable(db *sql.DB) {
	var exists bool
	if err := db.QueryRow("SELECT EXISTS (SELECT FROM pg_tables WHERE  schemaname = 'public' AND tablename = 'users' );").Scan(&exists); err != nil {
		log.Println("Failed to execute query!", err)
		return
	}
	if !exists {
		// Create the 'users' table.
		results, err := db.Query("CREATE TABLE users (id VARCHAR(36) PRIMARY KEY, username VARCHAR(50) NOT NULL, password VARCHAR(50) NOT NULL);")
		if err != nil {
			log.Println("Failed to create 'users' table!", err)
			return
		}
		log.Println("Table 'users' created successfully!", results)

		for _, user := range mocks.Users {
			queryStmt := `INSERT INTO users (id,username,password) VALUES ($1, $2, $3) RETURNING id;`

			err := db.QueryRow(queryStmt, &user.Id, &user.Username, &user.Password).Scan(&user.Id)
			if err != nil {
				log.Println("Failed to insert default user in table!", err)
				return
			}
		}
		log.Println("Default users included in table!", results)
	} else {
		log.Println("Table 'users' already exists!")
	}
}
