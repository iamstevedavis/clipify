Write-Host "Building Clipify.exe..."

# Build the executable
go build -ldflags "-H windowsgui" -o bin/Clipify.exe ./cmd
if ($?) {
    Write-Host "Build successful!" -ForegroundColor Green

    # Copy the assets folder to the bin directory
    $sourceAssets = "assets"
    $destinationAssets = "bin/assets"
    if (Test-Path $destinationAssets) {
        Remove-Item -Recurse -Force $destinationAssets
    }
    Copy-Item -Recurse -Force $sourceAssets $destinationAssets

    Write-Host "Assets copied to bin directory." -ForegroundColor Green
} else {
    Write-Host "Build failed!" -ForegroundColor Red
}