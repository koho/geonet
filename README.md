# Geolocation Networking Data

## Data source

| Item    | Source                                                                               |
|---------|--------------------------------------------------------------------------------------|
| geoip   | https://github.com/v2fly/geoip/releases/latest/download/geoip.dat                    |
| geosite | https://github.com/Loyalsoldier/v2ray-rules-dat/releases/latest/download/geosite.dat |

## Common query parameters

| Query   | Type   | Description    |
|---------|--------|----------------|
| country | string | IP Geolocation |
| format  | string | Output format  |

## Formatters

### `text` formatter

No other options

### `shell` formatter

| Query  | Type   | Description               |
|--------|--------|---------------------------|
| opt    | string | Options for shell command |

### `bind` formatter

| Query   | Type   | Description                   |
|---------|--------|-------------------------------|
| include | string | File path of the zone options |

### `dnsmasq` formatter

| Query | Type   | Description        |
|-------|--------|--------------------|
| dns   | string | DNS server address |

## API

### geoip

| Query | Type | Description                |
|-------|------|----------------------------|
| type  | int  | `4` for IPv4; `6` for IPv6 |

Supported formatters:

- `text`
- `shell`

```shell
curl http://127.0.0.1:8080/geoip?country=CN&type=4&format=shell&opt=dev%20pppoe-wan
```

### geosite

Supported formatters:

- `bind`
- `dnsmasq`

```shell
curl http://127.0.0.1:8080/geosite?country=CN&format=bind&include=/etc/bind/named.china.zones
```
