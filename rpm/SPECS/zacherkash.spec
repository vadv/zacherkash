%define version unknown
%define debug_package %{nil}

Name:           zacherkash
Version:        %{version}
Release:        1%{?dist}
Summary:        zacherkash
License:        BSD
URL:            http://git.itv.restr.im/infra/zacherkash
Source1:        zacherkash.init
Source2:        zacherkash-logrotate.in
Source3:        zacherkash.yaml
Source:         zacherkash-%{version}.tar.gz
BuildRequires:  make
BuildRequires:  git

%define restream_dir /opt/restream/
%define restream_bin_dir %{restream_dir}/zacherkash/bin

%description
This package provides log parser.

%prep
%setup

%pre
getent group zacherkash > /dev/null || groupadd -r zacherkash
getent passwd zacherkash > /dev/null || \
    useradd -r -g zacherkash -d /var/run/zacherkash -s /sbin/nologin \
    -c "zacherkash user" zacherkash

mkdir -p /var/log/zacherkash
chown zacherkash:zacherkash /var/log/zacherkash


%build
make

%install
mkdir -p %{buildroot}%{restream_bin_dir}
%{__install} -m 0755 -p bin/zacherkash %{buildroot}%{restream_bin_dir}
%{__mkdir} -p %{buildroot}/%{_sysconfdir}/init.d
%{__mkdir} -p %{buildroot}/%{_sysconfdir}/logrotate.d
%{__install} -m 0755 -p %{SOURCE1} %{buildroot}/%{_sysconfdir}/init.d/zacherkash
%{__install} -m 0644 -p %{SOURCE2} %{buildroot}/%{_sysconfdir}/logrotate.d/zacherkash
%{__install} -m 0644 -p %{SOURCE3} %{buildroot}/%{_sysconfdir}/zacherkash.yaml

%clean
rm -rf %{buildroot}

%files
%defattr(-,root,root,-)
%{restream_bin_dir}/zacherkash
%doc README.md
%config(noreplace) %{_sysconfdir}/zacherkash.yaml
%{_sysconfdir}/init.d/zacherkash
%{_sysconfdir}/logrotate.d/zacherkash
