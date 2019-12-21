import urllib.parse
import datetime
import hashlib
import flask
import re

def get_time():
    return str(datetime.datetime.today().strftime("%Y-%m-%d %H:%M:%S"))

def db_data_get(data):
    global set_data

    set_data = data

def db_change(data):
    if set_data == 'mysql':
        data = data.replace('random()', 'rand()')
        data = data.replace('%', '%%')
        data = data.replace('?', '%s')

    return data

def ip_check(d_type = 0):
    if d_type == 0:
        if flask.session and ('state' and 'id') in flask.session:
            ip = flask.session['id']
        else:
            try:
                ip = flask.request.environ.get('HTTP_X_REAL_IP', flask.request.environ.get('HTTP_X_FORWARDED_FOR', flask.request.remote_addr))

                if ip == '::1' or ip == '127.0.0.1':
                    ip = flask.request.environ.get('HTTP_X_FORWARDED_FOR', flask.request.remote_addr)
            except:
                ip = '-'
    else:
        try:
            ip = flask.request.environ.get('HTTP_X_REAL_IP', flask.request.environ.get('HTTP_X_FORWARDED_FOR', flask.request.remote_addr))

            if ip == '::1' or ip == '127.0.0.1':
                ip = flask.request.environ.get('HTTP_X_FORWARDED_FOR', flask.request.remote_addr)
        except:
            ip = '-'

    return str(ip)

def link_fix(main_link):
    if re.search('^:', main_link):
        main_link = re.sub('^:', '', main_link)

    main_link = re.sub('^사용자:', 'user:', main_link)
    main_link = re.sub('^파일:', 'file:', main_link)
    main_link = re.sub('^분류:', 'category:', main_link)

    other_link = re.search('[^\\\\]?(#[^#]+)$', main_link)
    if other_link:
        other_link = other_link.groups()[0]

        main_link = re.sub('(#[^#]+)$', '', main_link)
    else:
        other_link = ''

    main_link = re.sub('\\\\#', '%23', main_link)

    return [main_link, other_link]

def savemark(data):
    data = re.sub("\[date\(now\)\]", get_time(), data)

    ip = ip_check()
    if not re.search("\.", ip):
        name = '[[user:' + ip + '|' + ip + ']]'
    else:
        name = ip

    data = re.sub("\[name\]", name, data)

    return data

def url_pas(data):
    return urllib.parse.quote(data).replace('/','%2F')

def sha224(data):
    return hashlib.sha224(bytes(data, 'utf-8')).hexdigest()

def md5_replace(data):
    return hashlib.md5(data.encode()).hexdigest()
