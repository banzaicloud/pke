$enable_serial_logging = false

raise "vagrant-vbguest plugin must be installed" unless Vagrant.has_plugin? "vagrant-vbguest"

Vagrant.configure("2") do |config|
    # Sync time with the local host
    config.vm.provider 'virtualbox' do |vb|
        vb.customize [ "guestproperty", "set", :id, "/VirtualBox/GuestAdd/VBoxService/--timesync-set-threshold", 1000 ]
    end

    # sync build folder
    config.vm.synced_folder '.', '/vagrant', disabled: true
    config.vm.synced_folder 'scripts/vagrant/', '/scripts/', create: true
    config.vm.synced_folder 'build/', '/banzaicloud/', create: true

    $num_instances = 4

    # centos 8 nodes
    (1..$num_instances).each do |n|
        config.vm.define "centos#{n}" do |node|
            node.vm.box = "centos/8"

            # Monkey patch for https://github.com/dotless-de/vagrant-vbguest/issues/367
            class Foo < VagrantVbguest::Installers::CentOS
                def has_rel_repo?
                    unless instance_variable_defined?(:@has_rel_repo)
                        rel = release_version
                        @has_rel_repo = communicate.test("yum repolist")
                    end
                    @has_rel_repo
                end

                def install_kernel_devel(opts=nil, &block)
                    cmd = "yum update kernel -y"
                    communicate.sudo(cmd, opts, &block)

                    cmd = "yum install -y kernel-devel"
                    communicate.sudo(cmd, opts, &block)

                    cmd = "shutdown -r now"
                    communicate.sudo(cmd, opts, &block)

                    begin
                        sleep 5
                    end until @vm.communicate.ready?
                end
            end
            node.vbguest.installer = Foo

            node.vm.network "private_network", ip: "192.168.64.#{n+10}"
            node.vm.hostname = "centos#{n}"
            node.vm.provider "virtualbox" do |vb|
                vb.name = "centos#{n}"
                vb.memory = "2048"
                vb.cpus = "2"
                vb.customize ["modifyvm", :id, "--audio", "none"]
                vb.customize ["modifyvm", :id, "--memory", "2048"]
                vb.customize ["modifyvm", :id, "--cpus", "2"]
            end

            node.vm.provision "shell" do |s|
                s.inline = <<-SHELL
                dnf install -y yum-utils wget curl chrony vim net-tools socat
                echo 'sync time'
                systemctl enable --now chronyd
                swapoff -a
                modprobe ip_tables
                echo 'ip_tables' >> /etc/modules-load.d/iptables.conf
                echo 'set host name resolution'
                cat >> /etc/hosts <<EOF
192.168.64.11 centos1
192.168.64.12 centos2
192.168.64.13 centos3
192.168.64.14 centos4
EOF
                cat /etc/hosts

                hostnamectl set-hostname centos#{n}

                SHELL
            end
        end
    end

    # Ubuntu LTS nodes
    (1..$num_instances).each do |n|
        config.vm.define "ubuntu#{n}" do |node|
            node.vm.box = "ubuntu/focal64"
            node.vm.network "private_network", ip: "192.168.64.#{n+20}"
            node.vm.hostname = "ubuntu#{n}"

            node.vm.provider "virtualbox" do |vb|
                vb.name = "ubuntu#{n}"
                vb.memory = "2048"
                vb.cpus = "2"
                vb.customize ["modifyvm", :id, "--audio", "none"]
                vb.customize ["modifyvm", :id, "--memory", "2048"]
                vb.customize ["modifyvm", :id, "--cpus", "2"]
            end

            node.vm.provision "shell" do |s|
                s.inline = <<-SHELL

                apt-get update
                apt-get install -y ntp wget curl vim net-tools socat
                echo 'sync time'
                systemctl start ntp
                systemctl enable ntp

                echo 'set host name resolution'
                cat >> /etc/hosts <<EOF
