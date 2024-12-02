CREATE TABLE usuarios (
    id SERIAL PRIMARY KEY,
    nome_de_usuario VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL,
    nome VARCHAR(255) NOT NULL,
    senha VARCHAR(255) NOT NULL,
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

CREATE TABLE sessoes (
    sessao_id VARCHAR PRIMARY KEY,
    nome_de_usuario VARCHAR NOT NULL,
    expira_em TIMESTAMP NOT NULL
);

-- Insert into usuarios table
INSERT INTO
    usuarios (
        nome,
        nome_de_usuario,
        email,
        senha,
        imagem
    )
VALUES (
        'Pedro Santos',
        'pedr0',
        'pedro@email.com',
        'pedro santos',
        '/resources/images/avatars/pedro-avatar.jpg'
    ),
    (
        'Maria Antonia',
        'mari4',
        'maria@email.com',
        'maria antonia',
        '/resources/images/avatars/maria-avatar.jpg'
    ),
    (
        'João Marcos',
        'joa0',
        'joao@email.com',
        'joao marcos',
        '/resources/images/avatars/joao-avatar.jpg'
    );

-- Insert into anuncios table
INSERT INTO
    anuncios (
        nome,
        criado_em,
        anunciante_id,
        valor,
        descricao,
        imagem
    )
VALUES (
        'Carro',
        '2024-05-28 15:00',
        1,
        15750.00,
        'Veiculo usado mas bem conservado',
        '/resources/images/anuncios/anuncio-carro.jpg'
    ),
    (
        'Moto',
        '2024-01-29 10:00',
        1,
        7000.00,
        'Aceito troca por outra moto.',
        '/resources/images/anuncios/anuncio-moto.jpg'
    ),
    (
        'Livro',
        '2024-10-20 07:00',
        2,
        50.00,
        'Livro em perfeitas condições,apenas algumas marcas de realçe',
        '/resources/images/anuncios/anuncio-livro.jpg'
    ),
    (
        'Tenis',
        '2024-06-15 20:00',
        3,
        150.00,
        'Tenis com pouco uso',
        '/resources/images/anuncios/anuncio-tenis.webp'
    ),
    (
        'Camisa',
        '2024-03-18 10:00',
        3,
        47.99,
        'Camisa super estilosa',
        '/resources/images/anuncios/anuncio-camisa.avif'
    );