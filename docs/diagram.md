# Diagram

@startuml diagram

hide footbox

actor Usuário as u
actor Vendedor as v
participant "Pagina de lista\n anuncios" as pp
participant "Pagine de detalhes\n do anuncio" as dp
participant "Pagina de\n negociação" as pn
participant "Pagina de anuncios\n do vendedor" as pa

u -> pp : procura por **produtos**
u -> pp : seleciona **anuncio** para **detalhes**
pp -> dp : obtem acesso a\n informações do produto
u -> dp : **favorita** produto
u -> pn : abre uma **negociação** com o vendedor\n sobre o interesse em adquirir o produto
activate pn
pn -> v : notifica vendedor sobre negociação aberta
u <-> v : negociam sobre valores, entrega, etc.

alt negociação aceita

    v -> pn : sinaliza negociação como aceita/concluída
    pn -> pa : sinaliza **anuncio** como __concluído__
    pn --> u : retorna sobre confirmação da negociação
    |||
else necogiação rejeitada

    v -> pn : sinaliza negociação como rejeitada
    pn --> u : retorna sobre negativa da negociação
    deactivate pn

end

@enduml
