FROM scratch

ENV FOOMO_BERT_MAIN_MODULE Foomo

COPY bin/foomo-bert-linux-amd64 /usr/sbin/foomo-bert

# install ca root certificates
# https://curl.haxx.se/docs/caextract.html
# http://blog.codeship.com/building-minimal-docker-containers-for-go-applications/
ADD https://curl.haxx.se/ca/cacert.pem /etc/ssl/certs/ca-certificates.crt

ENTRYPOINT ["/usr/sbin/foomo-bert"]
