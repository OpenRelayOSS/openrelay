%define oruser_uid     811
%define oruser_gid     %{oruser_uid}
%define client_name    replay

Name:		openrelay
Version:	0.9
Release:	9
Summary:	OpenRelay - Realtime Relay Messaging Server
Group:		Applications/Internet
License:	GPLv2
URL:		http://openrelay.info/
Source0:	%{name}-%{version}-%{release}.x86_64.tar.gz
Requires:       libunwind 
Requires(pre,postun):  %{_sbindir}/groupadd
Requires(pre,postun):  %{_sbindir}/useradd
Requires(pre,post,preun,postun):  %{_bindir}/systemctl

%description
OpenRelay

%global debug_package %{nil}

%prep
rm -rf %{buildroot}

%setup -n %{name}-%{version}-%{release}.x86_64

%build

%install
install -d %{buildroot}/var/log/%{name} %{buildroot}/usr/local/%{name}/bin/ %{buildroot}/usr/local/%{name}/lib/ %{buildroot}%{_sysconfdir}/sysconfig/ %{buildroot}%{_sysconfdir}/systemd/system/
install %{name} %{buildroot}/usr/local/%{name}/bin/
install %{client_name} %{buildroot}/usr/local/%{name}/bin/
install %{name}-boot.sh %{buildroot}/usr/local/%{name}/bin/
install %{client_name}-boot.sh %{buildroot}/usr/local/%{name}/bin/
install gipcheck %{buildroot}/usr/local/%{name}/bin/
install lipcheck %{buildroot}/usr/local/%{name}/bin/
install libsodium.so.* %{buildroot}/usr/local/%{name}/lib/
install libzmq.so.* %{buildroot}/usr/local/%{name}/lib/
install libczmq.so.* %{buildroot}/usr/local/%{name}/lib/
cd %{buildroot}/usr/local/%{name}/lib
ln -s libsodium.so.*.*.* libsodium.so.23
ln -s libzmq.so.*.*.* libzmq.so.5
ln -s libczmq.so.*.*.* libczmq.so.4
cd -
install -S .rpmsave -b %{name}.env %{buildroot}%{_sysconfdir}/sysconfig/
install %{name}.service %{buildroot}%{_sysconfdir}/systemd/system/

%clean
rm -rf %{buildroot}

%pre
if [ $1 -eq 1 ] ; then
    %{_sbindir}/groupadd -g %{oruser_gid} oruser >/dev/null 2>&1 || :
    %{_sbindir}/useradd -u %{oruser_uid} -s /sbin/nologin -g oruser oruser >/dev/null 2>&1 || :
fi

%post
if [ $1 -eq 1 ] ; then
    %{_bindir}/systemctl daemon-reload >/dev/null 2>&1 || :
fi

%preun
if [ $1 -eq 0 ] ; then
    %{_bindir}/systemctl --no-reload disable %{name}.service > /dev/null 2>&1 || :
    %{_bindir}/systemctl stop %{name}.service > /dev/null 2>&1 || :
fi

%postun
%{_bindir}/systemctl daemon-reload >/dev/null 2>&1 || :
if [ $1 -ge 1 ] ; then
    %{_bindir}/systemctl try-restart %{name}.service >/dev/null 2>&1 || :
elif [ $1 -eq 0 ] ; then
    %{_sbindir}/userdel oruser >/dev/null 2>&1 || :
    %{_sbindir}/groupdel oruser >/dev/null 2>&1 || :
fi

%files
%defattr(0644, root, root, -)
%{_sysconfdir}/systemd/system/%{name}.service
%config(noreplace) %{_sysconfdir}/sysconfig/%{name}.env
%defattr(0755, oruser, oruser, 0755)
/usr/local/%{name}/bin/%{name}
/usr/local/%{name}/bin/%{client_name}
/usr/local/%{name}/bin/%{name}-boot.sh
/usr/local/%{name}/bin/%{client_name}-boot.sh
/usr/local/%{name}/bin/gipcheck
/usr/local/%{name}/bin/lipcheck
/var/log/%{name}
%defattr(-, -, -, -)
/usr/local/%{name}/lib/libsodium.so.23
/usr/local/%{name}/lib/libzmq.so.5
/usr/local/%{name}/lib/libczmq.so.4
%defattr(0775, oruser, oruser, -)
/usr/local/%{name}/lib/libsodium.so.*.*.*
/usr/local/%{name}/lib/libzmq.so.*.*.*
/usr/local/%{name}/lib/libczmq.so.*.*.*

%license LICENSE

%changelog

