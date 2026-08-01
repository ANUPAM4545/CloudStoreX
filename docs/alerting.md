# CloudStoreX AlertManager Configuration & Alerting Guide

## 1. Alert Severity Hierarchy

CloudStoreX defines four distinct alert severity levels in `deploy/observability/alertmanager/rules/alert.rules.yml`:

| Severity | Target Response Time | Notification Channel | Examples |
| :--- | :--- | :--- | :--- |
| **Critical** | < 15 minutes (24/7 On-Call) | PagerDuty / VictorOps + Slack (`#sre-incidents`) | `StorageProviderUnreachable`, `HTTPHighErrorRateCritical`, `PostgresDatabaseUnreachable` |
| **High** | < 1 hour | Slack (`#sre-alerts`) + Jira Incident | `StorageUploadLatencyHigh`, `BackgroundWorkerQueueDepthHigh`, `ProviderErrorRateHigh` |
| **Medium** | < 24 hours (Business Hours) | Slack (`#platform-warnings`) | `StorageQuotaThresholdWarning`, `PolicyEngineFallbackRateHigh` |
| **Low** | Weekly Review | Email / Dashboard Only | `MetadataCatalogGrowthWarning` |

---

## 2. Notification Routing & Grouping

In `deploy/observability/alertmanager/alertmanager.yml`, notifications are routed and deduplicated:

```yaml
route:
  group_by: ['alertname', 'workspace', 'provider']
  group_wait: 30s
  group_interval: 5m
  repeat_interval: 4h
  receiver: 'slack-sre-alerts'
  routes:
    - match:
        severity: critical
      receiver: 'pagerduty-critical'
      continue: true
```

- **Deduplication**: Multiple alerts for the same `alertname` and `provider` within 30 seconds are merged into a single alert notification.
- **Repeat Interval**: Ongoing firing alerts re-notify after 4 hours.

---

## 3. Creating Alert Silences

During planned maintenance or database rollouts, create a temporary silence via AlertManager CLI/API:

```bash
amtool silence add alertname=PostgresDatabaseUnreachable --author="SRE On-Call" --comment="Planned PG 16 upgrade" --duration=2h
```
