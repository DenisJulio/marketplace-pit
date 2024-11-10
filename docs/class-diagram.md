# Class Diagram

@startuml class_diagram

class Anuncio {
    Long id
    String nome
    LocalDateTime criadoEm
    Long anuncianteId
    Double valor
    String descricao
    String imagem
}

class Usuario {
    String email
    Long id
    String nome
}

class Oferta {
    Long id
    LocalDateTime criadoEm
    Long ofertanteId
    Long anuncianteId
    Long anuncioId
}

class Mensagem {
    Long id
    String conteudo
    LocalDateTime criadoEm
    Long ofertaId
    Long remetenteId
}

Anuncio "1" --o "1" Usuario : anunciante
Oferta "0..*" --o "1" Anuncio : anuncio
Oferta "1" --o "1" Usuario : ofertante
Oferta "1" --o "1" Usuario : anunciante
Mensagem "1..*" --* "1" Oferta : oferta

@enduml
