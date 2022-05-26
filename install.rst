Installation Instructions
=========================


Linux (Ubuntu 22.04, Debian 13, Linux Mint 20.3, Raspbian 11, Elementary OS 5.1)
--------------------------------------------------------------------------------

::

    curl -s https://packagecloud.io/install/repositories/six519/titik/script.deb.sh | sudo bash && sudo apt install titik

Linux (Red Hat 7/8, CentOS 7/8, Fedora 36, openSUSE 42.3, Scientific Linux 7, Oracle Linux 8)
---------------------------------------------------------------------------------------------

::

    curl -s https://packagecloud.io/install/repositories/six519/titik/script.rpm.sh | sudo bash && sudo yum install titik -y

Mac OS X
--------

Need to install Homebrew_ first

.. _Homebrew: https://brew.sh/

::

    brew tap six519/tap
    brew install titik

FreeBSD 13.1
------------

::
    
    sudo pkg install go
    sudo pkg install git
    sudo pkg add "https://github.com/six519/titik2/blob/master/bin/titik-3.0.0.pkg?raw=true"

Windows 10 (64 bit only)
------------------------

Download the installer from the link below then restart the system after installation

::

    https://github.com/six519/titik2/blob/master/bin/titik-3.3.0-setup.exe?raw=true
