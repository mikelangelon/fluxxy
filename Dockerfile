FROM golang
 
ADD . /go/src
RUN go install game/cmd/fluxxy
ENTRYPOINT /go/bin/fluxxy
 
EXPOSE 8080
