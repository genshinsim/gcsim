set argument="%2"

set filename=%~1
set output=%filename:txt=json%

"gcsim.exe" -c="%cd%/config/%filename%" -substatOptim=true -out="%cd%/optimized_config/%filename%" %argument% || exit /b %errorlevel%

"gcsim.exe" -c="%cd%/optimized_config/%filename%" -out="%cd%/viewer_gz/%output%" -gz="true" %argument%