CREATE DATABASE IF NOT EXISTS playground;

USE playground;

CREATE TABLE IF NOT EXISTs users (
	id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
	name VARCHAR(255) NOT NULL,
	email VARCHAR(255) NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT uc_user UNIQUE (name, email)
);

