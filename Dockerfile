FROM azuresdk/azure-cli-python
ENV PORT 4392
EXPOSE 4392
ADD . /go/src/github.com/kris-nova/terraformctl
RUN cd /go/src/github.com/kris-nova/terraformctl && \
    make && \
    mv /go/src/github.com/kris-nova/terraformctl/terraform-bin /usr/local/bin/terraform && \
    curl -L https://aka.ms/InstallAzureCli | bash && \
    mv /go/src/github.com/kris-nova/terraformctl/.azure ~/.azure

ENTRYPOINT /go/src/github.com/kris-nova/terraformctl/bin/terraformctl serve

