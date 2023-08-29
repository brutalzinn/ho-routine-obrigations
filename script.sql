CREATE TABLE Obrigations (
    `id` BINARY(16) DEFAULT UUID(),
    `name` varchar(100) NOT NULL,
    `mandatory` BOOLEAN default(false),
    `qr_code` varchar(100) NOT NULL,
	`create_at` TIMESTAMP NOT NULL DEFAULT '0000-00-00 00:00:00',
  	`update_at` TIMESTAMP NOT NULL ON UPDATE CURRENT_TIMESTAMP,
  	PRIMARY KEY(id)
)

INSERT INTO Obrigations (name,mandatory,qr_code) values ("cozinha", true, "hora-do-almoco")
INSERT INTO Obrigations (name,mandatory,qr_code) values ("sala", true, "hora-da-pausa")