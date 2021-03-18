FROM golang:1.16-alpine

LABEL author="https://github.com/yolong-lin"

WORKDIR /code
COPY ./main.go ./main.go

# RUN cd /code
RUN go mod init example
RUN go mod tidy

EXPOSE 8080

CMD [ "go", "run", "." ]