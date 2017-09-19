FROM golang:onbuild
ENV PORT 4392
EXPOSE 4392
ADD . /go/src/github.com/kris-nova/terraformctl
RUN cd /go/src/github.com/kris-nova/terraformctl && make
ENTRYPOINT /go/src/github.com/kris-nova/terraformctl/bin/terraformctl serve

