# Manufacture

修改后的参数

```yaml
port：7000
tls:
  # Enable TLS (default: false)
  enabled: true
ca:
  # Name of this CA
  name: manufacturer-ca
  
csr:
   cn: manufacturer-ca
   keyrequest:
     algo: ecdsa
     size: 256
   names:
      - C: CN
        ST: "NanJing"
        L:
        O: cBPM
        OU: ADMIN
   hosts:
     - localhost
     - 127.0.0.1
   ca:
      expiry: 131400h
      pathlength: 1
operations:
    # host and port for the operations server
    listenAddress: 127.0.0.1:9000
```

