set argument="%2"

set filename=%~1
set output=%filename:txt=json%

"gcsim.exe" -c="%cd%/config/%filename%" -out="%cd%/viewer_gz/%output%" -gz="true" %argument%
