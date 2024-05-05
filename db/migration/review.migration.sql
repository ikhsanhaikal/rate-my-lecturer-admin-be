CREATE TABLE reviews (
	id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
	classId BIGINT NOT NULL,
	reviewerId BIGINT NOT NULL,
	comment VARCHAR(255),
	rating FLOAT NOT NULL, 

	CONSTRAINT fk_class FOREIGN KEY (classId)
											REFERENCES classes(id)
											ON UPDATE CASCADE
											ON DELETE CASCADE,

	CONSTRAINT fk_user FOREIGN KEY (reviewerId)
											REFERENCES users(id)
											ON UPDATE CASCADE
											ON DELETE CASCADE
);