FROM azuresdk/azure-cli-python
ENV PORT 4392
ENV GOPATH /go
ENV GOROOT /usr/local/go
ENV PATH /usr/local/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/local/go/bin
EXPOSE 4392
COPY . /go/src/github.com/kris-nova/terraformctl
RUN cd /go/src/github.com/kris-nova/terraformctl && \
    mv /go/src/github.com/kris-nova/terraformctl/terraform-bin /usr/local/bin/terraform && \
    mv /go/src/github.com/kris-nova/terraformctl/.azure ~/.azure

ENTRYPOINT /go/src/github.com/kris-nova/terraformctl/terraformctl serve

