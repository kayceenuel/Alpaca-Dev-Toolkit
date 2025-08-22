# --- stage 1: Build go app --- 
FROM golang:1.23.4-alpine AS builder

#Install git (needed for Go modules)
RUN apk add --no-cache git 

#Set working directory 
WORKDIR /app 

# Copy go mod files and download dependencies first (for better caching) 
COPY go.mod go.sum ./
RUN go.mod download 

# copy the rest of the code 
COPY . .

#Build the go binary 
RUN go build -o myapp .