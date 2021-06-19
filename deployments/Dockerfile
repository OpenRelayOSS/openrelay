FROM centos:8
MAINTAINER "kyadet"

RUN dnf -y install epel-release
RUN dnf -y install libunwind 
EXPOSE 7000
EXPOSE 7001
EXPOSE 7002
EXPOSE 7003
EXPOSE 7004
EXPOSE 7005
EXPOSE 7006
EXPOSE 8000

ADD build/rpms/RPMS/x86_64/openrelay-0.9-9.x86_64.rpm /
RUN rpm -ivh openrelay-0.9-9.x86_64.rpm
ENTRYPOINT ["/usr/local/openrelay/bin/openrelay-boot.sh"]
