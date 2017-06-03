from multiprocessing import Process
import time

import pytest

from service import app


@pytest.fixture(scope='module')
def testserver():
    port = 6868

    def _serve():
        app.run(host='127.0.0.1', port=port, quiet=True)

    p = Process(target=_serve)
    print('Starting test server on port %d...' % port)
    p.start()
    time.sleep(1)
    yield 'http://127.0.0.1:%d/' % port
    p.terminate()
