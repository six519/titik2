Name: titik          
Version: 3       
Release: 2
Summary: Titik Interpreter        

License: MIT      
URL: https://github.com/six519/titik2     
Source0: titik.tar.gz       

BuildArch: x86_64

%description
Titik Programming Language/Interpreter Linux Release. This is for educational/experimental purpose only.

%prep
%setup -n titik

%install
rm -rf $RPM_BUILD_ROOT
mkdir -p $RPM_BUILD_ROOT/bin/
cp -r * $RPM_BUILD_ROOT/bin/


%files
/bin/titik

