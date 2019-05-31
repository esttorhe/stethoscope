FROM golang:1.12 AS builder

ENV APP_NAME=stethoscope
WORKDIR /opt/app

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -o stethoscope
RUN mkdir -p /opt/built/ && \
   mv stethoscope /opt/built && \
   mv *.yml /opt/built

FROM esttorhe/ca-alpine:latest

# Copy the certificates
COPY --from=builder /etc/ssl /etc/ssl/

# Copy the application files
WORKDIR /opt/app

COPY --from=builder /opt/built .
ENV PORT=7000
ENV LOG=INFO
EXPOSE 7000

ENTRYPOINT /opt/app/stethoscope ${LOG}