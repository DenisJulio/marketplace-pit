CREATE TABLE usuarios (
    id SERIAL PRIMARY KEY,
    nome_de_usuario VARCHAR(255) NOT NULL UNIQUE,
    nome VARCHAR(255) NOT NULL,
    imagem VARCHAR(255)
);

CREATE TABLE anuncios (
    id SERIAL PRIMARY KEY,
    nome VARCHAR(255) NOT NULL,
    criado_em TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    anunciante_id INTEGER NOT NULL,
    valor DECIMAL(10, 2) NOT NULL,
    descricao TEXT,
    imagem VARCHAR(255),
    FOREIGN KEY (anunciante_id) REFERENCES usuarios (id)
);

CREATE TABLE ofertas (
    id SERIAL PRIMARY KEY,
    criado_em TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    ofertante_id INTEGER NOT NULL,
    anunciante_id INTEGER NOT NULL,
    anuncio_id INTEGER NOT NULL,
    FOREIGN KEY (ofertante_id) REFERENCES usuarios (id),
    FOREIGN KEY (anunciante_id) REFERENCES usuarios (id),
    FOREIGN KEY (anuncio_id) REFERENCES anuncios (id)
);

CREATE TABLE mensagens (
    id SERIAL PRIMARY KEY,
    criado_em TIMESTAMP NOT NULL,
    conteudo TEXT NOT NULL,
    oferta_id INTEGER NOT NULL,
    remetente_id INTEGER NOT NULL,
    FOREIGN KEY(oferta_id) REFERENCES ofertas (id),
    FOREIGN KEY(remetente_id) REFERENCES usuarios (id)
);

-- Insert into usuarios table
INSERT INTO usuarios (nome, nome_de_usuario, imagem) VALUES
('Pedro Santos', 'pedr0', '/resources/images/avatars/pedro-avatar.jpg'),
('Maria Antonia', 'mari4', '/resources/images/avatars/maria-avatar.jpg'),
('João Marcos', 'joa0', '/resources/images/avatars/joao-avatar.jpg');

-- Insert into anuncios table
INSERT INTO anuncios (nome, criado_em, anunciante_id, valor, descricao, imagem) VALUES
('Carro', '2024-05-28 15:00', 1, 15750.00, 'Veiculo usado mas bem conservado', '/resources/images/anuncios/anuncio-carro.jpg'),
('Moto', '2024-01-29 10:00', 1, 7000.00, 'Aceito troca por outra moto.', '/resources/images/anuncios/anuncio-moto.jpg'),
('Livro', '2024-10-20 07:00', 2, 50.00, 'Livro em perfeitas condições,apenas algumas marcas de realçe', '/resources/images/anuncios/anuncio-livro.jpg'),
('Tenis', '2024-06-15 20:00', 3, 150.00, 'Tenis com pouco uso', '/resources/images/anuncios/anuncio-tenis.webp'),
('Camisa', '2024-03-18 10:00', 3, 47.99, 'Camisa super estilosa', '/resources/images/anuncios/anuncio-camisa.avif');

-- Insert into ofertas table
INSERT INTO ofertas (criado_em, ofertante_id, anunciante_id, anuncio_id) VALUES
('2024-11-01 10:00', 3, 2, 3),
('2024-11-03 10:00', 3, 1, 1),
('2024-11-25 11:00', 2, 1, 2),
('2024-12-04 12:00', 1, 3, 4);

-- Insert into mensagens table
INSERT INTO mensagens (criado_em, conteudo, oferta_id, remetente_id) VALUES
('2024-11-01 11:00', 'Ola, gostaria de saber mais sobre o produto.', 1, 3),
('2024-11-03 11:00', 'Ola, gostaria de saber mais sobre o produto.', 2, 3),
('2024-11-25 12:00', 'Ola, gostaria de saber mais sobre o produto.', 3, 2),
('2024-12-04 13:00', 'Ola, gostaria de saber mais sobre o produto.', 4, 1);
