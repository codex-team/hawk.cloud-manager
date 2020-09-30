# cloud-manager

## Wireguard

### server

- Stores peer configuration
- Bootstraps new peers
- Sends new config

### client

- Syncs peer config
- Initiates new peer request

## Usage

Run manager:
```shell
$ make
$ ./manager -addr <address to listen, format: 0.0.0.0:50051> -config /path/to/config.yaml
```

Run agent:
```shell
$ make agent
$ ./agent -config /path/to/store/config -manager <cloud-manager address, format: 0.0.0.0:50051> -pubkey /path/to/public/key -privkey /path/to/private/key -interval <time interval to check config changes, format: 5s>
```

Run integration tests:
```shell
$ make int
```

Run unit tests:
```shell
$ make ut
```

## Roadmap

MVP:
- [x] Read config from yaml
- [ ] Serve as HTTP (maybe GRPC?) server
- [ ] Support methods:
    - [ ] Get running config
    - [ ] Bootstrap request
- [ ] Agent pull config every n seconds

Mid-Term:
- Read config from Consul
- Auth via mTLS

Long-term:
- Support ACL and/or RBAC
- Web ui/API for requesting access for admins/devs
- ENV updater (separate agent)
- Cert manager (separate server + agent)
- ...GitOps manager for all workflows
