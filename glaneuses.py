# -*- coding: utf-8 -*-
import requests
import xmlrpclib


username = 'mkouhei'
apis = {'udd': ('http://udd.debian.org/dmd/'
                '?email1=mkouhei%40palmtb.net&format=json'),
        'pypi': 'http://pypi.python.org/pypi',
        'github': 'https://api.github.com/users/%s/events' % username}

udd = requests.get(apis.get('udd')).json()

client = xmlrpclib.ServerProxy(apis.get('pypi'))
pypi = client.user_packages(username)
pypi_list = [(pkg[1], client.package_releases(pkg[1])[0]) for pkg in pypi]

github_activity = requests.get(apis.get('github')).json()

print(udd)
print(pypi_list)
print(github_activity)
