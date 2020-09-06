# API

Definitions:

- Cloud Manager - main service, responsible for storing and distributuing WireGuard configuration and changes
- Cloud Agent - service running on target host, responsible for applying WireGuard configuration on host

## `POST /topology`

This endpoint return WireGuard configuration (aka topology) for Cloud Agent

### Request

```http
POST /topology HTTP/1.1
Content-Type: application/json

```

```jsonc
{
  "public_key:" "yAnz5TF+lXXJte14tji3zlMNq+hd2rYUIgJBgB3fBmk=", //WireGuard public key
  "signature": "ASDVASTDCYUADLIA...=" //* ed25519 signature of some random (optional, opinionated, implement later), simpliest way to authenticate Cloud Agent on Cloud Manager
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
    "listen_port": 51820, // WireGuard Listen port
    "peers": [
      {
        // https://github.com/WireGuard/wgctrl-go/blob/master/wgtypes/types.go#L235
        "public_key": "yAnz5TF+lXXJte14tji3zlMNq+hd2rYUIgJBgB3fBmk=", // WireGuard peer public key
        "endpoint": "wg.example.com:51820", // WireGuard peer endpoint
        "keep_alive_interval": "25s", // WireGuard peer keep alive interval
        "allowed_ips": ["10.11.0.76/32"] // WireGuard peer allowed IPs
      }
    ]
  }
}
```

## `GET /config`

Returns Manager's current configuration, gets it from initialized storage. See `PeerConfig`

## Request

```http
GET /config HTTP/1.1

```

```jsonc

```

## Response

```http
HTTP/1.1 200
Content-Type: application/json
```

```jsonc
// Serialized `PeerConfig`
{
  "hosts": [
    {
      "name": "hawk-main",
      "public_key": "yAnz5TF+lXXJte14tji3zlMNq+hd2rYUIgJBgB3fBmk=",
      "endpoint": "172.17.123.12",
      "allowed_ips": ["10.11.0.76/32"]
    }
  ],
  "groups": [
    {
      "name": "hawk-developers",
      "hosts": ["hawk-main"]
    }
  ]
}
```

## `PUT /config`

Updates Manager's current configuration. See `PeerConfig`

## Request

```http
PUT /config HTTP/1.1
Content-Type: application/json

```

```jsonc
{
  "hosts": [
    {
      "name": "hawk-main",
      "public_key": "yAnz5TF+lXXJte14tji3zlMNq+hd2rYUIgJBgB3fBmk=",
      "endpoint": "172.17.123.12",
      "allowed_ips": ["10.11.0.76/32"]
    },
    {
      "name": "hawk-admin",
      "public_key": "GAnz5TF+lXXJte1asdASDi3zlMNq+hd2rYUIgJBgB3fBmk=",
      "endpoint": "172.17.123.13",
      "allowed_ips": ["10.11.0.77/32"]
    }
  ],
  "groups": [
    {
      "name": "hawk-developers",
      "hosts": ["hawk-main", "hawk-admin"]
    }
  ]
}
```

## Response

```http
HTTP/1.1 200
```

```jsonc

```
