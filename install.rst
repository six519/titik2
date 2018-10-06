Installation Instructions
=========================


Linux (Ubuntu, Debian)
----------------------

::

    wget -qO - 'https://bintray.com/user/downloadSubjectPublicKey?username=six519' | sudo apt-key add -
    echo "deb https://dl.bintray.com/six519/debian all main" | sudo tee -a /etc/apt/sources.list
    sudo apt update
    sudo apt install titik

Linux (Red Hat, CentOS, Fedora)
-------------------------------

::

    wget https://bintray.com/six519/rpm/rpm -O bintray-six519-rpm.repo
    sudo mv bintray-six519-rpm.repo /etc/yum.repos.d/
    sudo yum update
    sudo yum install titik

Mac OS X
--------

Need to install Homebrew_ first

.. _Homebrew: https://brew.sh/

::

    brew tap six519/tap
    brew install titik

FreeBSD
-------

::
    
    sudo pkg install go
    sudo pkg install git
    sudo pkg add https://dl.bintray.com/six519/Generic/titik-2.0.1.txz

Windows (64 bit only)
---------------------

Need to install Chocolatey_ first

.. _Chocolatey: https://chocolatey.org/

::

    choco install titik