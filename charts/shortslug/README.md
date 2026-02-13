# ShortSlug Helm Chart

Install:
```bash
helm install shortslug ./charts/shortslug \
  --set image.tag=v1.0.2
```

Values of interest:
- `image.repository`
- `image.tag`
- `env.BRAND_NAME`
- `persistence.enabled`
- `ingress.enabled`
