FROM golang as builder 
RUN mkdir /build 
ADD . /build 
WORKDIR /build
RUN go build


#ACTUAL
FROM alpine
RUN apk update && apk add ca-certificates

#copy 
COPY --from=builder /build/build /usr/local/bin/goroute
RUN chmod +x /usr/local/bin/goroute

CMD /usr/local/bin/goroute