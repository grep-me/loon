ARG BASE_IMAGE=artifactory.scitec.com/ib/ironbank/redhat/ubi/ubi9:latest

FROM $BASE_IMAGE AS ubi9-installer

WORKDIR /
USER root

# Install RPMs:
RUN dnf install -y --nogpgcheck wget \
    unzip \
    tar \
    gcc \
    clang \
    git && \
    dnf -q clean all && \
    rm -rf /var/cache/dnf

# Grab & Make ZMQ:
RUN cd /tmp && wget https://github.com/zeromq/libzmq/releases/download/v4.3.5/zeromq-4.3.5.zip && \
    unzip zeromq-4.3.5.zip && rm -rf zeromq-4.3.5.zip && \
    cd zeromq-4.3.5 && \
    ./configure && \
    make && \
    make install

# Install GO: (Incase of Dev)
RUN cd /tmp && wget https://go.dev/dl/go1.23.0.src.tar.gz && \
    rm -rf /usr/local/go && \ 
    tar -xf go1.23.0.src.tar.gz && \
    mv ./go /usr/local

# Grab Loon Tool:
RUN --mount=type=secret,id=gitlab_username \
    --mount=type=secret,id=gitlab_token \
    cd /tmp && \
    git clone https://$(cat /run/secrets/gitlab_username):$(cat /run/secrets/gitlab_token)@gitlab.scitec.com/mdpap/cyber/fuzzing.git


FROM ubi9-installer AS loon-tool

COPY --from=ubi9-installer /tmp/fuzzing/loon /app
COPY --from=ubi9-installer /usr/local /usr/local
COPY --from=ubi9-installer /tmp /tmp
COPY --from=ubi9-installer /usr/local/lib /usr/local/lib

WORKDIR /app
USER root
ENV LD_LIBRARY_PATH=/usr/local/lib:/usr/lib:/lib
ENV PATH=$PATH:/usr/local/go/bin:/app

ENTRYPOINT ["loon"]





