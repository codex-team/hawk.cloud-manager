# API

## `/topology`

Get WG peers for current host

### Request

```http
POST /topology HTTP/1.1
Content-Type: application/json

{
  "PublicKey:" "yAnz5TF+lXXJte14tji3zlMNq+hd2rYUIgJBgB3fBmk=", // Host's public key
  "Signature": "ASDVASTDCYUADLIA...=" //* ed25519 signature
}
```

### Response

```http
HTTP/1.1 200
Content-Type: application/json

{
    "config": {
        // https://github.com/WireGuard/wgctrl-go/blob/master/wgtypes/types.go#L206
        "ListenPort": 51820,
        "Peers": [
            // https://github.com/WireGuard/wgctrl-go/blob/master/wgtypes/types.go#L235
            "PublicKey": "yAnz5TF+lXXJte14tji3zlMNq+hd2rYUIgJBgB3fBmk=",
            "Endpoint": "wg.example.com:51820",
            "PersistentKeepAliveInterval": "25s",
            "AllowedIPs": "10.11.0.76/32"
        ]
    }
}
```