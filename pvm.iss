; Script generated by the Inno Setup Script Wizard.
; SEE THE DOCUMENTATION FOR DETAILS ON CREATING INNO SETUP SCRIPT FILES!

#define MyAppName "pvm"
#define MyAppVersion "1.0.0"
#define MyAppPublisher "UnreadCode, Inc."
#define MyAppURL "https://github.com/unreadcode/pvm"
#define MyAppExeName "pvm.exe"
#define MyAppAssocName MyAppName + " File"
#define MyIcon "bin\php.ico"
#define MyAppId "7582EF6C-90E2-4A3E-98E3-81706DF20EE3"
#define PVM_DIR "."


[Setup]
; NOTE: The value of AppId uniquely identifies this application. Do not use the same AppId value in installers for other applications.
; (To generate a new GUID, click Tools | Generate GUID inside the IDE.)
AppId={#MyAppId}
AppName={#MyAppName}
AppVersion={#MyAppVersion}
AppVerName={#MyAppName} {#MyAppVersion}
AppPublisher={#MyAppPublisher}
AppPublisherURL={#MyAppURL}
AppSupportURL={#MyAppURL}
AppUpdatesURL={#MyAppURL}
DefaultDirName={commonpf}\{#MyAppName}
DisableDirPage=no
ChangesAssociations=yes
DefaultGroupName={#MyAppName}
LicenseFile={#PVM_DIR}\LICENSE
; Uncomment the following line to run in non administrative install mode (install for current user only.)
PrivilegesRequired=admin
OutputDir={#PVM_DIR}\dist\{#MyAppVersion}
OutputBaseFilename={#MyAppName}-Setup
SetupIconFile={#PVM_DIR}\{#MyIcon}
Compression=lzma
SolidCompression=yes
ChangesEnvironment=yes
DisableProgramGroupPage=yes
WizardStyle=modern
UninstallDisplayIcon={app}\{#MyIcon}
VersionInfoVersion={#MyAppVersion}
VersionInfoCopyright=Copyright (C) 2024 UnreadCode Inc.
VersionInfoProductName={#MyAppName}

[Languages]
Name: "english"; MessagesFile: "compiler:Default.isl"
Name: "zh_cn"; MessagesFile: "compiler:Languages\ChineseSimplified.isl"

[Files]
Source: "{#PVM_DIR}\bin\*"; DestDir: "{app}"; Flags: ignoreversion recursesubdirs createallsubdirs;

; NOTE: Don't use "Flags: ignoreversion" on any shared system files

[Icons]
Name: "{group}\{#MyAppName}"; Filename: "{app}\{#MyAppExeName}"; IconFilename: "{#MyIcon}"

[Code]

function InitializeUninstall(): Boolean;
var spath, upath, phpPath: string;
begin
    // remove php_path
    RegQueryStringValue(HKEY_LOCAL_MACHINE, 'SYSTEM\CurrentControlSet\Control\Session Manager\Environment', 'PHP_PATH', phpPath);
    RemoveDir(phpPath);
    // Clean the registry
    RegDeleteValue(HKEY_LOCAL_MACHINE, 'SYSTEM\CurrentControlSet\Control\Session Manager\Environment', 'PVM_ROOT')
    RegDeleteValue(HKEY_LOCAL_MACHINE, 'SYSTEM\CurrentControlSet\Control\Session Manager\Environment', 'PHP_PATH')
    RegDeleteValue(HKEY_CURRENT_USER, 'Environment', 'PVM_ROOT')
    RegDeleteValue(HKEY_CURRENT_USER, 'Environment', 'PHP_PATH')
    // system environment variable
    RegQueryStringValue(HKEY_LOCAL_MACHINE, 'SYSTEM\CurrentControlSet\Control\Session Manager\Environment', 'Path', spath);
    StringChangeEx(spath,'%PVM_ROOT%','',True);
    StringChangeEx(spath,'%PHP_PATH%','',True);
    StringChangeEx(spath,';;',';',True);
    RegWriteExpandStringValue(HKEY_LOCAL_MACHINE, 'SYSTEM\CurrentControlSet\Control\Session Manager\Environment', 'Path', spath);

    // user environment variable
    RegQueryStringValue(HKEY_CURRENT_USER, 'Environment', 'Path', upath);
    StringChangeEx(upath,'%PVM_ROOT%','',True);
    StringChangeEx(upath,'%PHP_PATH%','',True);
    StringChangeEx(upath,';;',';',True);
    RegWriteExpandStringValue(HKEY_CURRENT_USER, 'Environment', 'Path', upath);
    Result := True;
end;



var PhpPathPage: TInputDirWizardPage;

procedure InitializeWizard;
begin
  PhpPathPage := CreateInputDirPage(wpSelectDir,
  'Select PHP Path', 'Where is PHP installed?',
  'Please select PHP installation directory.', False, '');
  PhpPathPage.Add('This directory will automatically be added to your system path.');
  PhpPathPage.Values[0] := ExpandConstant('{commonpf}\php');
end;

procedure CurStepChanged(CurStep: TSetupStep);
var path: string;
begin
  if CurStep = ssPostInstall then
  begin
      RegWriteExpandStringValue(HKEY_LOCAL_MACHINE, 'SYSTEM\CurrentControlSet\Control\Session Manager\Environment', 'PVM_ROOT', ExpandConstant('{app}'));
      RegWriteExpandStringValue(HKEY_LOCAL_MACHINE, 'SYSTEM\CurrentControlSet\Control\Session Manager\Environment', 'PHP_PATH', PhpPathPage.Values[0]);
      RegWriteExpandStringValue(HKEY_CURRENT_USER, 'Environment', 'PVM_ROOT', ExpandConstant('{app}'));
      RegWriteExpandStringValue(HKEY_CURRENT_USER, 'Environment', 'PHP_PATH', PhpPathPage.Values[0]);

      // system environment variable
      RegQueryStringValue(HKEY_LOCAL_MACHINE, 'SYSTEM\CurrentControlSet\Control\Session Manager\Environment', 'Path', path);
      if Pos('%PVM_ROOT%', path) = 0 then begin
        path := path + ';%PVM_ROOT%';
        StringChangeEx(path,';;',';',True);
        RegWriteExpandStringValue(HKEY_LOCAL_MACHINE, 'SYSTEM\CurrentControlSet\Control\Session Manager\Environment', 'Path', path);
      end;
      if Pos('%PHP_PATH%', path) = 0 then begin
        path := path + ';%PHP_PATH%';
        StringChangeEx(path,';;',';',True);
        RegWriteExpandStringValue(HKEY_LOCAL_MACHINE, 'SYSTEM\CurrentControlSet\Control\Session Manager\Environment', 'Path', path);
      end;
      // user environment variable
      RegQueryStringValue(HKEY_CURRENT_USER, 'Environment', 'Path', path);
      if Pos('%PVM_ROOT%', path) = 0 then begin
        path := path + ';%PVM_ROOT%';
        StringChangeEx(path,';;',';',True);
        RegWriteExpandStringValue(HKEY_CURRENT_USER, 'Environment', 'Path', path);
      end;

      if Pos('%PHP_PATH%', path) = 0 then begin
        path := path + ';%PHP_PATH%';
        StringChangeEx(path,';;',';',True);
        RegWriteExpandStringValue(HKEY_CURRENT_USER, 'Environment', 'Path', path);
      end;
  end;
end;

[UninstallDelete]
Type: files; Name: "{app}\*"
Type: filesandordirs; Name: "{app}\*"