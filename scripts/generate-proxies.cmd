: # This is a special script which intermixes both sh
: # and cmd code. It is written this way because it is
: # used in system() shell-outs directly in otherwise
: # portable code. See https://stackoverflow.com/questions/17510688
: # for details.
:; echo Running in bash mode ; cd $( cd "$(dirname "${BASH_SOURCE[0]}")" ; pwd -P ); ./generate-proxies.sh ; exit
@ECHO OFF
ECHO Running in cmd mode
ECHO %%~dp0 is "%~dp0"
CMD /C "%~dp0generate-proxies.bat"