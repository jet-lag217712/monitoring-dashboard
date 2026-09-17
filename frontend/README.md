# Equate Monitoring Dashboard

React 19 + Vite frontend for the local Equate appliance. The browser talks to
nginx, nginx proxies approved requests to the local Backend API, and the API
reads PostgreSQL. The frontend never connects directly to collectors, MQTT, or
PostgreSQL.

## Local development

```bash
cp .env.example .env.local
npm install
npm run dev
```

Open the Vite URL. Sign-in uses the appliance-local username and password
backed by the host PAM broker.

## Environment

| Variable | Purpose |
|---|---|
| `VITE_API_BASE_URL` | API origin; appliance default is `/api` |
| `VITE_AUTH_MODE` | Must be `appliance_local` |
| `VITE_APP_VERSION` | Release version shown in the footer; same value as `equate version` (`buildVersion`). Appliance builds pass the release `--version`. Local/dev defaults to `git describe`. |

## Production image

```bash
docker build \
  --build-arg VITE_API_BASE_URL=/api \
  --build-arg VITE_AUTH_MODE=appliance_local \
  --build-arg VITE_APP_VERSION=1.4.0 \
  -t equate/frontend:local .
```

The production appliance serves the built frontend from nginx over its local
HTTPS endpoint. No external identity provider is part of the supported
appliance configuration.
