Backs up (saves) modified dates (timestamps) for all files including those in subfolders in a directory with the ability to later restore them from the saved state, might be useful for cloud backup software and services which don't support restoring file timestamps when getting your files back

Get executable files from the [**Releases**](https://github.com/Owyn/backup_file_timestamps/releases) section

**Usage:**  
`backup-file-timestamps.exe` - uses current directory (or the one dragged onto it) to go and save timestamps  
`backup-file-timestamps.exe -save "C:\my folder"` - save modified time dates inside `C:\my folder` and its subfolders  
`backup-file-timestamps.exe -restore "C:\my folder"` - restores modified time dates for all files inside `C:\my folder` and its subfolders  


**Console-less usage (GUI):**

Navigate to `shell:sendto` (enter into explorer's address bar or the Run command)  
![shell-sendto-1](https://github.com/user-attachments/assets/5e82a84d-a15f-4dcb-813c-63dd4ac2c00d) or ![TjL86](https://github.com/user-attachments/assets/a79f55b8-d801-424b-89fd-380751b084a9)


and place two files executable files there - `backup-file-timestamps.exe` and `restore file timestamps.bat` (or shortcuts to them, doesn't matter)   


now after that you can simply click "send to" in context menu on any folder  
![context_menu translated](https://github.com/user-attachments/assets/200c6db1-8742-43fa-960b-63837d4ba2d0)

**Screenshots**

Saving:  
![done](https://github.com/user-attachments/assets/c16eb7ab-1555-4c60-a228-3443da3bfb51)  
Restoring:  
![done_restore](https://github.com/user-attachments/assets/9c063f22-15ea-4ec6-b0a5-8c4cf89702f9)  

