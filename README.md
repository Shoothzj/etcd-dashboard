# etcd-dashboard
etcd dashboard project for fun.

## backend api command
### put key
```bash
curl -X PUT -d '{"key":"test","value":"test"}' http://localhost:10001/api/etcd/keys
```
### list keys
```bash
curl localhost:10001/api/etcd/keys
```
### get key
```bash
curl localhost:10001/api/etcd/keys/k1
```
