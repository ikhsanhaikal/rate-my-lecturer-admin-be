CREATE TABLE reviews (
	id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
	classId BIGINT NOT NULL,
	userId BIGINT NOT NULL,
	comments VARCHAR(255),
	rating FLOAT NOT NULL, 

	CONSTRAINT fk_class FOREIGN KEY (classId)
											REFERENCES classes(id)
											ON UPDATE CASCADE
											ON DELETE CASCADE,

	CONSTRAINT fk_user FOREIGN KEY (userId)
											REFERENCES users(id)
											ON UPDATE CASCADE
											ON DELETE CASCADE
);