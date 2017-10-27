@echo off

rem TODO generate random suffix
echo image_suffix=xxxxx

echo --- Building docker image
docker build -f .\Dockerfile-windows --tag buildkiteagent%image_suffix% . || goto :error

echo +++ Running tests
docker run --rm buildkiteagent%image_suffix% go test -v ./... || goto :error
goto :EOF

:error
set previous_errorlevel=%errorlevel%
docker rmi buildkiteagent%image_suffix%
echo Failed with error #%previous_errorlevel%.
exit /b %previous_errorlevel%