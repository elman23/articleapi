CREATE TABLE articles (id VARCHAR(36) PRIMARY KEY, title VARCHAR(100) NOT NULL, description VARCHAR(50) NOT NULL, content VARCHAR(50) NOT NULL);

SELECT EXISTS (SELECT FROM pg_tables WHERE  schemaname = 'public' AND tablename = 'articles' );

CREATE TABLE IF NOT EXISTS  articles  (id VARCHAR(36) PRIMARY KEY, title VARCHAR(100) NOT NULL, description VARCHAR(50) NOT NULL, content VARCHAR(50) NOT NULL);

INSERT INTO articles (id,title,description,content) VALUES ('8617bf49-39a9-4268-b113-7b6bcd189ba2', 'Article 1', 'Article Description 1', 'Article Content 1');