
FROM golang:1.18
RUN mkdir -p /app
COPY . /app
WORKDIR /app/
RUN go build && mv ntpChecker /
RUN ls
WORKDIR /app
RUN rm -rf *
ENTRYPOINT ["/ntpChecker"]