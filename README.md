Backs up (saves) modified dates (timestamps) for all files including those in subfolders in a directory with the ability to later restore them from the saved state, might be useful for cloud backup software and services which don't support restoring file timestamps when getting your files back

**Requirements:**  
Python (2\3+ version)  

**Usage:**  
`backup-file-timestamps.py` - uses current directory to go and save timestamps  
`backup-file-timestamps.py -save "C:\my folder"` - save modified time dates inside `C:\my folder` and its subfolders  
`backup-file-timestamps.py -restore "C:\my folder"` - restores modified time dates for all files inside `C:\my folder` and its subfolders  


**Console-less usage (GUI):**

Navigate to `shell:opento` (enter into explorer's address bar)  
![image](https://user-images.githubusercontent.com/1309656/123554015-ce337380-d786-11eb-88bc-48a8c214a88d.png)

and place two `.bat` files there (or shortcuts to them, doesn't matter) after editing those with notepad to have the correct path pointing to your `backup-file-timestamps.py` file where you placed it after saving)  
![image](https://user-images.githubusercontent.com/1309656/123554153-6af61100-d787-11eb-8e88-686efc11967a.png)

after that you can simply click "send to" in context menu on a folder  
![image](https://user-images.githubusercontent.com/1309656/123554307-07201800-d788-11eb-954f-e65aa21a0b11.png)
