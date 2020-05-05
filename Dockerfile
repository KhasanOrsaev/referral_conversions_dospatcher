FROM golang:1.13
#RUN mkdir /go/src
ARG GIT_LOGIN

ADD . /app
WORKDIR /app

RUN git config --global url."https://git.fin-dev.ru/scm".insteadof "https://git.fin-dev.ru"
RUN echo "${GIT_LOGIN}" > ~/.netrc
RUN go env -w GONOSUMDB="git.fin-dev.ru" # && go build -o main ./cmd/dispatcher/main.go

#CMD ["./main"]
