FROM azuresdk/azure-cli-python
ENV PORT 4392
EXPOSE 4392
ADD rootfs /go/src/github.com/kris-nova/terraformctl/rootfs
RUN cd /go/src/github.com/kris-nova/terraformctl/rootfs && \
    mv /go/src/github.com/kris-nova/terraformctl/rootfs/terraform-bin /usr/local/bin/terraform && \
    mv /go/src/github.com/kris-nova/terraformctl/rootfs/.azure ~/.azure

ENTRYPOINT /go/src/github.com/kris-nova/terraformctl/rootfs/terraformctl serve

