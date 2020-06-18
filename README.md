# cloud-manager

## Wireguard

### server

- Stores peer configuration
- Bootstraps new peers
- Sends new config

### client

- Syncs peer config
- Initiates new peer request

## Roadmap

MVP:
- [ ] Read config from yaml
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