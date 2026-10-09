#!/usr/bin/env bash
set -euo pipefail

# Este script se ejecuta en el runner Linux de CI, antes de crear BuildKit.
# Conserva las otras opciones del daemon y prioriza la caché pública de Google.
sudo python3 - <<'PY'
import json
from pathlib import Path

config_path = Path("/etc/docker/daemon.json")
config = json.loads(config_path.read_text()) if config_path.exists() else {}
mirror = "https://mirror.gcr.io"
previous_mirrors = config.get("registry-mirrors", [])
config["registry-mirrors"] = [mirror] + [url for url in previous_mirrors if url != mirror]
config_path.parent.mkdir(parents=True, exist_ok=True)
config_path.write_text(json.dumps(config, indent=2) + "\n")
PY

sudo systemctl restart docker
