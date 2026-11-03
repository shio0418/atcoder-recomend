terraform {
  required_providers {
    vultr = {
      source  = "vultr/vultr"
      version = "~> 2.19"
    }
    cloudflare = {
        source = "cloudflare/cloudflare"
        version = "~> 5.0"
    }
  }
}

variable "vultr_api_key" {
  type      = string
  sensitive = true
}

provider "vultr" {
  api_key = var.vultr_api_key
}

resource "vultr_ssh_key" "my_key" {
  name    = "my-laptop"
  ssh_key = file("~/.ssh/id_ed25519.pub") # 自分の公開鍵のパスに合わせる
}

data "vultr_os" "ubuntu" {
  filter {
    name   = "name"
    values = ["Ubuntu 22.04 LTS x64"]
  }
}

resource "vultr_startup_script" "setup" {
  name = "atcoder-recomend-setup"
  script = filebase64("setup.sh")
  type = "boot"
}

resource "vultr_instance" "app_server" {
  plan        = "vc2-2c-4gb" # ← ダッシュボードで実際のplan IDを確認して合わせる
  region      = "nrt"        # Tokyo
  os_id       = data.vultr_os.ubuntu.id
  label       = "atcoder-recomend"
  hostname    = "atcoder-recomend"
  ssh_key_ids = [vultr_ssh_key.my_key.id]
  script_id = vultr_startup_script.setup.id
}

output "instance_ip" {
  value = vultr_instance.app_server.main_ip
}

variable "cloudflare_api_token" {
  type      = string
  sensitive = true
}

provider "cloudflare" {
  api_token = var.cloudflare_api_token
}

data "cloudflare_zone" "this" {
  filter = {
    name = "shio0418.dev"
  }
}

resource "cloudflare_dns_record" "app" {
  zone_id = data.cloudflare_zone.this.id
  name = "atcoder-recommend"
  type = "A"
  ttl = 1
  content = vultr_instance.app_server.main_ip
  proxied = true
}