FROM ubuntu:22.04
RUN apt-get update && apt-get install curl software-properties-common unzip -y
RUN curl -fsSL https://apt.releases.hashicorp.com/gpg | apt-key  add -
RUN apt-add-repository "deb [arch=amd64] https://apt.releases.hashicorp.com $(lsb_release -cs) main"
RUN apt-get update && apt-get install vault -y

ARG pluginversion=0.3.0
RUN mkdir -p /opt/vault-plugins
RUN curl -L -o vault-plugin-secrets-keycloak.zip \
    https://github.com/Serviceware/vault-plugin-secrets-keycloak/releases/download/v${pluginversion}/vault-plugin-secrets-keycloak_${pluginversion}_linux_amd64.zip && \
    unzip vault-plugin-secrets-keycloak.zip -d /opt/ && \
    mv /opt/vault-plugin-secrets-keycloak_v${pluginversion} /opt/vault-plugins/vault-plugin-secrets-keycloak

RUN chmod +x /opt/vault-plugins/vault-plugin-secrets-keycloak
EXPOSE 8200

CMD ["vault", "server", "-dev", "-dev-root-token-id=root", "-dev-plugin-dir=/opt/vault-plugins", "-dev-listen-address=0.0.0.0:8200"]