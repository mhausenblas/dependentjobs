FROM golang:1.8
ADD . /go/src/dependentjobs
RUN go get gopkg.in/yaml.v2
WORKDIR /go/src/dependentjobs
RUN go install
CMD dependentjobs

