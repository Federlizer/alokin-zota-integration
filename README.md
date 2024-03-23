# Alokin - Zota integration

### Config

To configure the application, you will need to set the following environment variables in the `Dockerfile`:

```Dockerfile
# The secret API key provided by Zota
ENV ZOTA_SECRET_KEY=11111111-1111-1111-1111-111111111111
# The endpoint ID provided by Zota
ENV ZOTA_ENDPOINT_ID=111111
# The merchant ID provided by Zota
ENV ZOTA_MERCHANT_ID=EXAMPLE-MERCHANT-ID
# The base URL provided by Zota
ENV ZOTA_BASE_URL=https://api.zotapay-sandbox.com
```
