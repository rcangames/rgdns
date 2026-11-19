# RGDNS architecture — first iteration
RGDNS is the DNS engineering project of Rcan Games. This repository begins with a local authoritative DNS prototype and a static company website.

## Prototype
- Load zone records from JSON.
- Listen on UDP and TCP.
- Support A, AAAA, CNAME, NS and TXT record types.
- Respond authoritatively for configured zones.

## Planned work
1. Registry data model and authenticated registration API.
2. DNSSEC signing, key rotation and access controls.
3. Monitoring, backups, audit logs and incident response.
4. Management dashboard and production deployment.

## Global namespace reality
A DNS server or GitHub repository cannot independently create a globally recognized `.rg` top-level domain. This requires registry arrangements and delegation through the global DNS root, plus legal, security and operational readiness.

The sample records use 192.0.2.0/24, a documentation-only address range. This server is a local prototype, not a public recursive resolver.
