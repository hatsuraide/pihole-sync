# pihole-sync

[![Release version](https://img.shields.io/github/v/release/Kizuren/pihole-sync)](https://github.com/Kizuren/pihole-sync/releases/latest)
[![Tests](https://img.shields.io/github/actions/workflow/status/Kizuren/pihole-sync/go.yml?branch=main&label=tests)](https://github.com/Kizuren/pihole-sync/actions/workflows/go.yml?query=branch%3Amain)
![Go version](https://img.shields.io/github/go-mod/go-version/Kizuren/pihole-sync)
[![Docker image size](https://img.shields.io/docker/image-size/Kizuren/pihole-sync/latest)](https://hub.docker.com/r/Kizuren/pihole-sync)

Synchronize Pi-hole v6.x configuration to replicas.

This project is not a part of the [official Pi-hole project](https://github.com/pi-hole), but uses the api provided by Pi-hole instances to perform the synchronization actions.

> pihole-sync is a fork of [nebula-sync](https://github.com/lovelaze/nebula-sync) by Fredrik Berntsson, with added Pi-hole 2FA/TOTP support and a native OS service so it can run continuously without Docker.

## Features
- **Full sync**: Use Pi-hole Teleporter for full synchronization.
- **Selective sync**: Selective feature synchronization.
- **Cron schedule**: Run on cron schedule.
- **2FA/TOTP**: Authenticate against Pi-hole instances that have two-factor authentication enabled.
- **Native service**: Install as a systemd/launchd/Windows service that starts on boot and restarts on failure — no Docker required.

## Installation

### Native binary + service (recommended, no Docker required)

Download the binary from the [latest release](https://github.com/Kizuren/pihole-sync/releases/latest) or build from source:
```
go install github.com/Kizuren/pihole-sync@latest
```

Create an env file (see [Configuration](#configuration) below and [`examples/.env_full`](examples/.env_full)), then install pihole-sync as a background service. It will start on boot and restart itself automatically if it ever crashes — the same job Docker's restart policy + cron would otherwise do, without needing Docker at all:

```bash
# install as a systemd (Linux) / launchd (macOS) / Windows service
sudo pihole-sync service install --env-file /etc/pihole-sync/.env

# start it now
sudo pihole-sync service start

# check status / stop / uninstall
sudo pihole-sync service status
sudo pihole-sync service stop
sudo pihole-sync service uninstall
```

`CRON` must be set in the env file for the service to keep syncing on a schedule; without it, pihole-sync performs a single sync and exits.

You can also just run it in the foreground yourself (e.g. under your own supervisor, or for a one-off manual sync):
```bash
# run
pihole-sync run

# read envs from file
pihole-sync run --env-file .env
```

### Docker Compose (optional)

```yaml
---
services:
  pihole-sync:
    image: ghcr.io/Kizuren/pihole-sync:latest
    container_name: pihole-sync
    environment:
    - PRIMARY=http://ph1.example.com|password
    - REPLICAS=http://ph2.example.com|password,http://ph3.example.com|password
    - FULL_SYNC=true
    - RUN_GRAVITY=true
    - CRON=0 * * * *
```

### Docker CLI

```bash
docker run --rm \
  --name pihole-sync \
  -e PRIMARY="http://ph1.example.com|password" \
  -e REPLICAS="http://ph2.example.com|password" \
  -e FULL_SYNC=true \
  -e RUN_GRAVITY=true \
  ghcr.io/Kizuren/pihole-sync:latest
```

## Examples
Env and docker-compose examples can be found [here](https://github.com/Kizuren/pihole-sync/tree/main/examples)

## Running as a service

`pihole-sync service` manages pihole-sync as a native background service via [systemd](https://systemd.io/) (Linux), [launchd](https://en.wikipedia.org/wiki/Launchd) (macOS) or the Windows Service Manager — whichever the OS provides. This is the recommended way to run pihole-sync continuously without Docker: the service starts on boot and is restarted automatically if it ever crashes, the same guarantees a Docker restart policy would give you.

| Command                                          | Description                                                          |
|---------------------------------------------------|------------------------------------------------------------------|
| `pihole-sync service install --env-file <path>`  | Registers pihole-sync as a service that runs `pihole-sync run` using the given env file |
| `pihole-sync service start`                      | Starts the installed service                                       |
| `pihole-sync service stop`                       | Stops the running service                                          |
| `pihole-sync service restart`                    | Restarts the service                                                |
| `pihole-sync service status`                     | Shows whether the service is running                               |
| `pihole-sync service uninstall`                  | Removes the installed service                                      |

Installing/uninstalling/starting/stopping a service requires administrator/root privileges (`sudo` on Linux/macOS, an elevated shell on Windows). The `--env-file` is read by the running service, not baked into the service definition, so secrets stay out of the systemd unit/Windows registry — protect the file with restrictive permissions (e.g. `chmod 600`) instead.

## Configuration

The following environment variables can be specified:

### Required Environment Variables

| Name      | Default | Example                                          | Description                                              |
|-----------|---------|--------------------------------------------------|----------------------------------------------------------|
| `PRIMARY` | n/a     | `http://ph1.example.com\|password`                       | Specifies the primary Pi-hole configuration              |
| `REPLICAS`| n/a     | `http://ph2.example.com\|password,http://ph3.example.com\|password` | Specifies the list of replica Pi-hole configurations     |
| `FULL_SYNC` | n/a   | `true`                                           | Specifies whether to perform a full synchronization      |

> **Note:** When `FULL_SYNC=true`, the system will perform a full Teleporter import/export from the primary Pi-hole to the replicas. This will synchronize all settings and configurations.

> **Docker secrets:** `PRIMARY` and `REPLICAS` environment variables support Docker secrets when defined as `PRIMARY_FILE` and `REPLICAS_FILE`. See note regarding default user and Docker secrets example below.

### Optional Environment Variables

| Name                               | Default | Example         | Description                                        |
|------------------------------------|---------|-----------------|----------------------------------------------------|
| `CRON`                             | n/a     | `0 * * * *`     | Specifies the cron schedule for synchronization    |
| `RUN_GRAVITY`                      | false   | true            | Specifies whether to run gravity after syncing     |
| `TZ`                               | n/a     | `Europe/London` | Specifies the timezone for logs and cron           |
| `CLIENT_SKIP_TLS_VERIFICATION`     | false   | true            | Skips TLS certificate verification                 |
| `CLIENT_RETRY_DELAY_SECONDS`       | 1       | 5               | Seconds to delay between connection attempts       |
| `CLIENT_TIMEOUT_SECONDS`           | 20      | 60              | Http client timeout in seconds                     |
| `PRIMARY_TOTP`                     | n/a     | `JBSWY3DPEHPK3PXP` | TOTP secret for the primary, if it has 2FA enabled |
| `REPLICAS_TOTP`                    | n/a     | `secretA,,secretC` | TOTP secrets for the replicas, comma-separated in the same order as `REPLICAS`; leave an entry empty for replicas without 2FA |

> **2FA/TOTP:** If a Pi-hole instance has two-factor authentication enabled, set its TOTP secret (the same base32 seed used to set up your authenticator app) via `PRIMARY_TOTP`/`REPLICAS_TOTP` (or the `_FILE` variants, e.g. `PRIMARY_TOTP_FILE`) so pihole-sync can generate a valid code on each authentication request. Instances without 2FA don't need these set.

> **Note:** The following optional settings apply only if `FULL_SYNC=false`. They allow for granular control of synchronization if a full sync is not wanted.

| Name                              | Default | Description                            |
|-----------------------------------|---------|----------------------------------------|
| `SYNC_CONFIG_DNS`                  | false   | Synchronize DNS settings               |
| `SYNC_CONFIG_DHCP`                 | false   | Synchronize DHCP settings              |
| `SYNC_CONFIG_NTP`                  | false   | Synchronize NTP settings               |
| `SYNC_CONFIG_RESOLVER`             | false   | Synchronize resolver settings          |
| `SYNC_CONFIG_DATABASE`             | false   | Synchronize database settings          |
| `SYNC_CONFIG_MISC`                 | false   | Synchronize miscellaneous settings     |
| `SYNC_CONFIG_DEBUG`                | false   | Synchronize debug settings             |
| `SYNC_GRAVITY_DHCP_LEASES`         | false   | Synchronize DHCP leases                |
| `SYNC_GRAVITY_GROUP`               | false   | Synchronize groups                     |
| `SYNC_GRAVITY_AD_LIST`             | false   | Synchronize ad lists                   |
| `SYNC_GRAVITY_AD_LIST_BY_GROUP`    | false   | Synchronize ad lists by group          |
| `SYNC_GRAVITY_DOMAIN_LIST`         | false   | Synchronize domain lists               |
| `SYNC_GRAVITY_DOMAIN_LIST_BY_GROUP`| false   | Synchronize domain lists by group      |
| `SYNC_GRAVITY_CLIENT`              | false   | Synchronize clients                    |
| `SYNC_GRAVITY_CLIENT_BY_GROUP`     | false   | Synchronize clients by group           |


#### Config filters
> Allows including or excluding specific config keys.\
**Note:** `The SYNC_CONFIG_*_INCLUDE` and `SYNC_CONFIG_*_EXCLUDE` settings are mutually exclusive within each section. Additionally, config filters are only applied if `FULL_SYNC=false`.\
Config keys are relative to the section and are **case sensitive**. For example, the key `dns.upstreams` should be referred to as `upstreams`, and `dns.cache.size` should be referred to as `cache.size`.

| Name                              | Example                    | Description                                     |
|-----------------------------------|----------------------------|-------------------------------------------------|
| `SYNC_CONFIG_DNS_INCLUDE`         | upstreams,interface        | DNS config keys to include                     |
| `SYNC_CONFIG_DNS_EXCLUDE`         | upstreams,interface        | DNS config keys to exclude                     |
| `SYNC_CONFIG_DHCP_INCLUDE`        | active,start               | DHCP config keys to include                    |
| `SYNC_CONFIG_DHCP_EXCLUDE`        | active,start               | DHCP config keys to exclude                    |
| `SYNC_CONFIG_NTP_INCLUDE`         | ipv4,sync                  | NTP config keys to include                     |
| `SYNC_CONFIG_NTP_EXCLUDE`         | ipv4,sync                  | NTP config keys to exclude                     |
| `SYNC_CONFIG_RESOLVER_INCLUDE`    | resolveIPv4,networkNames   | Resolver config keys to include                |
| `SYNC_CONFIG_RESOLVER_EXCLUDE`    | resolveIPv4,networkNames   | Resolver config keys to exclude                |
| `SYNC_CONFIG_DATABASE_INCLUDE`    | DBimport,maxDBdays         | Database config keys to include                |
| `SYNC_CONFIG_DATABASE_EXCLUDE`    | DBimport,maxDBdays         | Database config keys to exclude                |
| `SYNC_CONFIG_MISC_INCLUDE`        | nice,delay_startup         | Misc config keys to include                    |
| `SYNC_CONFIG_MISC_EXCLUDE`        | nice,delay_startup         | Misc config keys to exclude                    |
| `SYNC_CONFIG_DEBUG_INCLUDE`       | database,networking        | Debug config keys to include                   |
| `SYNC_CONFIG_DEBUG_EXCLUDE`       | database,networking        | Debug config keys to exclude                   |

### Webhooks

Pihole Sync can invoke webhooks depending if a sync succeeded or failed. URL is required for the webhook to trigger. Both success and failure webhooks use the same enviroment variable pattern. Webhooks have a timeout of 10 seconds.

> **Note:** Replace `<OUTCOME>` with either `SUCCESS` or `FAILURE`.

| Name                                 | Default | Example                            | Description |
|--------------------------------------|---------|------------------------------------|-------------|
| `WEBHOOK_SYNC_<OUTCOME>_URL`            | n/a     | `https://www.example.com/webhook`  | URL to invoke for the webhook |
| `WEBHOOK_SYNC_<OUTCOME>_METHOD`         | `POST`  | `GET`                              | The HTTP method for the webhook |
| `WEBHOOK_SYNC_<OUTCOME>_BODY`           | n/a     | `this is my webhook body`          | The body of the webhook request |
| `WEBHOOK_SYNC_<OUTCOME>_HEADERS`        | n/a     | `header1:foo,header2:bar`          | HTTP headers to set for the webhook request in the format `key:value` separated by comma. Any whitespace will be used verbatim, no string trimming. |

Additionally, you can skip TLS verification for all webhooks if necessary:

| Name                                            | Default | Example         | Description                                        |
|-------------------------------------------------|---------|-----------------|----------------------------------------------------|
| `WEBHOOK_CLIENT_SKIP_TLS_VERIFICATION`     | false   | true            | Skips TLS certificate verification                 |

#### Integration examples:

##### healthcheck.io:

```
WEBHOOK_SYNC_SUCCESS_URL=https://hc-ping.com/{your-slug-or-guid-here}
WEBHOOK_SYNC_FAILURE_URL=https://hc-ping.com/{your-slug-or-guid-here}/fail
```

##### Apprise:

```
WEBHOOK_SYNC_FAILURE_URL=http://localhost:8080/notify
WEBHOOK_SYNC_FAILURE_BODY=urls=mailto://user:pass@gmail.com&body=test message
```

##### A service that needs JSON:

```
WEBHOOK_SYNC_FAILURE_URL=https://www.example.com/notify.json
WEBHOOK_SYNC_FAILURE_BODY={"hello":"world"}
WEBHOOK_SYNC_FAILURE_HEADERS=Content-Type:application/json
```

## Notes / Known issues

### Default user of Docker container / Docker secrets example
By default, the Docker container runs as user `1001`. If you are using Docker secrets, the user that is running the container will need read permissions to the files that the Docker secrets reference. If the user does not have the right permissions you will receive an error `Failed to initialize service error="open /run/secrets/primary: permission denied"`. To avoid this error, either make sure to `chown 1001 ./your/secretfiles && chmod 400 ./your/secretfiles` or use the [`user` directive in Docker Compose](https://docs.docker.com/reference/compose-file/services/#user) to change the user that the container runs as to a user of your choice - and then make sure to update your secret files' ownership to that user. In the example [docker-compose-with-secrets.yml](examples/docker-compose-with-secrets.yml), user `1234` owns `./secrets/primary.txt` and `./secrets/replicas.txt` and both have `-r--------` permissions.

### App passwords and authentication errors
When using Pi-hole's app passwords ("Configure app password" in the Web interface / API settings page) with pihole-sync, you should enable the Pi-hole setting `webserver.api.app_sudo` on your `REPLICAS` servers or you may receive authentication errors. To configure this setting, perform one of the following:
- From the command line, execute `sudo pihole-FTL --config webserver.api.app_sudo true`. Repeat for each replica.
- From the Pi-hole web UI, go to Settings -> All Settings. Toggle the "Modified settings / All settings" slider in the upper right to show "All settings". Choose the "Webserver and API" section. Check the "Enabled" box under `webserver.api.app_sudo` and then click "Save & Apply". Repeat for each replica.
- With the text editor of your choice, edit the `/etc/pihole/pihole.toml` file. Find the `app_sudo = false` line under `[webserver.api]` and modify it to read `app_sudo = true`. Save the file and restart the pihole-FTL service. Repeat for each replica.


## Disclaimer

This project is based on code of nebula-sync.

This project is an unofficial, community-maintained project and is not affiliated with the [official Pi-hole project](https://github.com/pi-hole). It aims to add sync/replication features not available in the core Pi-hole product but operates independently of Pi-hole LLC. Although tested across various environments, using any software from the Internet involves inherent risks. See the [license](https://github.com/Kizuren/pihole-sync/blob/main/LICENSE) for more details.

Pi-hole and the Pi-hole logo are [registered trademarks](https://pi-hole.net/trademark-rules-and-brand-guidelines) of Pi-hole LLC.


