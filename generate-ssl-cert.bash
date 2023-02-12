# Create directory for certs if it does not exist
mkdir -p load-balancer/certs


# Generate SSL certs and private key
openssl req -x509 -out load-balancer/certs/localhost.crt -keyout load-balancer/certs/localhost.key \
    -days 7 \
    -newkey rsa:2048 -nodes -sha256 \
    -subj '/CN=localhost' -extensions EXT -config <( \
    printf "[dn]\nCN=localhost\n[req]\ndistinguished_name = dn\n[EXT]\nsubjectAltName=DNS:localhost\nkeyUsage=digitalSignature\nextendedKeyUsage=serverAuth")


# Add ssl cert to trusted cert list
OS="`uname`"
case "${OS}" in
    'Linux'*)
        sudo cp load-balancer/certs/localhost.crt /usr/local/share/ca-certificates/localhost.crt
        sudo update-ca-certificates
        ;;
    'Darwin'*)
        sudo security add-trusted-cert -d -r trustRoot -k /Library/Keychains/System.keychain load-balancer/certs/localhost.crt
        ;;
esac