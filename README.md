# frog
镜像同步程序

![](./frog.jpg)

[![Build Status](https://travis-ci.org/ckeyer/frog.png?branch=master)](https://travis-ci.org/ckeyer/frog)
[![Go Report Card](https://goreportcard.com/badge/github.com/ckeyer/frog)](https://goreportcard.com/report/github.com/ckeyer/frog)
[![GoDoc](https://godoc.org/github.com/ckeyer/frog?status.png)](http://godoc.org/github.com/ckeyer/frog)
[![license](https://img.shields.io/badge/license-GPL%20V3.0-blue.svg?maxAge=2592000)](https://github.com/ckeyer/frog/blob/master/LICENSE)


# 配置文件

config.yaml
```
global:
  period: 10h
  deleteeverytime: true
registries:
  - username: <your username>
    password: <your password>
    server: registry.cn-beijing.aliyuncs.com
tasks:
  - origin: ckeyer/dev
    target: registry.cn-beijing.aliyuncs.com/wa/dev
    tags: ["vue", "k8s19", "tf-py3", "ng2", "rpm", "py3", "php7", "node", "nginx"]
  - origin: alpine
    target: registry.cn-beijing.aliyuncs.com/wa/alpine
    tags: ["edge", "3.5", "3.6", "3.7"]
  - origin: nginx
    target: registry.cn-beijing.aliyuncs.com/wa/nginx
    tags: ["1.13", "1.12"]
  - origin: redis
    target: registry.cn-beijing.aliyuncs.com/wa/redis
    tags: ["3.2", "4.0"]
  - origin: ubuntu
    target: registry.cn-beijing.aliyuncs.com/wa/ubuntu
    tags: ["17.10", "18.04", "16.04"]
  - origin: mongo
    target: registry.cn-beijing.aliyuncs.com/wa/mongo
    tags: ["3.0", "3.2", "3.4", "3.6", "3.7"]
  - origin: mysql
    target: registry.cn-beijing.aliyuncs.com/wa/mysql
    tags: ["5.5", "5.6", "5.7", "8.0"]
  - origin: centos
    target: registry.cn-beijing.aliyuncs.com/wa/centos
    tags: ["6", "7"]
  - origin: php
    target: registry.cn-beijing.aliyuncs.com/wa/php
    tags: ["5.6", "7.2", "7"]
  - origin: quay.io/coreos/etcd
    target: registry.cn-beijing.aliyuncs.com/wa/etcd
    tags: ["3.2", "3.1", "v3.3"]
  - origin: gcr.io/google-containers/kube-apiserver-amd64
    target: registry.cn-beijing.aliyuncs.com/wa/kube-apiserver-amd64
    tags: ["v1.9.6", "v1.9.5", "v1.9.4", "v1.10.0"]
  - origin: gcr.io/google-containers/kube-controller-manager-amd64
    target: registry.cn-beijing.aliyuncs.com/wa/kube-controller-manager-amd64
    tags: ["v1.9.6", "v1.9.5", "v1.9.4", "v1.10.0"]
  - origin: gcr.io/google-containers/kube-scheduler-amd64
    target: registry.cn-beijing.aliyuncs.com/wa/kube-scheduler-amd64
    tags: ["v1.9.6", "v1.9.5", "v1.9.4", "v1.10.0"]
  - origin: gcr.io/google-containers/kube-proxy-amd64
    target: registry.cn-beijing.aliyuncs.com/wa/kube-proxy-amd64
    tags: ["v1.9.6", "v1.9.5", "v1.9.4", "v1.10.0"]
  - origin: gcr.io/google-containers/kube-nethealth-amd64
    target: registry.cn-beijing.aliyuncs.com/wa/kube-nethealth-amd64
    tags: ["1.0"]
  - origin: gcr.io/google-containers/kube-dnsmasq-amd64
    target: registry.cn-beijing.aliyuncs.com/wa/kube-dnsmasq-amd64
    tags: ["1.4.1", "1.4", "1.3", "1.2"]
  - origin: gcr.io/google-containers/pause-amd64
    target: registry.cn-beijing.aliyuncs.com/wa/pause-amd64
    tags: ["3.0", "3.1"]
  - origin: calico/cni
    target: registry.cn-beijing.aliyuncs.com/wa/calico-cni
    tags: ["v2.0.3", "v2.0.2", "v2.0.1", "v1.11.3", "v1.11.2", "v1.11.1"]
  - origin: calico/kube-controllers
    target: registry.cn-beijing.aliyuncs.com/wa/calico-kube-controllers
    tags: ["v2.0.0", "v2.0.2", "v2.0.1", "v1.0.3", "v1.0.2", "v1.0.1"]
  - origin: calico/node
    target: registry.cn-beijing.aliyuncs.com/wa/calico-node
    tags: ["v3.0.4", "v3.0.3", "release-v3.1", "release-v3.0", "v2.6.8", "v1.0.2", "v1.0.1"]
  - origin: calico/ctl
    target: registry.cn-beijing.aliyuncs.com/wa/calico-ctl
    tags: ["v2.0.2", "v2.0.1", "v2.0.0", "release-v3.0", "release-v3.1", "release-v2.6", "release-v2.5"]
  - origin: istio/envoy
    target: registry.cn-beijing.aliyuncs.com/wa/istio-envoy
    tags: ["latest"]
  - origin: istio/pilot
    target: registry.cn-beijing.aliyuncs.com/wa/istio-pilot
    tags: ["0.7.1", "0.7.0", "0.6.0", "0.5.1"]
  - origin: istio/mixer
    target: registry.cn-beijing.aliyuncs.com/wa/istio-mixer
    tags: ["0.7.1", "0.7.0", "0.6.0", "0.5.1"]
  - origin: istio/istio-ca
    target: registry.cn-beijing.aliyuncs.com/wa/istio-ca
    tags: ["0.7.1", "0.7.0", "0.6.0", "0.5.1"]
  - origin: istio/proxy
    target: registry.cn-beijing.aliyuncs.com/wa/istio-proxy
    tags: ["0.7.1", "0.7.0", "0.6.0", "0.5.1"]
  - origin: istio/proxy_init
    target: registry.cn-beijing.aliyuncs.com/wa/istio-proxy_init
    tags: ["0.7.1", "0.7.0", "0.6.0", "0.5.1"]
```
