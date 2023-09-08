CREATE TABLE Obrigations (
    `id` INT NOT NULL AUTO_INCREMENT, 
    `name` varchar(100) NOT NULL,
    `name_normalize` varchar(100) NOT NULL,
    `mandatory` BOOLEAN default(false),
    `qr_code` varchar(100) NOT NULL,
	`create_at` TIMESTAMP NOT NULL DEFAULT '0000-00-00 00:00:00',
  	`update_at` TIMESTAMP NOT NULL ON UPDATE CURRENT_TIMESTAMP,
  	PRIMARY KEY(id)
);

CREATE TABLE Devices (
    `id` INT NOT NULL AUTO_INCREMENT, 
    `name` varchar(100) NOT NULL,
    `token_firebase` TEXT NOT NULL,
	`create_at` TIMESTAMP NOT NULL DEFAULT '0000-00-00 00:00:00',
  	`update_at` TIMESTAMP NOT NULL ON UPDATE CURRENT_TIMESTAMP,
  	PRIMARY KEY(id)
);


CREATE TABLE Pending (
    `id` INT NOT NULL AUTO_INCREMENT, 
    `id_device` INT NOT NULL,
    `id_obrigation` INT NOT NULL,
    `waiting` BOOLEAN default(true),
	`expire_at` TIMESTAMP NOT NULL DEFAULT '0000-00-00 00:00:00',
	`create_at` TIMESTAMP NOT NULL DEFAULT '0000-00-00 00:00:00',
  	`update_at` TIMESTAMP NOT NULL ON UPDATE CURRENT_TIMESTAMP,
  	PRIMARY KEY(id),
    FOREIGN KEY (id_obrigation) REFERENCES Obrigations(id),
    FOREIGN KEY (id_device) REFERENCES Devices(id)
);

INSERT INTO Obrigations (name,name_normalize,mandatory,qr_code) values ("Comparecer a cozinha para tomar o medicamento da manha.","cozinha", true, "hora_do_medicamento");
INSERT INTO Obrigations (name,name_normalize,mandatory,qr_code) values ("Comparecer a cozinha para almocar","cozinha", true, "hora_do_almoco");
INSERT INTO Obrigations (name,name_normalize, mandatory,qr_code) values ("Comparecer a consulta médica.","consulta_medica", true, "hora_da_consulta_medica");
INSERT INTO Obrigations (name,name_normalize, mandatory,qr_code) values ("Pausa para assistir casimiro","pausa_assistir_casimiro", false, "hora_de_ver_casimiro");