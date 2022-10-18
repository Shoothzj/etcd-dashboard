FROM shoothzj/base

COPY dist /opt/etcd-dashboard

WORKDIR /opt/etcd-dashboard

EXPOSE 10001

CMD ["/usr/bin/dumb-init", "/opt/etcd-dashboard/etcd-dashboard"]
