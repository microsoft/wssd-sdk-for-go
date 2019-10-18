# Copyright (c) Microsoft Corporation.
# Licensed under the MIT license.

# Create the CA
cfssl gencert -initca ca-csr.json | cfssljson -bare ca

# Generate the wssd service certs
cfssl gencert \
  -ca=ca.pem \
  -ca-key=ca-key.pem \
  -config=ca-config.json \
  -profile=server \
  wssdagent-csr.json | cfssljson -bare wssd
