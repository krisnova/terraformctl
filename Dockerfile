FROM azuresdk/azure-cli-python
ENV PORT 4392
ENV GOPATH /go
ENV GOROOT /usr/local/go
ENV PATH /usr/local/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/local/go/bin
EXPOSE 4392
ADD . /go/src/github.com/kris-nova/terraformctl
RUN cd /go/src/github.com/kris-nova/terraformctl && \
    /go/src/github.com/kris-nova/terraformctl/go/bin/go build -o bin/terraformctl -ldflags "-X github.com/kris-nova/terraformctl/cmd.GitSha=${GIT_SHA} -X github.com/kris-nova/terraformctl/cmd.Version=${VERSION}" main.go && \
    mv /go/src/github.com/kris-nova/terraformctl/terraform-bin /usr/local/bin/terraform && \
    curl -L https://aka.ms/InstallAzureCli | bash && \
    mv /go/src/github.com/kris-nova/terraformctl/.azure ~/.azure

ENTRYPOINT /go/src/github.com/kris-nova/terraformctl/bin/terraformctl serve

