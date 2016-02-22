FROM debian:jessie

ENV CHRONOGRAF_VERSION=0.10.0-rc1

RUN apt-get update && apt-get install -y curl
RUN curl -s -o /tmp/cronograpf_latest_amd64.deb https://s3.amazonaws.com/get.influxdb.org/chronograf/chronograf_${CHRONOGRAF_VERSION}_amd64.deb && \
  dpkg -i /tmp/cronograpf_latest_amd64.deb && \
  rm /tmp/cronograpf_latest_amd64.deb && \
  rm -rf /var/lib/apt/lists/*

# Allow Chronograf to accept connections from other hosts
ENV CHRONOGRAF_BIND=0.0.0.0:10000
EXPOSE 10000

CMD [ "/opt/chronograf/chronograf" ]
