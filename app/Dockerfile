FROM golang:alpine
RUN mkdir /code
ADD . /code/
WORKDIR /code

ARG ORBITAL_DB_ADDRESS
ENV ORBITAL_DB_ADDRESS=${ORBITAL_DB_ADDRESS}

RUN go build -o main .

CMD ["./main"]
