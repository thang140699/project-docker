# Build 

FROM golang:1.18.1 AS build
WORKDIR /authen-docker 
COPY go.mod go.sum ./
RUN go mod download 
COPY . . 
# COPY authen-docker/controller/*.go ./
RUN go build -o authen /authen-docker/controller/*.go
WORKDIR /authen-docker
CMD ["./controller"]

#  Deploy
FROM ubuntu:22.04 AS runtime

RUN mkdir /app-authen
WORKDIR /app-authen
COPY --from=build /authen-docker/authen .
RUN chmod +x ./authen
ENV PORT 8081
CMD ./authen --configPrefix=wedding --configSource=/opt/event/authenticator/.env






# FROM ubuntu:22.04 AS runtime
# # RUN apt updated
# WORKDIR /opt/event/authenticator

# COPY ./controller/config  .

# COPY ./controller/.env .

# RUN chmod +x ./config

# ENV PORT 8081

# CMD ./config -configPrefix=wedding -configSource=/opt/event/authenticator/.env

 