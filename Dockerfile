# Define a imagem base (neste caso, uma imagem oficial do Go com a versão 1.20)
FROM golang:1.20

ENTRYPOINT ["tail", "-f", "/dev/null"]
# # Define o diretório de trabalho dentro do contêiner
# WORKDIR /app

# # Copia o código fonte para o diretório de trabalho
# COPY . .

# # Define a variável de ambiente para o projeto
# ENV GO111MODULE=on

# # Executa o comando de build do projeto Go
# RUN go build -o main .

# # Define o comando para executar o binário
# CMD ["./main"]