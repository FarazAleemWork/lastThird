# WhenIsTahajjud

## What Is This App?
This app gives users the last third of the night which is a special time for muslims 

## Live Site (Possibly online Azure is expensive or you can spin it up with the docker image or with the main.tf & config/cloudinit files)

## Technologies Used 

- **Golang** - Backend API, Static file server and Fetch data from external API
- **Docker** - Containerized image (can be found on docker hub farazaleemwork)
- **Docker Compose** - Multicontainer orchestration
- **Terraform** - IaC for Azure VM, VN, Subnet, RG, IP provisioning
- **cloud-init & config** - Automated downloading of docker, images and running the containers
- **Nginx** - Reverse proxy server
- **Github Actions** - Currently pushes changes to docker hub and creates a new image version

## Technologies to Be Added
- **Prometheus**
- **Azure Monitor + Log Analytics**
