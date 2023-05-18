# Article REST API

## Links

REST API development: https://dev.to/janirefdez/create-a-rest-api-with-go-1j52

Database connection: https://dev.to/janirefdez/connect-rest-api-to-database-with-go-d8m

Create and connect to PostgreSQL: https://towardsdatascience.com/how-to-run-postgresql-and-pgadmin-using-docker-3a6a8ae918b5

## Table creation

```
CREATE TABLE articles (id VARCHAR(36) PRIMARY KEY, title VARCHAR(100) NOT NULL, description VARCHAR(50) NOT NULL, content VARCHAR(50) NOT NULL);

SELECT EXISTS (SELECT FROM pg_tables WHERE  schemaname = 'public' AND tablename = 'articles' );

CREATE TABLE IF NOT EXISTS  articles  (id VARCHAR(36) PRIMARY KEY, title VARCHAR(100) NOT NULL, description VARCHAR(50) NOT NULL, content VARCHAR(50) NOT NULL);

INSERT INTO articles (id,title,description,content) VALUES ('8617bf49-39a9-4268-b113-7b6bcd189ba2', 'Article 1', 'Article Description 1', 'Article Content 1');
```

## Environment variables

Create a `.env` file in the home of the project. The content is:
```
POSTGRES_USER=myuser
POSTGRES_PASSWORD=mypassword
POSTGRES_DB=articlesDB
PGADMIN_DEFAULT_EMAIL=user@domain.com
PGADMIN_DEFAULT_PASSWORD=password
```