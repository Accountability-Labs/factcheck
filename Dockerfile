FROM golang:1.21

WORKDIR /src/
COPY . .
RUN make
EXPOSE 8080

CMD ["/src/factcheck"]
