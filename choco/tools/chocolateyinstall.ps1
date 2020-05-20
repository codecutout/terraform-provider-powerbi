$ErrorActionPreference = 'Stop'; # stop on all errors

$toolsDir = "$(Split-Path -parent $MyInvocation.MyCommand.Definition)"
$packageArgs = @{
    packageName   = $env:ChocolateyPackageName
    unzipLocation = Join-Path $env:APPDATA "terraform.d\plugins"
    file          = Join-Path $toolsDir 'terraform-plugins.zip'
  }

Install-ChocolateyZipPackage @packageArgs
