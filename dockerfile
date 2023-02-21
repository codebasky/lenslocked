FROM golang:alpine
WORKDIR /app
COPY . .
RUN go build -o lenslocked ./cmd

FROM alpine
RUN addgroup -g 1000 dudes \
    && adduser -u 1000 -G dudes -D dude
COPY --from=0 --chown=dude:dudes /app/lenslocked /lenslocked
EXPOSE 3000
USER dude
CMD ["/lenslocked"]
