# API

Definitions:
- Cloud Manager - main service, responsible for storing and distributuing WireGuard configuration and changes
- Cloud Agent -  service running on target host, responsible for applying WireGuard configuration on host

## `/topology`

This endpoint return WireGuard configuration (aka topology) for Cloud Agent

### Request

```http
POST /topology HTTP/1.1
Content-Type: application/json

```
```jsonc
{
  "PublicKey:" "yAnz5TF+lXXJte14tji3zlMNq+hd2rYUIgJBgB3fBmk=", //WireGuard public key
  "Signature": "ASDVASTDCYUADLIA...=" //* ed25519 signature of some random (optional, opinionated, implement later), simpliest way to authenticate Cloud Agent on Cloud Manager
}
```

### Response

```http
HTTP/1.1 200
Content-Type: application/json
```
```jsonc
{
    "config": {
        // https://github.com/WireGuard/wgctrl-go/blob/master/wgtypes/types.go#L206
        // Field comment are available at the link
        "ListenPort": 51820, // WireGuard Listen port
        "Peers": [{
            // https://github.com/WireGuard/wgctrl-go/blob/master/wgtypes/types.go#L235
            "PublicKey": "yAnz5TF+lXXJte14tji3zlMNq+hd2rYUIgJBgB3fBmk=", // WireGuard peer public key
            "Endpoint": "wg.example.com:51820", // WireGuard peer endpoint
            "PersistentKeepAliveInterval": "25s", // WireGuard peer keep alive interval
            "AllowedIPs": "10.11.0.76/32" // WireGuard peer allowed IPs
        }]
    }
}
```