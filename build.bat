@REM
@REM  Licensed to the Apache Software Foundation (ASF) under one or more
@REM  contributor license agreements.  See the NOTICE file distributed with
@REM  this work for additional information regarding copyright ownership.
@REM  The ASF licenses this file to You under the Apache License, Version 2.0
@REM  (the "License"); you may not use this file except in compliance with
@REM  the License.  You may obtain a copy of the License at
@REM
@REM      http://www.apache.org/licenses/LICENSE-2.0
@REM
@REM  Unless required by applicable law or agreed to in writing, software
@REM  distributed under the License is distributed on an "AS IS" BASIS,
@REM  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
@REM  See the License for the specific language governing permissions and
@REM  limitations under the License.

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
