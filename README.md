# udppunch

udp punch for wireguard, inspired by [natpunch-go](https://github.com/malcolmseyd/natpunch-go)

## usage

server side

```bash
./punch-server-linux-amd64 -port 56000
```

client side

> make sure wireguard is up

```bash
./dist/punch-client-linux-amd64 -server xxxx:56000 -iface wg0
```

## resource

- [natpunch-go](https://github.com/malcolmseyd/natpunch-go) (because of [#7](https://github.com/malcolmseyd/natpunch-go/issues/7) not support macOS, so I build this)
- [wireguard-vanity-address](https://github.com/yinheli/wireguard-vanity-address) generate keypairs with a given prefix
- [UDP hole punching](https://en.wikipedia.org/wiki/UDP_hole_punching)


## notes

```bash
./punch-server-linux-arm64 -port 56000

wg-quick up ./conf/w1.conf
wg-quick down ./conf/w1.conf
./punch-client-linux-arm64 -server w0:56000 -iface w1

wg-quick up ./conf/w2.conf
wg-quick down ./conf/w2.conf
./punch-client-linux-arm64 -server w0:56000 -iface w2
```

先执行 client 的是 从端



