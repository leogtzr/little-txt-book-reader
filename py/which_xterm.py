import shutil

if shutil.which('xterm'):
    print('OK xterm')

if shutil.which('vim'):
    print('OK vim')

if not shutil.which('vim23'):
    print('Alv')
