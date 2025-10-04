terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "4.46.0"
    }
  }
}

provider "azurerm" {
  # Configuration options
  features {}
  subscription_id = "baa15b7f-fff5-47dd-b2d4-adaf10dd077b"
}

resource "azurerm_resource_group" "tahajjud-app-rg" {
  name     = "tahajjud-app-resources"
  location = "East US 2"
  tags = {
    environment = "test"
  }
}

resource "azurerm_virtual_network" "tahajjud-vn" {
  name                = "tahajjud-network"
  resource_group_name = azurerm_resource_group.tahajjud-app-rg.name
  location            = azurerm_resource_group.tahajjud-app-rg.location
  address_space       = ["10.10.10.0/24"]

  tags = {
    environment = "test"
  }
}

resource "azurerm_subnet" "tahajjud-sn" {
  name                 = "tahajjud-subnet"
  resource_group_name  = azurerm_resource_group.tahajjud-app-rg.name
  virtual_network_name = azurerm_virtual_network.tahajjud-vn.name

  address_prefixes = [
    "10.10.10.0/25",
    "10.10.10.128/25"
  ]
}

resource "azurerm_network_security_group" "tahajjud-nsg" {
  name                = "tahajjud-network-security-group"
  location            = azurerm_resource_group.tahajjud-app-rg.location
  resource_group_name = azurerm_resource_group.tahajjud-app-rg.name

  tags = {
    environment = "test"
  }
}

resource "azurerm_network_security_rule" "tjrule1" {
  name                        = "tjrule1-AllowHTTP"
  priority                    = 100
  direction                   = "Inbound"
  access                      = "Allow"
  protocol                    = "Tcp"
  source_port_range           = "*"
  destination_port_range      = "80"
  source_address_prefix       = "*"
  destination_address_prefix  = "*"
  resource_group_name         = azurerm_resource_group.tahajjud-app-rg.name
  network_security_group_name = azurerm_network_security_group.tahajjud-nsg.name
}

resource "azurerm_network_security_rule" "tjrule2" {
  name                        = "tjrule2-AllowHTTPS"
  priority                    = 101
  direction                   = "Inbound"
  access                      = "Allow"
  protocol                    = "Tcp"
  source_port_range           = "*"
  destination_port_range      = "443"
  source_address_prefix       = "*"
  destination_address_prefix  = "*"
  resource_group_name         = azurerm_resource_group.tahajjud-app-rg.name
  network_security_group_name = azurerm_network_security_group.tahajjud-nsg.name
}

resource "azurerm_network_security_rule" "tjrule3" {
  name                        = "tjrule2-AllowSSH"
  priority                    = 102
  direction                   = "Inbound"
  access                      = "Allow"
  protocol                    = "*"
  source_port_range           = "*"
  destination_port_range      = "22"
  source_address_prefix       = "*"
  destination_address_prefix  = "*"
  resource_group_name         = azurerm_resource_group.tahajjud-app-rg.name
  network_security_group_name = azurerm_network_security_group.tahajjud-nsg.name
}

resource "azurerm_subnet_network_security_group_association" "tj-nsg-assoc" {
  subnet_id                 = azurerm_subnet.tahajjud-sn.id
  network_security_group_id = azurerm_network_security_group.tahajjud-nsg.id
}

resource "azurerm_public_ip" "tj-ip" {
  name                = "tahajjud-public-ip"
  resource_group_name = azurerm_resource_group.tahajjud-app-rg.name
  location            = azurerm_resource_group.tahajjud-app-rg.location
  allocation_method   = "Static"

  tags = {
    environment = "test"
  }
}

resource "azurerm_network_interface" "tj-nic" {
  name                = "tahajjud-nic"
  location            = azurerm_resource_group.tahajjud-app-rg.location
  resource_group_name = azurerm_resource_group.tahajjud-app-rg.name

  ip_configuration {
    name                          = "internal"
    subnet_id                     = azurerm_subnet.tahajjud-sn.id
    private_ip_address_allocation = "Dynamic"
    public_ip_address_id          = azurerm_public_ip.tj-ip.id

  }

  tags = {
    environment = "test"
  }
}

resource "azurerm_linux_virtual_machine" "tj-linux-vm" {
  name                  = "tahajjud-vm"
  location              = azurerm_resource_group.tahajjud-app-rg.location
  resource_group_name   = azurerm_resource_group.tahajjud-app-rg.name
  size                  = "Standard_B1s"
  admin_username        = "adminuser"
  network_interface_ids = [azurerm_network_interface.tj-nic.id]


  custom_data = base64encode(templatefile("cloud-init.yml", {
    compose_content = file("${path.module}/configs/docker-compose.yml"),
    nginx_content   = file("${path.module}/configs/nginx.conf")
  }))

  admin_ssh_key {
    username   = "adminuser"
    public_key = file("~/.ssh/tjazurekey.pub")
  }

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
    version   = "latest"
  }

  provision_vm_agent = true

}