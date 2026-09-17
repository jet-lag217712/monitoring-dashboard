# Update channel (Azure)

Static HTTPS hosting for Equate appliance `.eqa` packages and channel manifests.
See [`docs/releases/appliance-updates.md`](../../docs/releases/appliance-updates.md).

Examples:

- [`examples/manifest.stable.json`](examples/manifest.stable.json)
- [`examples/update-channel.conf`](examples/update-channel.conf)
- Schema: [`channel-manifest.schema.json`](channel-manifest.schema.json)

Publish:

```bash
make appliance-publish-azure STORAGE_ACCOUNT=<account> VERSION=<semver> ARCH=amd64
```
