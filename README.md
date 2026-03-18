# Cloudflare DDNS Updater

Use Cloudflare a DDNS provider with this tool on crontab.

## IPv4 and IPv6 Support

This tool supports both IPv4 (A records) and IPv6 (AAAA records) in a single instance/container:

- **IPv4**: Automatically detected via `api.ipify.org` and updates A records
- **IPv6**: Automatically detected via `api64.ipify.org` and updates AAAA records
- If IPv6 is not available, the tool will log a warning and continue with IPv4 only

You only need **one container** to handle both IPv4 and IPv6.

```
$> ./cf-ddns --help
usage: cf-ddns --cf-email=CF-EMAIL --cf-api-key=CF-API-KEY --cf-zone-id=CF-ZONE-ID [<flags>] <hostnames>...

Cloudflare DynDNS Updater

Flags:
  --help                   Show context-sensitive help (also try --help-long and --help-man).
  --ip-address=IP-ADDRESS  Skip resolving external IP and use provided IP (IPv4)
  --ipv6-address=IPV6-ADDRESS  Skip resolving external IPv6 and use provided IPv6
  --no-verify              Don't verify ssl certificates
  --cf-email=CF-EMAIL      Cloudflare Email
  --cf-api-key=CF-API-KEY  Cloudflare API key
  --cf-zone-id=CF-ZONE-ID  Cloudflare Zone ID

Args:
  <hostnames>  Hostnames to update
```

## Docker Usage

### Using pre-built image

```bash
docker run --rm \
  -e CF_EMAIL=your@email.com \
  -e CF_API_KEY=your_api_key \
  -e CF_ZONE_ID=your_zone_id \
  ghcr.io/favonia/cloudflare-ddns:latest \
  your.domain.com
```

### Building locally

```bash
docker build -t cf-ddns .
docker run --rm \
  -e CF_EMAIL=your@email.com \
  -e CF_API_KEY=your_api_key \
  -e CF_ZONE_ID=your_zone_id \
  cf-ddns \
  your.domain.com
```

### Environment Variables

| Variable | Description | Required |
|----------|-------------|----------|
| `CF_EMAIL` | Cloudflare account email | Yes |
| `CF_API_KEY` | Cloudflare API key (Global or Zone Edit) | Yes |
| `CF_ZONE_ID` | Cloudflare Zone ID (found in DNS settings) | Yes |

### Docker Compose Example

```yaml
version: '3.8'
services:
  cf-ddns:
    image: ghcr.io/favonia/cloudflare-ddns:latest
    environment:
      CF_EMAIL: your@email.com
      CF_API_KEY: your_api_key
      CF_ZONE_ID: your_zone_id
    command: your.domain.com
    # Optional: run periodically
    restart: unless-stopped
```

## Compiling for MIPS (Ubnt Edgerouter Lite)
```
GOOS=linux GOARCH=mips64 go build -o cf-ddns-mips64
```
