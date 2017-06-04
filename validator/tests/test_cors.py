import requests


def test_preflight(testserver):
    r = requests.options(testserver)
    assert r.status_code == 200
    assert r.headers.get('Access-Control-Allow-Origin') == '*'


def test_cors_on_400(testserver):
    r = requests.post(testserver + 'v1/validate/')
    assert r.status_code == 400, (r.status_code, r.content)
    assert r.headers.get('Access-Control-Allow-Origin') == '*'


def test_cors_on_404(testserver):
    r = requests.get(testserver + 'asdfasdfasdf')
    assert r.status_code == 404, (r.status_code, r.content)
    assert r.headers.get('Access-Control-Allow-Origin') == '*'
