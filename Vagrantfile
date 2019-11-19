Vagrant.configure("2") do |config|
  config.vm.box = "bento/ubuntu-18.04"
  config.vm.box_check_update = false

  config.vm.network "forwarded_port", guest: 8080, host: 8080
  config.vm.provider "virtualbox" do |vb|
    vb.memory = "1024"
  end

  config.vm.provision "shell", path: "install/support.sh"
  config.vm.provision "shell", path: "install/docker.sh"
  config.vm.provision "shell", path: "install/go.sh"
end