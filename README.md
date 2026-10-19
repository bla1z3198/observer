# Hysteria 2 Auth & Subscription Server

A lightweight, standalone HTTP service written in Go for managing client authentication and subscription distribution for **Hysteria 2** proxy servers.

---

## Features

- **Hysteria 2 HTTP Auth API:** Handles incoming authentication verification requests directly from the Hysteria 2 daemon.
- **TLS Subscription Endpoint:** Distributes base64-encoded client subscription links over HTTPS compatible with clients supporting subscription profiles.
- **Minimal Resource Footprint:** Runs comfortably on micro VPS instances (as low as 1 vCPU and 512 MB – 1 GB RAM).
- **File-Driven Configuration:** Separates runtime configuration, user token maps, and client connection links into clean JSON and text files.

---

## Architecture Overview

```text
               +---------------------------+
               |     Connecting Client     |
               +-------------+-------------+
                             | 
                             | 1. Connects via QUIC
                             v
               +-------------+-------------+
               |      Hysteria 2 Node      |
               +-------------+-------------+
                             |
                             | 2. POST /auth (HTTP)
                             v
               +-------------+-------------+
               |    Go Auth/Sub Server     |
               |                           |
               | - /auth (Hysteria check)  |
               | - /sub/ (Client links)    |
               +-------------+-------------+