# Cloudflare DDNS Updater

Use Cloudflare a DDNS provider with this tool on crontab.

## IPv4 and IPv6 Support

This tool supports both IPv4 (A records) and IPv6 (AAAA records) in a single instance/container:

- **IPv4**: Automatically detected via `ipv4.ip.sb` and updates A records
- **IPv6**: Automatically detected via `ipv6.ip.sb` and updates AAAA records
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
  ety001/cf-ddns:latest \
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

### Docker Compose Example (with cron schedule)

```yaml
version: '3.8'
services:
  cf-ddns:
    image: ety001/cf-ddns:latest
    environment:
      CF_EMAIL: your@email.com
      CF_API_KEY: your_api_key
      CF_ZONE_ID: your_zone_id
    command: your.domain.com
    restart: unless-stopped
    # Optional: run every 5 minutes
    # Use with docker-compose + cron or a scheduler
```

### Cron Example

Run every 5 minutes:

```bash
*/5 * * * * docker run --rm ety001/cf-ddns:latest \
  --cf-email your@email.com \
  --cf-api-key your_api_key \
  --cf-zone-id your_zone_id \
  your.domain.com
```

### Environment Variables

| Variable | Description | Required |
|----------|-------------|----------|
| `CF_EMAIL` | Cloudflare account email | Yes |
| `CF_API_KEY` | Cloudflare API key (Global or Zone Edit) | Yes |
| `CF_ZONE_ID` | Cloudflare Zone ID (found in DNS settings) | Yes |
| `IPV4_ENDPOINT` | IPv4 detection endpoint (default: `https://ipv4.ip.sb`) | No |
| `IPV6_ENDPOINT` | IPv6 detection endpoint (default: `https://ipv6.ip.sb`) | No |

### Custom IP Detection Endpoints

You can override the default IP detection services by setting the `IPV4_ENDPOINT` and `IPV6_ENDPOINT` environment variables:

```bash
docker run --rm \
  -e CF_EMAIL=your@email.com \
  -e CF_API_KEY=your_api_key \
  -e CF_ZONE_ID=your_zone_id \
  -e IPV4_ENDPOINT=https://api.ipify.org?format=json \
  -e IPV6_ENDPOINT=https://api64.ipify.org?format=json \
  ety001/cf-ddns:latest \
  your.domain.com
```

The custom endpoint must return a JSON response with an `ip` field:
```json
{"ip": "1.2.3.4"}
```

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
