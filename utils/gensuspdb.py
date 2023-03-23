import sqlite3

# it would be nice to scrape this from somewhere online instead of having it hardcoded,
# for the sake of accuracy and completeness...
# Common system processes, to attempt detection of malware masquerading as system processes
# Example: lsass.exe in c:\Users\John\Documents and running in userland is probably not legit
systemProcs = [
    ("lsass.exe", "C:\\Windows\\System32", "Local Security Authority Subsystem Service"),
    ("csrss.exe", "C:\\Windows\\System32", "Client Server Runtime Subsystem"),
    ("explorer.exe", "C:\\Windows\\explorer.exe", "User session"),
    ("svchost.exe", "C:\\Windows\\System32", "Service Host")
]

# Suspicious functions from libraries
# Networking, RunPE stuff, ...
imports = [
    ("UrlDownloadToFile", "urlmon.dll", "File Download"),
    ("GetAsyncKeyState", "user32.dll", "Keylogging"),
    ("CreateProcess", "kernel32.dll", "RunPE"),
    ("GetThreadContext", "kernel32.dll", "RunPE"),
    ("NtUnmapViewOfSection", "kernel32.dll", "RunPE"),
    ("VirtualAllocEx", "kernel32.dll", "RunPE/Injection"),
    ("WriteProcessMemory", "kernel32.dll", "RunPE/Injection"),
    ("SetThreadContext", "kernel32.dll", "RunPE"),
    ("ResumeThread", "kernel32.dll", "RunPE/Injection"),

]

# Suspicious libraries:
# Networking, RunPE stuff, ...
libraries = [
    ("wininet.dll", "Web / Networking"),
    ("winhttp.dll", "Web / Networking"),
    ("ws2_32.dll", "Communication / Networking - Raw sockets"),
    ("urlmon.dll", "Web / Networking"),
    ("advapi32.dll", "Registry access")
]

with sqlite3.connect("susp.db") as dbconn:
    cursor = dbconn.cursor()
    # Create tables
    cursor.execute("CREATE TABLE process(name, path, desc)")
    cursor.execute("CREATE TABLE import(name, lib, desc)")
    cursor.execute("CREATE TABLE library(name, desc)")
    # Populate tables
    cursor.executemany("INSERT INTO process VALUES(?, ?, ?)", systemProcs)
    dbconn.commit() # commit after each insertion
    cursor.executemany("INSERT INTO import VALUES(?, ?, ?)", imports)
    dbconn.commit()
    cursor.executemany("INSERT INTO library VALUES(?, ?)", libraries)
    dbconn.commit()
    
