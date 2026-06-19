FROM golang:alpine

ARG PORT
ENV PORT $PORT
EXPOSE $PORT

WORKDIR /app
COPY . . 
RUN go mod tidy
RUN go build -o main cmd/app/main.go
CMD [ "./main" ]