%if 0%{?rhel} == 7
  %define dist .el7
%endif
%define _unpackaged_files_terminate_build 0
Name:	gomound
Version: 0.1
Release: 3%{?dist}
Summary: A golang daemon to display zookeeper information.

License: GPLv2
URL: https://github.com/Jmainguy/gophermound
Source0: gomound.tar.gz
Requires(pre): shadow-utils


%description
A golang daemon to display zookeeper information over http.

%prep
%setup -q -n gophermound
%install
mkdir -p $RPM_BUILD_ROOT/usr/sbin
mkdir -p $RPM_BUILD_ROOT/opt/gomound
mkdir -p $RPM_BUILD_ROOT/usr/lib/systemd/system
mkdir -p $RPM_BUILD_ROOT/etc/rc.d/init.d/
install -m 0755 $RPM_BUILD_DIR/gophermound/gomound %{buildroot}/usr/sbin
install -m 0644 $RPM_BUILD_DIR/gophermound/service/gomound.service %{buildroot}/usr/lib/systemd/system
install -m 0755 $RPM_BUILD_DIR/gophermound/service/gomound.sysv %{buildroot}/etc/rc.d/init.d/gomound

%files
/usr/sbin/gomound
%if 0%{?rhel} == 6
  /etc/rc.d/init.d/gomound
%else
  /usr/lib/systemd/system/gomound.service
%endif
%dir /opt/gomound
%doc

%pre
getent group gomound >/dev/null || groupadd -r gomound
getent passwd gomound >/dev/null || \
    useradd -r -g gomound -d /opt/gomound -s /sbin/nologin \
    -c "User to run gomound service" gomound
exit 0
%post
chown -R gomound:gomound /opt/gomound
if [ -f /usr/bin/systemctl ]; then
  systemctl daemon-reload
fi

%changelog
