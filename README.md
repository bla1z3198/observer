# observer - hysteria2 auth & subscription server

<img src="https://img.shields.io/badge/System-Linux-FCC624?style=for-the-badge&logo=linux&logoColor=black" alt="Linux" /> <img src="https://img.shields.io/badge/Language-Go-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go" /> <img src="https://img.shields.io/badge/Hysteria-v2-FF4F00?style=for-the-badge&logo=fastapi&logoColor=white" alt="Hysteria 2"/>

A lightweight, standalone HTTP service written in Go for managing client authentication and subscription distribution for **Hysteria 2** proxy servers.

---

## Features

- **Hysteria 2 HTTP Auth API:** Handles incoming authentication verification requests directly from the Hysteria 2 daemon.
- **TLS Subscription Endpoint:** Distributes base64-encoded client subscription links over HTTPS compatible with clients supporting subscription profiles.
- **Minimal Resource Footprint:** Runs comfortably on micro VPS instances (as low as 1 vCPU and 512 MB – 1 GB RAM).
- **File-Driven Configuration:** Separates runtime configuration, user token maps, and client connection links into clean JSON and text files.

---
## How to start?
```
Download archive from releases

Create config.json, users.json and links.txt

Use the setup.sh (but check hy2 cfg carefully after) 
```
## Files
```
config.json - is a observer config
users.json - token and id of user (table)
links.txt - your hysteria2 links (auth via token)
```
## Example of hy2 config
```
listen: :443

tls:
  cert: /etc/letsencrypt/live/YOUR_DOMAIN_HERE/fullchain.pem
  key: /etc/letsencrypt/live/YOUR_DOMAIN_HERE/privkey.pem

masquerade:
  type: file
  file:
    dir: root/var/www/masq
  listenHTTPS: :443
  forceHTTPS: true                                              

auth:
  type: http
  http:
    url: http://ADDRESS OF OBSERVER (AUTH)/auth
```