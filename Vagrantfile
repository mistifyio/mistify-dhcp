
VAGRANTFILE_API_VERSION = "2"

Vagrant.configure(VAGRANTFILE_API_VERSION) do |config|
	config.vm.box = "mistify/trusty64-vmware"
	config.vm.box_url  = "http://www.akins.org/boxes/mistify-ubuntu-vmware.box"
	config.ssh.forward_agent = true

	config.vm.synced_folder ".", "/home/vagrant/go/src/github.com/mistifyio/mistify-dhcp", create: true

    #config.vm.network "private_network", type: "dhcp"

    config.vm.network "private_network",auto_config: false

	config.vm.provider "vmware_fusion" do |v|
		# GUI is needed because when you use bridging inside Linux,
		# Fusion must ask for admin password
		v.gui = true
		v.vmx["memsize"] = "1024"
		v.vmx["numvcpus"] = "2"
		v.vmx["vhv.enable"] = "TRUE"
	end

    config.vm.provision "shell", privileged: false, inline: <<EOF
cd /home/vagrant/go/src/github.com/mistifyio/mistify-dhcp
go get github.com/tools/godep
godep go install
sudo apt-get update
sudo apt-get install -y dhcping
bash test.sh
EOF
end
