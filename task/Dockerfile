FROM golang:alpine

ARG TASKPORT
ENV TASKPORT $TASKPORT
EXPOSE $TASKPORT

WORKDIR /app
COPY . . 
RUN go mod tidy
RUN go build -o main cmd/app/main.go
CMD [ "./main" ]