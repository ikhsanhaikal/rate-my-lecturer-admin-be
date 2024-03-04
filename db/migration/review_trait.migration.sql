CREATE TABLE reviews_and_traits (
	id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,

	reviewId BIGINT NOT NULL,
	traitId BIGINT NOT NULL,

	CONSTRAINT fk_review FOREIGN KEY (reviewId)
											REFERENCES reviews(id)
											ON UPDATE CASCADE
											ON DELETE CASCADE,

	CONSTRAINT fk_trait FOREIGN KEY (traitId)
											REFERENCES traits(id)
											ON UPDATE CASCADE
											ON DELETE CASCADE
);