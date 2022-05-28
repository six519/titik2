Name: titik          
Version: 3       
Release: 3
Summary: Titik Interpreter        

License: MIT      
URL: https://github.com/six519/titik2     
Source0: titik.tar.gz       

BuildArch: x86_64
Requires: golang

%description
Titik Programming Language/Interpreter Linux Release. This is for educational/experimental purpose only.

%prep
%setup -n titik

%install
go build .
rm -rf $RPM_BUILD_ROOT
mkdir -p $RPM_BUILD_ROOT/bin/
cp -r titik2 $RPM_BUILD_ROOT/bin/titik


%files
/bin/titik

