Installation Instructions
=========================


Linux (Ubuntu, Debian)
----------------------

::

    curl -s https://packagecloud.io/install/repositories/six519/titik/script.deb.sh | sudo bash && sudo apt install titik

Linux (Red Hat, CentOS, Fedora)
-------------------------------

::

    curl -s https://packagecloud.io/install/repositories/six519/titik/script.rpm.sh | sudo bash && sudo yum install titik -y

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
    sudo pkg add "https://github.com/six519/titik2/blob/master/bin/titik-3.0.0.pkg?raw=true"

Windows (64 bit only)
---------------------

Download the installer from the link below then restart the system after installation

::

    https://github.com/six519/titik2/blob/master/bin/titik-3.0.0-setup.exe?raw=true