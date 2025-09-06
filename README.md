# 🐸 FileHopper

**FileHopper** is a modern, lightweight, and cross-platform file server written in **Go**.
It provides **SFTP** and **WebDAV** support out of the box, comes with a **web-based admin UI**, and features a **secure QR-code based setup flow** for first-time configuration.

> _"The lightweight file server that leaps ahead."_

---

## ✨ Features

- 🔒 **Secure by default** – SFTP (via SSH) and WebDAV over HTTPS
- 🖥️ **Web-based Admin UI** – manage users, sessions, and settings
- 📱 **QR-Code Setup Flow** – scan and configure your server in minutes
- 📂 **User Management** – per-user home directories, roles, and permissions
- 🪶 **Lightweight Deployment** – single Go binary or Docker container
- 🖧 **Cross-Platform** – works on Linux, Windows, and macOS
- 🔌 **Extensible** – future support for S3/MinIO backends, quotas, and more

---

## 🚀 Quick Start

### Run with Go

```bash
go run ./cmd/filehopper
```

### Run with Docker

```bash
docker run -d \
  -p 22:22 \
  -p 443:443 \
  -v /srv/filehopper:/data \
  filehopper/filehopper:latest
```

After the first start, FileHopper will generate a **setup token** and display a **QR code**.
Scan it with your phone or open the printed URL in your browser to finish the initial configuration.

---

## 🛠️ Roadmap

### MVP

- [ ] SFTP server (user/password auth, chroot home dirs)
- [ ] WebDAV server (Basic Auth, HTTPS)
- [ ] WebUI (basic user management)
- [ ] QR-code based setup flow

### Phase 2

- [ ] TLS via Let's Encrypt
- [ ] Quotas & storage limits
- [ ] Audit logs & active session overview
- [ ] Drag & drop file uploads in WebUI

### Phase 3

- [ ] External storage backends (S3, MinIO, CephFS)
- [ ] LDAP/AD integration
- [ ] Sharing links & versioning
- [ ] REST API for automation

---

## 🧑‍💻 Contributing

Contributions are welcome!
Please open an issue or submit a pull request if you’d like to help improve FileHopper.

---

## 📜 License

MIT License – feel free to use, modify, and distribute.

---

## 🐸 About

FileHopper is designed to be a **modern alternative to traditional file servers**,
with a focus on **simplicity, security, and user experience**.
