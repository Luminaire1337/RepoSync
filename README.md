# RepoSync

A minimalist Go-based HTTP webhook listener that automatically pulls updates from your Git repository.

## Overview

RepoSync is a lightweight webhook server designed to keep your Git repositories in sync with remote changes. When configured with a GitHub webhook, it automatically pulls the latest changes whenever commits are pushed to the repository.

## Features

- Minimal footprint and dependencies
- Secure webhook verification using HMAC-SHA256
- Simple configuration through environment variables

## Installation

```bash
git clone https://github.com/Luminaire1337/RepoSync.git
cd RepoSync
go build
```

## Configuration

RepoSync requires the following environment variables:

- `GITHUB_SECRET`: The secret key configured in your GitHub webhook
- `REPO_DIR`: The path to your Git repository directory
- `LISTEN_ADDR` (optional): The address and port to listen on (default: ":8080")

## Systemd Service Setup

To run RepoSync as a systemd service:

1. Build the executable:

   ```bash
   go build
   ```

2. Create a systemd service file:

   ```bash
   sudo nano /etc/systemd/system/reposync.service
   ```

3. Add the following content to the service file:

   ```
   [Unit]
   Description=RepoSync GitHub webhook server
   After=network.target

   [Service]
   Type=simple
   User=YOUR_USERNAME
   WorkingDirectory=/home/YOUR_USERNAME/RepoSync
   ExecStart=/home/YOUR_USERNAME/RepoSync/RepoSync
   Restart=always
   RestartSec=5
   Environment="GITHUB_SECRET=your_webhook_secret"
   Environment="REPO_DIR=/home/YOUR_USERNAME/DifferentRepo"

   [Install]
   WantedBy=multi-user.target
   ```

4. Enable and start the service:

   ```bash
   sudo systemctl daemon-reload
   sudo systemctl enable --now reposync
   ```

## GitHub Webhook Setup

Configure a webhook in your GitHub repository:

- Payload URL: `http://your-server.com:8080/webhook`
- Content type: `application/json`
- Secret: same value as your `GITHUB_SECRET`
- Events: Just the push event

## Security

RepoSync validates incoming webhook requests by verifying the HMAC-SHA256 signature provided by GitHub in the `X-Hub-Signature-256` header.

## License

See the [LICENSE](https://github.com/Luminaire1337/RepoSync/blob/master/LICENSE) file for details.
