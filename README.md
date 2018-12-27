## Simple Test App (v2)

A simple web service that allow a quick test of micro-service deployment and routing behaviors. It exposes two interfaces:
- root interface (`/`) returns an HTML page and triggers loading of the information about the host and the application
- information interface (`/info`) returns an HTML table with the information about the host and the application
- health interface (`/healthz`) that returns OK once the application start running

The usage:

```bash
simple-test-app <port>
```

If `<port>` is omitted, the application starts on port `8282`. The title of the application can be customized by defining environment variable `TITLE`. If the variable isn't defined, the application title is `Simple Test Application (X)` where `X` stands for hardcoded code version.
