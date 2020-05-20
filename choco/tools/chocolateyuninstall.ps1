$ErrorActionPreference = 'Stop'; # stop on all errors

$packageArgs = @{
    packageName   = $env:ChocolateyPackageName
    zipFileName   = "$($env:ChocolateyPackageName)Install.zip"
}

Uninstall-ChocolateyZipPackage @packageArgs
