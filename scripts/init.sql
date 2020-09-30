######################################################
# Creating tables

CREATE TABLE accounts (
    id int PRIMARY KEY UNIQUE AUTO_INCREMENT,
    document_number VARCHAR(11) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE operations (
    id int PRIMARY KEY UNIQUE,
    description VARCHAR(50) NOT NULL
);

CREATE TABLE transactions (
    id int PRIMARY KEY UNIQUE AUTO_INCREMENT,
    account_id int NOT NULL,
    operation_id int NOT NULL,
    amount DOUBLE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (account_id) REFERENCES accounts(id),
    FOREIGN KEY (operation_id) REFERENCES operations(id)
);

######################################################
# Inserting operations values

INSERT INTO `operations` (`id`, `description`) VALUES (1, 'COMPRA A VISTA');
INSERT INTO `operations` (`id`, `description`) VALUES (2, 'COMPRA PARCELADA');
INSERT INTO `operations` (`id`, `description`) VALUES (3, 'SAQUE');
INSERT INTO `operations` (`id`, `description`) VALUES (4, 'PAGAMENTO');