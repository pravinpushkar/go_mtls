## Creating ca certificate, server certificate and client certificate:
openssl ecparam -name prime256v1 -genkey -noout -out ca.key 
openssl req -new -sha256 -key ca.key -out ca.csr -config root.conf -extensions v3_req
openssl x509 -req -sha256 -days 365 -in ca.csr -signkey ca.key -outform PEM -out ca.pem -extfile root.conf -extensions v3_req

openssl ecparam -name prime256v1 -genkey -noout -out server.key
openssl req -new -sha256 -key server.key -out server.csr -config server.conf -extensions v3_req
openssl x509 -req -in server.csr -CA ca.pem -CAkey ca.key -CAcreateserial -outform PEM -out server.pem -days 365 -sha256 -extfile server.conf -extensions v3_req

openssl ecparam -name prime256v1 -genkey -noout -out client.key
openssl req -new -sha256 -key client.key -out client.csr -config client.conf -extensions v3_req
openssl x509 -req -in client.csr -CA ca.pem -CAkey ca.key -CAcreateserial -outform PEM -out client.pem -days 365 -sha256 -extfile client.conf -extensions v3_req


## Use curl verify the server https endpoint:
```
curl --trace trace.log -k --cacert <ca cert file path>  --cert <client cert file path> --key <client cert key file path>  https://localhost:8443/hello
```