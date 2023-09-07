variable "swarm-manager-host" {
  type        = string
  sensitive   = true
  description = "Address of swarm manager"
}

variable "app-name" {
  type        = string
  description = "Name of app"
}

variable "app-version" {
  type        = string
  description = "Version of Docker image of app"
  default     = "latest"
}

variable "app-host" {
  type        = string
  description = "Hostname of app"
  default     = "gomvn.example.com"
}

variable "app-config-b64" {
  type        = string
  description = "Application config file, base64-encoded"
  sensitive   = true
}

variable "cpu-limit" {
  type        = number
  description = "CPU limit in nanoseconds"
  default     = 100 * 10000000
}

variable "memory-limit" {
  type        = number
  description = "Memory limit in bytes"
  default     = 64 * 1024 * 1024
}
