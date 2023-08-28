CREATE TABLE Obrigations (
    `id` INT NOT NULL AUTO_INCREMENT, 
    `name` varchar(100) NOT NULL,
    `mandatory` BOOLEAN default(false),
    `qr_code` varchar(100) NOT NULL,
	`create_at` TIMESTAMP NOT NULL DEFAULT '0000-00-00 00:00:00',
  	`update_at` TIMESTAMP NOT NULL ON UPDATE CURRENT_TIMESTAMP,
  	PRIMARY KEY(id)
)