FROM golang:1.18.1 AS build
# RUN go .env 
# RUN mkdir /config
WORKDIR /authen-docker 
COPY go.mod go.sum ./
RUN go mod download 
# RUN mkdir /app

COPY . . 
RUN go build -o authen /authen-docker/controller/*.go
WORKDIR /authen-docker
CMD ["./controller"]


FROM ubuntu:22.04 AS runtime
RUN mkdir /app-authen
WORKDIR /app-authen
 
 COPY --from=build /authen-docker/authen .
#  COPY build/env  /authen-docker/controller/app-authen/.env
 COPY authen/cmd/* /authen-docker/
 
 RUN bash -c "chmod +x ./config"

ENV PORT 8081

CMD ./config -configPrefix=wedding -configSource=/opt/event/authenticator/.env

 



# FROM ubuntu:22.04 AS runtime
# # RUN apt updated
# WORKDIR /opt/event/authenticator

# COPY ./controller/config  .

# COPY ./controller/.env .

# RUN chmod +x ./config

# ENV PORT 8081

# CMD ./config -configPrefix=wedding -configSource=/opt/event/authenticator/.env

 