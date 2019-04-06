My first experiment into the PROS system.

What I did to get started:
1. Downloaded `pros-windows-3.1.4.0-64bit.exe` from https://github.com/purduesigbots/pros-cli/releases.
2. Installed PROS 3.1.4.0 for all users.  This also put the pros-cli binary, `prosv5.exe`, in the system PATH.
3. Restarted my terminal application, entered this directory, and ran:
    ```
    prosv5 conduct new .
    ```
    This created a full-blown PROS "Hello World" programmin environment with about 800-900 different
    files, and then compiled the whole thing.
4. I versioned the .gitignore first (because most of the files from the previous step are 
   automatically-generated and never change.)  After committing that, I versioned the rest of the
   files, and there were only 5 of them:

    include/main.h
    project.pros
    src/autonomous.cpp
    src/initialize.cpp
    src/opcontrol.cpp
5. I built and deployed the code with:
    ```
    prosv5 mut --slot=2
    ```
    The slots are 1-based.

    The command automatically attachd to the terminal in order to capture standard output,
    but I could also have hit CTRL+C (to exit from `prosv5 mut`) and then run:

    ```
    prosv5 terminal
    ```
    to attach to a running program.