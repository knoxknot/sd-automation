Vagrant.configure("2") do |config|
  # Operating Systen for the VM
  config.vm.box = "ubuntu/xenial64"

  # SSH Settings
  config.ssh.private_key_path = ["~/.ssh/csproject_key", "~/.vagrant.d/insecure_private_key"]
  config.ssh.insert_key = false

  # Upload Public Key into the VM
  config.vm.provision "file", source: "~/.ssh/csproject_key.pub", destination: "~/.ssh/authorized_keys"

  # Configure the VM with Hostname Master
  config.vm.define "master" do |master|
    master.vm.hostname = "master"
    master.vm.network "private_network", ip: "192.168.255.6"
    master.vm.provision :ansible do |ansible|
      ansible.inventory_path = "./infrastructure/config/hosts"
      ansible.verbose = "vvvv"
      ansible.raw_arguments  = ["--private-key=~/.ssh/csproject_key"]
      ansible.playbook = "./infrastructure/config/configure.yaml"
    end
    master.vm.provision "shell", path: "./infrastructure/config/bootstrap-cluster.sh" 
  end 
 
  # Synchronize Folders within the Host and VM
  config.vm.synced_folder ".", "/home/vagrant/csproject", :owner => 'vagrant', :mount_options => ["dmode=774", "fmode=774"]
 
  # VM Provider
  config.vm.provider "virtualbox" do |vb|
    # Disconnect serial mode on bootup
    vb.customize [ "modifyvm", :id, "--uartmode1", "disconnected" ]

    # Group the Machines
    vb.customize ["modifyvm", :id, "--groups", "/CSProject Cluster"]

    # Customize the number CPUS on the VM
    vb.cpus = "2"

    # Customize the amount of memory on the VM
    vb.memory = "1024"
  end
end