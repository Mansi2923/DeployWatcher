terraform {
  required_providers {
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 2.0"
    }
  }
}

provider "kubernetes" {
  config_path = "~/.kube/config"
}

# Namespace for DeployWatch
resource "kubernetes_namespace" "deploywatch" {
  metadata {
    name = "deploywatch"
  }
}

# Backend API Deployment
resource "kubernetes_deployment" "backend" {
  metadata {
    name      = "deploywatch-backend"
    namespace = kubernetes_namespace.deploywatch.metadata[0].name
  }

  spec {
    replicas = 2

    selector {
      match_labels = {
        app = "deploywatch-backend"
      }
    }

    template {
      metadata {
        labels = {
          app = "deploywatch-backend"
        }
      }

      spec {
        container {
          image = "deploywatch-backend:latest"
          name  = "backend"

          port {
            container_port = 8080
          }

          resources {
            limits = {
              cpu    = "0.5"
              memory = "512Mi"
            }
            requests = {
              cpu    = "250m"
              memory = "256Mi"
            }
          }
        }
      }
    }
  }
}

# Frontend Deployment
resource "kubernetes_deployment" "frontend" {
  metadata {
    name      = "deploywatch-frontend"
    namespace = kubernetes_namespace.deploywatch.metadata[0].name
  }

  spec {
    replicas = 2

    selector {
      match_labels = {
        app = "deploywatch-frontend"
      }
    }

    template {
      metadata {
        labels = {
          app = "deploywatch-frontend"
        }
      }

      spec {
        container {
          image = "deploywatch-frontend:latest"
          name  = "frontend"

          port {
            container_port = 80
          }

          resources {
            limits = {
              cpu    = "0.5"
              memory = "512Mi"
            }
            requests = {
              cpu    = "250m"
              memory = "256Mi"
            }
          }
        }
      }
    }
  }
}

# Backend Service
resource "kubernetes_service" "backend" {
  metadata {
    name      = "deploywatch-backend"
    namespace = kubernetes_namespace.deploywatch.metadata[0].name
  }

  spec {
    selector = {
      app = "deploywatch-backend"
    }

    port {
      port        = 80
      target_port = 8080
    }

    type = "ClusterIP"
  }
}

# Frontend Service
resource "kubernetes_service" "frontend" {
  metadata {
    name      = "deploywatch-frontend"
    namespace = kubernetes_namespace.deploywatch.metadata[0].name
  }

  spec {
    selector = {
      app = "deploywatch-frontend"
    }

    port {
      port        = 80
      target_port = 80
    }

    type = "ClusterIP"
  }
}

# Ingress for external access
resource "kubernetes_ingress_v1" "deploywatch" {
  metadata {
    name      = "deploywatch-ingress"
    namespace = kubernetes_namespace.deploywatch.metadata[0].name
    annotations = {
      "nginx.ingress.kubernetes.io/rewrite-target" = "/"
    }
  }

  spec {
    rule {
      host = "deploywatch.local"
      http {
        path {
          path = "/api"
          backend {
            service {
              name = kubernetes_service.backend.metadata[0].name
              port {
                number = 80
              }
            }
          }
        }
        path {
          path = "/"
          backend {
            service {
              name = kubernetes_service.frontend.metadata[0].name
              port {
                number = 80
              }
            }
          }
        }
      }
    }
  }
} 