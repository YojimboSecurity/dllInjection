# DLL Injection

I created this POC to test Windows Sysmon Create Remote Thread Event ID 8. This 
is based off of the python code found in `Grey Hat Python` By Justin Seitz.

# Testing Windows Sysmon CreateRemoteThread EventID 8

To test Windows Sysmon CreateRemoteThread make sure Sysmon is installed on the machine. Copy the following Sysmon config to a file named `sysmon_config.xml` on the windows machine.

```xml
<Sysmon schemaversion="4.1">
   <EventFiltering>
      <CreateRemoteThread onmatch="include">
          <StartFunction name="technique_id=T1055,technique_name=Process Injection" condition="contains">LoadLibrary</StartFunction>
          <TargetImage name="technique_id=T1055,technique_name=Process Injection" condition="is">C:\Windows\System32\rundll32.exe</TargetImage>
          <TargetImage name="technique_id=T1055,technique_name=Process Injection" condition="is">C:\Windows\System32\svchost.exe</TargetImage>
          <TargetImage name="technique_id=T1055,technique_name=Process Injection" condition="is">C:\Windows\System32\sysmon.exe</TargetImage>
      </CreateRemoteThread>
  </EventFiltering>
</Sysmon>
```

**Note:** This Sysmon config is from Olaf Hartong's sysmon-modular and can be found here https://github.com/olafhartong/sysmon-modular/blob/version-8/8_create_remote_thread/include_dll_injection.xml

Load this config file with the following command.

```powershell
Sysmon.exe -c sysmon_config.xml
```

List the running processes for `svchost` by executing the following command in PowerShell

```powerShell
Get-Process svchost
```
You should see something similar to the output below

```powershell

Handles  NPM(K)    PM(K)      WS(K) VM(M)   CPU(s)     Id ProcessName
-------  ------    -----      ----- -----   ------     -- -----------
    172      12     2064       1700 ...22     0.03    180 svchost
    565      32    42896      40688 ...10    56.16    312 svchost
    568      20     4872       8104 ...12     3.50    644 svchost
    514      15     3092       4972 ...84     2.22    704 svchost
   1654      54    14784      23488 ...37     8.23    884 svchost
   1305      65    68356      61512 ...64    32.67    892 svchost
    867      33    26652      13540 ...62     7.25    956 svchost
    494      28     5892       9724 ...43     2.13   1044 svchost
    502      41    14152      12932 ...58     5.34   1372 svchost
    427      21     5880      11980 ...86     2.36   1592 svchost
    180      15     3484       6540 ...37     2.06   1696 svchost
    117       8     1384       6344 ...84     0.02   1984 svchost
    114       9     1236       2916 ...77     0.66   2116 svchost
```

Select one of the PIDs, in this example, I will use `2116`.

Next, you will need a DLL to inject onto the process. For this, I have just used one that is available in `System32`, however, any will do. 

Before running `dllInjection.exe` make sure to have event viewer open and on Sysmon operational. I also recommend filtering for event Id 8.

Run `dllInjection.exe` as seen in the example below.

```powershell
dllInjection.exe -p 2116 -d C:\Windows\System32\kernel32.dll
```
After running `dllInjection.exe` you should see an event that looks like the following.

```xml
<Events>
  <Event xmlns="http://schemas.microsoft.com/win/2004/08/events/event">
    <System>
      <Provider Guid="{5770385F-C22A-43E0-BF4C-06F5698FFBD9}" Name="Microsoft-Windows-Sysmon"></Provider>
      <EventID>8</EventID>
      <Version>2</Version>
      <Level>4</Level>
      <Task>8</Task>
      <Opcode>0</Opcode>
      <Keywords>0x8000000000000000</Keywords>
      <TimeCreated SystemTime="2020-07-20T19:00:02.974935200Z"></TimeCreated>
      <EventRecordID>29220</EventRecordID>
      <Correlation></Correlation>
      <Execution ProcessID="1812" ThreadID="2300"></Execution>
      <Channel>Microsoft-Windows-Sysmon/Operational</Channel>
      <Computer>DevBox12.win2012r2.local</Computer>
      <Security UserID="S-1-5-18"></Security>
    </System>
    <EventData>
      <Data Name="RuleName">technique_id=T1055,technique_name=Process Injection</Data>
      <Data Name="UtcTime">2020-07-20 19:00:02.974</Data>
      <Data Name="SourceProcessGuid">{32162397-E9B2-5F15-D517-E60200000000}</Data>
      <Data Name="SourceProcessId">9468</Data>
      <Data Name="SourceImage">\\vboxsvr\go\dllInjection\git.yojimbosecurity.com\dllInjection\dllInjection.exe</Data>
      <Data Name="TargetProcessGuid">{32162397-3369-5F13-63CA-030000000000}</Data>
      <Data Name="TargetProcessId">2688</Data>
      <Data Name="TargetImage">C:\Windows\System32\svchost.exe</Data>
      <Data Name="NewThreadId">6652</Data>
      <Data Name="StartAddress">0x00007FFACF0F4960</Data>
      <Data Name="StartModule">C:\Windows\system32\KERNEL32.DLL</Data>
      <Data Name="StartFunction">LoadLibraryA</Data>
  </EventData>
</Event>
```