192.168.64.21 ubuntu1
192.168.64.22 ubuntu2
192.168.64.23 ubuntu3
192.168.64.24 ubuntu4
EOF
                cat /etc/hosts

                hostnamectl set-hostname ubuntu#{n}

                SHELL
            end
        end
    end

    config.vm.define "ubuntu-docker-bionic" do |node|
        node.vm.box = "ubuntu/bionic64"
        node.vm.network "private_network", ip: "192.168.64.30"
        node.vm.hostname = "ubuntu-docker"

        node.vm.provider "virtualbox" do |vb|
            vb.name = "ubuntu-docker"
            vb.memory = "2048"
            vb.cpus = "2"
            vb.customize ["modifyvm", :id, "--audio", "none"]
            vb.customize ["modifyvm", :id, "--memory", "2048"]
            vb.customize ["modifyvm", :id, "--cpus", "2"]
        end

        node.vm.provision "shell" do |s|
            s.inline = <<-SHELL

            apt-get update
            apt-get install -y ntp wget curl vim net-tools socat
            echo 'sync time'
            systemctl start ntp
            systemctl enable ntp

            echo 'set host name resolution'
            cat >> /etc/hosts <<EOF
192.168.64.30 ubuntu-docker
EOF
            cat /etc/hosts

            hostnamectl set-hostname ubuntu-docker

            curl -fsSL https://download.docker.com/linux/ubuntu/gpg | apt-key add -
            add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"
            apt-get update && apt-get install -y \
                containerd.io=1.2.13-1 \
                docker-ce=5:19.03.8~3-0~ubuntu-$(lsb_release -cs) \
                docker-ce-cli=5:19.03.8~3-0~ubuntu-$(lsb_release -cs)

            cat > /etc/docker/daemon.json <<EOF
{
  "exec-opts": ["native.cgroupdriver=systemd"],
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "100m"
  },
  "storage-driver": "overlay2"
}
EOF

            mkdir -p /etc/systemd/system/docker.service.d
            systemctl daemon-reload
            systemctl restart docker

            SHELL
        end
    end

    config.vm.define "ubuntu-docker-focal" do |node|
        node.vm.box = "ubuntu/focal64"
        node.vm.network "private_network", ip: "192.168.64.30"
        node.vm.hostname = "ubuntu-docker"

        node.vm.provider "virtualbox" do |vb|
            vb.name = "ubuntu-docker"
            vb.memory = "2048"
            vb.cpus = "2"
            vb.customize ["modifyvm", :id, "--audio", "none"]
            vb.customize ["modifyvm", :id, "--memory", "2048"]
            vb.customize ["modifyvm", :id, "--cpus", "2"]
        end

        node.vm.provision "shell" do |s|
            s.inline = <<-SHELL

            apt-get update
            apt-get install -y ntp wget curl vim net-tools socat
            echo 'sync time'
            systemctl start ntp
            systemctl enable ntp

            echo 'set host name resolution'
            cat >> /etc/hosts <<EOF
192.168.64.30 ubuntu-docker
EOF
            cat /etc/hosts

            hostnamectl set-hostname ubuntu-docker

            curl -fsSL https://download.docker.com/linux/ubuntu/gpg | apt-key add -
            add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"
            apt-get update && apt-get install -y \
                containerd.io=1.4.3-1 \
                docker-ce=5:20.10.1~3-0~ubuntu-$(lsb_release -cs) \
                docker-ce-cli=5:20.10.1~3-0~ubuntu-$(lsb_release -cs)

            cat > /etc/docker/daemon.json <<EOF
{
  "exec-opts": ["native.cgroupdriver=systemd"],
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "100m"
  },
  "storage-driver": "overlay2"
}
EOF

            mkdir -p /etc/systemd/system/docker.service.d
            systemctl daemon-reload
            systemctl restart docker

            SHELL
        end
    end

end
