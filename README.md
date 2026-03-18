# Cloudflare DDNS Updater

Use Cloudflare a DDNS provider with this tool on crontab.

## IPv4 and IPv6 Support

This tool supports both IPv4 (A records) and IPv6 (AAAA records) in a single instance/container:

- **IPv4**: Automatically detected via `ipv4.ip.sb` and updates A records
- **IPv6**: Automatically detected via `ipv6.ip.sb` and updates AAAA records
- IPv4 and IPv6 updates are completely independent - you can update only IPv4, only IPv6, or both

You only need **one container** to handle both IPv4 and IPv6.

```
$> ./cf-ddns --help
usage: cf-ddns --cf-email=CF-EMAIL --cf-api-key=CF-API-KEY --cf-zone-id=CF-ZONE-ID [<flags>]

Cloudflare DynDNS Updater

Flags:
  --help                   Show context-sensitive help (also try --help-long and --help-man).
  --ip-address=IP-ADDRESS  Skip resolving external IP and use provided IP (IPv4)
  --ipv6-address=IPV6-ADDRESS  Skip resolving external IPv6 and use provided IPv6
  --no-verify              Don't verify ssl certificates
  --interval=0             Run in loop mode, checking IP every N minutes (0 = run once)
  --cf-email=CF-EMAIL      Cloudflare Email
  --cf-api-key=CF-API-KEY  Cloudflare API key
  --cf-zone-id=CF-ZONE-ID  Cloudflare Zone ID
```

## Loop Mode

The `--interval` parameter allows the program to run continuously and check for IP changes at a specified interval (in minutes). This is useful when you don't want to use cron or external schedulers.

### How Loop Mode Works

- When `--interval` is set to a value greater than 0, the program enters loop mode
- The program will check your IP every N minutes
- **DNS records are only updated when the IP changes**, avoiding unnecessary API calls to Cloudflare
- The program runs once immediately on startup, then waits for the interval before the next check

### Loop Mode Example

Run in loop mode, checking IP every 10 minutes:

```bash
docker run --rm \
  -e CF_EMAIL=your@email.com \
  -e CF_API_KEY=your_api_key \
  -e CF_ZONE_ID=your_zone_id \
  -e IPV4_HOSTNAME=your.domain.com \
  -e IPV6_HOSTNAME=your.domain.com \
  ety001/cf-ddns:latest \
  --interval 10
```

### Loop Mode with Docker Compose

```yaml
version: '3.8'
services:
  cf-ddns:
    image: ety001/cf-ddns:latest
    environment:
      CF_EMAIL: your@email.com
      CF_API_KEY: your_api_key
      CF_ZONE_ID: your_zone_id
      IPV4_HOSTNAME: your.domain.com
      IPV6_HOSTNAME: your.domain.com
    command: ["--interval", "10"]  # Check every 10 minutes
    restart: unless-stopped
```

### Comparing Loop Mode vs Cron

| Method | Pros | Cons |
|--------|------|------|
| **Loop Mode (`--interval`)** | Simple setup, no external scheduler needed | Container stays running |
| **Cron** | Container runs only when needed, more traditional approach | Requires cron or external scheduler |

Choose loop mode for simplicity, or cron for a more traditional "run and exit" approach.

## Docker Usage

### Using pre-built image

Update both IPv4 and IPv6 (one-time run):

```bash
docker run --rm \
  -e CF_EMAIL=your@email.com \
  -e CF_API_KEY=your_api_key \
  -e CF_ZONE_ID=your_zone_id \
  -e IPV4_HOSTNAME=your.domain.com \
  -e IPV6_HOSTNAME=your.domain.com \
  ety001/cf-ddns:latest
```

Update only IPv4:

```bash
docker run --rm \
  -e CF_EMAIL=your@email.com \
  -e CF_API_KEY=your_api_key \
  -e CF_ZONE_ID=your_zone_id \
  -e IPV4_HOSTNAME=your.domain.com \
  ety001/cf-ddns:latest
```

Update only IPv6:

```bash
docker run --rm \
  -e CF_EMAIL=your@email.com \
  -e CF_API_KEY=your_api_key \
  -e CF_ZONE_ID=your_zone_id \
  -e IPV6_HOSTNAME=your.domain.com \
  ety001/cf-ddns:latest
```

### Building locally

```bash
docker build -t cf-ddns .
docker run --rm \
  -e CF_EMAIL=your@email.com \
  -e CF_API_KEY=your_api_key \
  -e CF_ZONE_ID=your_zone_id \
  -e IPV4_HOSTNAME=your.domain.com \
  -e IPV6_HOSTNAME=ipv6.domain.com \
  cf-ddns
```

### Cron Example

Run every 5 minutes:

```bash
*/5 * * * * docker run --rm ety001/cf-ddns:latest \
  --cf-email your@email.com \
  --cf-api-key your_api_key \
  --cf-zone-id your_zone_id \
  -e IPV4_HOSTNAME=your.domain.com \
  -e IPV6_HOSTNAME=your.domain.com
```

### Environment Variables

| Variable | Description | Required |
|----------|-------------|----------|
| `CF_EMAIL` | Cloudflare account email | Yes |
| `CF_API_KEY` | Cloudflare API key (Global or Zone Edit) | Yes |
| `CF_ZONE_ID` | Cloudflare Zone ID (found in DNS settings) | Yes |
| `IPV4_HOSTNAME` | Hostname for IPv4 (A record) | At least one of IPV4_HOSTNAME or IPV6_HOSTNAME |
| `IPV6_HOSTNAME` | Hostname for IPv6 (AAAA record) | At least one of IPV4_HOSTNAME or IPV6_HOSTNAME |
| `IPV4_ENDPOINT` | IPv4 detection endpoint (default: `https://ipv4.ip.sb`) | No |
| `IPV6_ENDPOINT` | IPv6 detection endpoint (default: `https://ipv6.ip.sb`) | No |

### Custom IP Detection Endpoints

You can override the default IP detection services by setting the `IPV4_ENDPOINT` and `IPV6_ENDPOINT` environment variables:

```bash
docker run --rm \
  -e CF_EMAIL=your@email.com \
  -e CF_API_KEY=your_api_key \
  -e CF_ZONE_ID=your_zone_id \
  -e IPV4_HOSTNAME=your.domain.com \
  -e IPV6_HOSTNAME=your.domain.com \
  -e IPV4_ENDPOINT=https://api.ipify.org?format=json \
  -e IPV6_ENDPOINT=https://api64.ipify.org?format=json \
  ety001/cf-ddns:latest
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
    image: ety001/cf-ddns:latest
    environment:
      CF_EMAIL: your@email.com
      CF_API_KEY: your_api_key
      CF_ZONE_ID: your_zone_id
      IPV4_HOSTNAME: your.domain.com
      IPV6_HOSTNAME: ipv6.domain.com
    # Optional: run periodically
    restart: unless-stopped
```

## Compiling for MIPS (Ubnt Edgerouter Lite)
```
GOOS=linux GOARCH=mips64 go build -o cf-ddns-mips64
```
