
from glob import iglob
import re

# TODO: Add a cleanup run for ./js/app/*.js
# TODO: Test

def between(string, c1, c2):
  return(string.split(c1)[1].split(c2)[2])


def compile_apps(i_app_dir, o_app_dir):
  # TODO: Ensure output dir exists
  added_CN = False
  imported = []
  for appFile in iglob('*.js', root_dir=input_app_dir):
    with open(i_app_dir + appFile, r) as i_f, open(o_app_dir + appFile, w) as o_f:
      for line in i_f:
        if line.startswith('import') and line.endswith('\'../main.js\';'):
          # XXX: Requires semicolons at end of line
          if not added_CN:
            o_f.write('const CN = globalThis.CN')
            added_CN = True
          imported.append(between(line, '{', '}'))
          # TODO: Strip whitespace
          # TODO: Fix for multi imports
        else:
          for toReplace in imported:
            re.sub('(?<=[ ^])' + toReplace, 'CN.' + toReplace, line)
          o_f.write(line)


def main():
  # TODO Add lint of js files
  # TODO: Modify exporting file
  compile_apps(i_app_dir = './js/app/', o_app_dir = './compiled/app/')


if __name__ == '__main__':
  main()
