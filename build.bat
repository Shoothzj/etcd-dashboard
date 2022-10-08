SET DIR=%~dp0
echo %DIR%
call "%DIR%\portal\build.bat"
mvn -B clean package -Dmaven.test.skip=true
if not exist "%DIR%\dist" mkdir "%DIR%\dist"
del "%DIR%\dist\*" /S /Q
xcopy "%DIR%\portal\build" "%DIR%\dist\static" /E /Q /D
copy "%DIR%\target\paas-dashboard-java-0.0.1-SNAPSHOT.jar" "%DIR%\dist\paas-dashboard.jar"
xcopy "%DIR%\target\lib" "%DIR%\dist\lib" /E /Q /D
xcopy "%DIR%\target\conf" "%DIR%\dist\conf" /E /Q /D
xcopy "%DIR%\docker-build\*" "%DIR%\dist\" /E /Q /D
cd %DIR%
