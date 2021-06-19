### i18n rewrite

[![Main](https://github.com/WLun001/i18n-rewrite/actions/workflows/main.yml/badge.svg)](https://github.com/WLun001/i18n-rewrite/actions/workflows/main.yml)
[![Go Matrix](https://github.com/WLun001/i18n-rewrite/actions/workflows/go-cross.yml/badge.svg)](https://github.com/WLun001/i18n-rewrite/actions/workflows/go-cross.yml)

A [Traefik](https://traefik.io) middleware plugin that rewrite path based on `Accept-Language` request header

### Development
```bash
# bash shell
TRAEFIK_PILOT_TOKEN=your-token docker-compose -f docker-compose.dev.yaml up

# fish shell
env TRAEFIK_PILOT_TOKEN=your-token docker-compose -f docker-compose.dev.yaml up
```
