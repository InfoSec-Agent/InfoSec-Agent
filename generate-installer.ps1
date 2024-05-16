# Get version from latest git tag
$VERSION=(git describe --tags $(git rev-list --tags --max-count=1)).Substring(1)

# Build executables
.\build.bat

# Generate installer with Inno Setup
iscc .\generate-installer.iss /DMyAppVersion=$VERSION