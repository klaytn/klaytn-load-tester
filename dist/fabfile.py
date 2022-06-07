from __future__ import with_statement
from fabric.api import run, roles, execute, env, task, hosts, cd, local, settings, abort, parallel, runs_once, serial
from fabric.contrib.console import confirm
from fabric.operations import put
from fabric.context_managers import lcd
import boto3
import json
from subprocess import call,check_call
import os,errno, pexpect

env.hosts = []
env.user = 'ubuntu'
env.key_filename = '/Users/tom/aws/klay-load.pem'
env.colorize_errors = True
env['abort_on_prompts'] = False

# 'klay-load-base' on Seoul region
imageId = 'ami-00ca7ffe117e2fe91'
instanceType = 'c4.2xlarge'

## master private ip
masterpublicip = "13.124.198.168"
masterprivateip = "172.31.31.6"

## klaynode rpc ip and port list
endpoints = ['http://172.31.31.6:8545', 'http://172.31.31.6:8546', 'http://172.31.31.6:8547', 'http://172.31.31.6:8548']
## test case names
# "transferSignedTx,transferUnsignedTx,cpuHeavyTx"
testcase = "transferSignedTx"

## private key list of pre-fund account
coinbasekeylist = ['349343aad78f398528907e62b62ce7e7e3c9f57c674e12bbd03857682353a73f',
                   '671a7c20cf6edb66f85e1202aaf064c60cfff91a362f59cc8704b73353453783',
                   '297068ea65e1d04a5533548427a3cb543eb250f4aee1f7a49951a73fd1eca41f',
                   '55fa85f7773e4d0a4e8261c289d555eb75abf51eb127bf3f6a13f06b598ce6b3']

def	get_privates(tagname):

    regions = [ 'ap-northeast-2' ]

    hosts = []
    for i in range (0, len(regions)):
        ec2 = boto3.client('ec2', region_name = regions[i])
        response = ec2.describe_instances(
            Filters = [
                {
                    'Name': 'tag-value',
                    'Values': [
                        tagname,
                    ]
                },
            ],
        )

        for reservation in response["Reservations"]:
            for instance in reservation["Instances"]:
                for interface in instance["NetworkInterfaces"]:
                    hosts.append(interface["PrivateIpAddress"])

    print(hosts)
    #env.hosts = hosts
    return hosts


def get_publics(tagname):
    regions = ['ap-northeast-2']

    hosts = []
    for i in range(0, len(regions)):
        ec2 = boto3.client('ec2', region_name=regions[i])
        response = ec2.describe_instances(
            Filters=[
                {
                    'Name': 'tag-value',
                    'Values': [
                        tagname,
                    ]
                },
            ],
        )

        for reservation in response["Reservations"]:
            for instance in reservation["Instances"]:
                for interface in instance["NetworkInterfaces"]:
                    hosts.append(instance["PublicIpAddress"])

    print(hosts)
    # env.hosts = hosts
    return hosts


@task
@serial
@runs_once
def get_master(tagname):
    hosts = get_publics(tagname)
    print("all hosts = " + str(hosts))
    env.hosts = [hosts[0]]
    print("hosts = " + str(env.hosts))
    return hosts[0]


@task
@serial
@runs_once
def get_publicip(tagname):
    hosts = get_publics(tagname)
    print("all hosts = " + str(hosts))
    env.hosts = hosts
    print("hosts = " + str(env.hosts))
    return hosts

@task
@serial
@runs_once
def get_privateip(tagname):
    hosts = get_privates(tagname)
    print("all hosts = " + str(hosts))
    env.hosts = hosts
    print("private host = " + str(env.hosts))
    return hosts


@task
def master():
    with settings(warn_only=True):
        with cd("~/locust"):
            print(env.hosts)
            run("pwd")
            run("hostname -f")
            run("sh master.sh")

@task
@serial
@runs_once
def startmaster():
    master()

@task
def startslave(rps):
    slave(rps)

@task
def slave(rps):
    with settings(warn_only=True):
        with cd("~/locust"):
            print(env.hosts)
            print(env.host)
            run("pwd")
            run("hostname -f")
            idx = env.hosts.index(env.host)
            run("sh slave.sh " + str(rps) + " " + str(masterprivateip) + " " + coinbasekeylist[idx % len(coinbasekeylist)] + " " + str(testcase) + " " + str(endpoints[idx % len(endpoints)]))


@task
def stopmaster():
    with settings(warn_only=True):
        with cd("~/locust"):
            print(env.hosts)
            print(env.host)
            run("pwd")
            run("hostname -f")
            run("pkill locust")

@task
def stopslave():
    with settings(warn_only=True):
        with cd("~/locust"):
            print(env.hosts)
            print(env.host)
            run("pwd")
            run("hostname -f")
            run("killall -2 klayslave")
            run("sleep 2")
            run("pkill klayslave")

@task
def buildslave():
    with settings(warn_only=True):
        with cd("~/go/src/github.com/ground-x/locust-load-tester"):
            print(env.host)
            run("pwd")
            run("hostname -f")
            run("git pull")
            with cd("klayslave"):
                run("/usr/local/go/bin/go build")
                run("mv ~/go/src/github.com/ground-x/locust-load-tester/klayslave/klayslave ~/locust")
