# Camunda History Viewer

A lightweight Go web application created as a supplementary tool for my main [Camunda LLM Ticket Refund](https://github.com/DigitLock/camunda-llm-ticket-refund) portfolio project. Built to provide quick, visual access to process execution history—making it easier to demonstrate and debug BPMN workflows without navigating through Camunda Cockpit's interface.

**Purpose:**
- 🎯 Supporting tool for [Camunda LLM Ticket Refund](https://github.com/DigitLock/camunda-llm-ticket-refund) project
- 📊 Enables better presentation of process execution flows
- 🔍 Simplifies debugging during development
- 🎓 Demonstrates Go backend development and REST API integration

## Features

- 📊 View last 10 process instances
- 🔍 Detailed activity timeline for each process
- ✅ Visual status indicators (success/failure)
- 🎨 Clean, responsive UI
- 🔐 Basic authentication support

## Prerequisites

- Go 1.21+
- Camunda Platform 7.x with REST API enabled and Basic Authentication configured

## Camunda Configuration

Before running the application, you need to enable Basic Authentication for Camunda REST API.

### Enable Basic Auth in Camunda

1. Access your Camunda container or installation directory
2. Edit the file: `camunda/webapps/engine-rest/WEB-INF/web.xml`
3. Uncomment the HTTP Basic Authentication filter section:
```xml
<!-- Http Basic Authentication Filter -->
<filter>
    <filter-name>camunda-auth</filter-name>
    <filter-class>
      org.camunda.bpm.engine.rest.security.auth.ProcessEngineAuthenticationFilter
    </filter-class>
    <async-supported>true</async-supported>
    <init-param>
      <param-name>authentication-provider</param-name>
      <param-value>org.camunda.bpm.engine.rest.security.auth.impl.HttpBasicAuthenticationProvider</param-value>
    </init-param>
    <init-param>
      <param-name>rest-url-pattern-prefix</param-name>
      <param-value></param-value>
    </init-param>
</filter>

<filter-mapping>
    <filter-name>camunda-auth</filter-name>
    <url-pattern>/*</url-pattern>
</filter-mapping>
```

4. Restart Camunda:
```bash
docker restart camunda  # if using Docker
# or restart Tomcat service
```

5. Verify authentication is working:
```bash
curl -i -u demo:demo "http://your-camunda-host:8080/engine-rest/engine"
```

You should receive a `200 OK` response with engine information.

## Installation

1. Clone the repository:
```bash
git clone https://github.com/yourusername/camunda-history-viewer.git
cd camunda-history-viewer
```

2. Copy `.env.example` to `.env` and configure:
```bash
cp .env.example .env
```

3. Edit `.env` with your Camunda settings:
```env
CAMUNDA_BASE_URL=http://your-camunda-host:8080/engine-rest
CAMUNDA_USER=your_username
CAMUNDA_PASSWORD=your_password
SERVER_PORT=3000
```

4. Install dependencies:
```bash
go mod download
```

5. Run the application:
```bash
go run main.go
```

6. Open your browser at `http://localhost:3000`

## Configuration

All configuration is done via environment variables in the `.env` file:

| Variable | Description | Default |
|----------|-------------|---------|
| `CAMUNDA_BASE_URL` | Camunda REST API endpoint | `http://localhost:8080/engine-rest` |
| `CAMUNDA_USER` | Camunda username | `demo` |
| `CAMUNDA_PASSWORD` | Camunda password | `demo` |
| `SERVER_PORT` | Port for the web server | `3000` |

## Project Structure
```
.
├── main.go              # Main application
├── templates/
│   ├── home.html        # Process list page
│   └── process.html     # Process detail page
├── .env.example         # Example environment variables
├── .gitignore
├── go.mod
└── README.md
```

## Troubleshooting

### Authentication Error (401 Unauthorized)
- Verify Basic Auth is enabled in Camunda's `web.xml`
- Check username and password in `.env` file
- Ensure Camunda container/service has been restarted after config changes

### Connection Refused
- Verify `CAMUNDA_BASE_URL` is correct
- Check if Camunda is running: `curl http://your-camunda-host:8080/camunda/`
- Ensure network connectivity between the app and Camunda

### Empty Process List
- Verify the `processDefinitionKey` in code matches your deployed process
- Check if there are any completed process instances in Camunda Cockpit

## Related Projects

- [Camunda LLM Ticket Refund](https://github.com/DigitLock/camunda-llm-ticket-refund) - Main portfolio project demonstrating BPMN process automation with LLM integration

## License

MIT