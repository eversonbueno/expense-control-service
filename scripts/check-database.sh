#!/bin/bash

# Script para verificar se as tabelas foram criadas corretamente no MySQL

echo "🔍 Verificando se as tabelas foram criadas corretamente..."
echo "=================================================="

# Verificar se o container MySQL está rodando
if ! docker ps | grep -q "finance_db"; then
    echo "❌ Container MySQL não está rodando!"
    echo "Execute: make up"
    exit 1
fi

echo "✅ Container MySQL está rodando"
echo ""

# Conectar no MySQL e verificar as tabelas
echo "📋 Listando todas as tabelas:"
docker exec finance_db mysql -u everson.bueno -pTeste@1234 expense-control-service -e "SHOW TABLES;"

echo ""
echo "📊 Verificando estrutura das tabelas:"
echo ""

# Verificar estrutura de cada tabela
tables=("usuarios" "contas" "categorias" "lancamentos" "orcamentos")

for table in "${tables[@]}"; do
    echo "🔍 Estrutura da tabela: $table"
    docker exec finance_db mysql -u everson.bueno -pTeste@1234 expense-control-service -e "DESCRIBE $table;"
    echo ""
done

echo "✅ Verificação concluída!"
echo ""
echo "💡 Para conectar manualmente no MySQL:"
echo "docker exec -it finance_db mysql -u everson.bueno -pTeste@1234 expense-control-service"
