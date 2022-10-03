# plausible-exporter

<img width="1438" alt="image" src="https://user-images.githubusercontent.com/6919894/193658233-18ecc4a2-52c7-4c48-b409-d315a4a44c41.png">

`plausible-exporter` is a Prometheus exporter for [Plausible Analytics](https://plausible.io).
It enables you to keep an overview of your websites statistics in Prometheus.

## Usage

`plausible-exporter` can be run as a Docker container (preferrably) or built from source and run as a static binary.

### docker-compose

```yaml
version: "3"
services:
  plausible-exporter:
    image: ghcr.io/riesinger/plausible-exporter:latest
    environment:
      - PLAUSIBLE_HOST=https://plausible.io
      - PLAUSIBLE_SITE_IDS=your-web.site
      - PLAUSIBLE_TOKEN=<the-token-you-created-in-plausible>
    ports:
      - 8080:8080
```

For configuration see [this section](#configuration).

### Binary

To run the exporter as binary, clone this repo, build the Go binary and run it:

```sh
git clone https://github.com/riesinger/plausible-exporter
make # or "make static"
./plausible-exporter
```

For configuration see [this section](#configuration).

### Configuration

The exporter can be configured via a `config.yaml` file placed in `/etc/plausible-exporter` or via environment variables.

When put into a config file, the variable names are `snake_cased`, when set via the environment, the variables must be `UPPER_SNAKE_CASED`.
All options can be set in the config file or environment variables, with environment variables taking precedence.

| Option               | Required | Description                                                      | Default        |
| -------------------- | -------- | ---------------------------------------------------------------- | -------------- |
| `plausible_host`     | ✅       | The hostname and protocol of your plausible server               | -              |
| `plausible_site_ids` | ✅       | The IDs of the sites you want to fetch from plausible, as a list | -              |
| `plausible_token`    | ✅       | A valid API token for your plausible server                      | -              |
| `listen_address`     | ❌       | Which host and port to listen to                                 | `0.0.0.0:8080` |

### Prometheus

This exporter creates 4 metrics:

* `plausible_visitors` - How many visitors were on your site on the current day
* `plausible_visit_duration` - How long an average visit to your site was (in seconds)
* `plausible_page_views` - How many page views your site had today
* `plausible_bounce_rate` - How many visitors left your site (in percent, 0-100)

In case you've configured multiple sites to be scraped, you can differentiate between them with the `site_id` label.

## License

This project is MIT-licensed, see [LICENSE.md](./LICENSE.md)
