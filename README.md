# Micro-Pinger

Micro-Pinger is a service written in Go (Golang) for monitoring various HTTP endpoints. It periodically sends requests to configured URLs and checks for expected responses. This documentation provides an overview of the main features and usage instructions.

## Features

- **HTTP Endpoint Monitoring**: Micro-Pinger monitors HTTP endpoints by sending requests according to the specified configurations.
- **Alerting Mechanism**: It provides customizable alerting mechanisms via various channels such as Slack, Telegram, etc., in case of failures or anomalies.
- **Configuration via YAML**: Configuration for services, including endpoints to monitor, alerting settings, request details, etc., can be provided using YAML files.
- **Throttling and Timeout**: The service is equipped with throttling and timeout mechanisms to ensure efficient resource usage and timely response.

## Usage

### Installation

Micro-Pinger can be installed and executed using the following steps:

1. Clone the repository:

   ```bash
   git clone https://github.com/jtrw/micro-pinger.git
   ```

2. Navigate to the project directory:

   ```bash
   cd micro-pinger
   ```

3. Build the executable:

   ```bash
   go build
   ```

4. Run the executable:

   ```bash
   ./micro-pinger
   ```

### Command-Line Options

Micro-Pinger supports the following command-line options:

- `-c, --config`: Path to the configuration file (default: config.yml).
- `-l, --listen`: Address to listen on (default: :8080).
- `-s, --secret`: Secret key for authentication (default: 123).
- `--pinsize`: Size of the PIN (default: 5).
- `--expire`: Maximum lifetime for a service (default: 24h).
- `--pinattempts`: Maximum attempts to enter PIN (default: 3).
- `--web`: Web UI location (default: /).

### Configuration

Configuration for Micro-Pinger is done via YAML files. Below is an example configuration format:

```yaml
services:
  - name: ExampleService
    url: http://example.com
    method: GET
    type: json
    body: ""
    interval: 5s
    headers:
      - name: Authorization
        value: Bearer TOKEN
    response:
      status: 200
      body: "Example Response"
      compare: contains
    alerts:
      - name: SlackAlert
        type: slack
        webhook: https://slack/webhook
        to: "#alerts"
        failure: 3
        success: 3
        send-on-resolve: true
```

### API Endpoints

Micro-Pinger exposes the following API endpoints:

- `/api/v1/check`: Initiates checks for configured services.

### Web Interface

Micro-Pinger provides a simple web interface accessible at the root URL. This interface can be customized using the --web option.
Now it's not implemented yet.

### Alerting Mechanism

Micro-Pinger sends alerts via configured channels (e.g., Slack, Telegram) in case of failures or anomalies in the monitored services.

#### Supports Messages

- Slack
- Telegram

### Example

Below is an example of using Micro-Pinger to monitor a service:

```bash
./micro-pinger --config=config.yml --listen=:8080 --secret=123
```

This command starts the Micro-Pinger service using the provided configuration file (config.yml), listening on port 8080, and with a secret key for authentication.
