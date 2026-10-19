#!/bin/bash

# Ensure script is run as root
if [ "$EUID" -ne 0 ]; then
  echo "Please run as root"
  exit 1
fi

echo "=== 1) Installing Hysteria 2 ==="
bash <(curl -fsSL https://get.hy2.sh/)

echo "=== 2) Setting up masquerade directory ==="
mkdir -p /var/www/masq
if [ -f "index.html" ]; then
  cp index.html /var/www/masq/
  echo "index.html moved to /var/www/masq/"
else
  echo "Warning: index.html not found in current directory. Please move it manually later."
fi

echo "=== 3) Configuring Hysteria 2 TLS ==="
echo "How would you like to configure TLS for Hysteria?"
echo "1) Auto-TLS (ACME / Let's Encrypt - requires a valid domain pointed to this server)"
echo "2) Custom Certificates (Provide paths manually)"
read -p "Select [1/2]: " tls_choice

HY2_CFG="/etc/hysteria/config.yaml"
cat <<EOF >$HY2_CFG
listen: :443

auth:
  type: http
  http:
    url: http://127.0.0.1:8080/auth

masquerade:
  type: file
  file:
    dir: /var/www/masq
EOF

if [ "$tls_choice" == "1" ]; then
  read -p "Enter your domain name (e.g., example.com): " domain_name
  read -p "Enter your email for Let's Encrypt: " acme_email

  # Append ACME TLS block
  cat <<EOF >>$HY2_CFG

acme:
  domains:
    - $domain_name
  email: $acme_email
EOF
else
  read -p "Enter path to TLS Certificate (e.g., /etc/hysteria/server.crt): " cert_path
  read -p "Enter path to TLS Key (e.g., /etc/hysteria/server.key): " key_path

  # Append Custom TLS block
  cat <<EOF >>$HY2_CFG

tls:
  cert: $cert_path
  key: $key_path
EOF
fi

echo "=== 4) Creating Go Auth Server .service file ==="
# Assuming the Go binary is compiled as 'hy2-auth' and exists in the current directory
cp hy2-auth /usr/local/bin/hy2-auth 2>/dev/null || echo "Warning: hy2-auth binary not found in current directory."
cp config.json users.json links.txt /usr/local/bin/ 2>/dev/null

cat <<EOF >/etc/systemd/system/hy2-auth.service
[Unit]
Description=Hysteria 2 Go Auth and Sub Service
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=/usr/local/bin
ExecStart=/usr/local/bin/hy2-auth
Restart=always
RestartSec=5s
LimitNOFILE=65535

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload

echo "=== 5) Enabling and starting services ==="
systemctl enable --now hy2-auth.service
systemctl enable --now hysteria-server.service

echo "=== 6) Done! ==="
echo "Check Hysteria 2 status: systemctl status hysteria-server"
echo "Check Auth Server status: systemctl status hy2-auth"
