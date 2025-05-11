```
 ______     __     ______   __   __   ______     __  __     __         ______  
/\  ___\   /\ \   /\__  _\ /\ \ / /  /\  __ \   /\ \/\ \   /\ \       /\__  _\ 
\ \ \__ \  \ \ \  \/_/\ \/ \ \ \'/   \ \  __ \  \ \ \_\ \  \ \ \____  \/_/\ \/ 
 \ \_____\  \ \_\    \ \_\  \ \__|    \ \_\ \_\  \ \_____\  \ \_____\    \ \_\ 
  \/_____/   \/_/     \/_/   \/_/      \/_/\/_/   \/_____/   \/_____/     \/_/ 
                                                                               
```

# Intro 
GitVault is a cli tool used to manage "vault" files stored in a git repo.  

# Install 

## Go 
```
go install github.com/pzolo85/git_vault
```

# Usage 
```
$ git_vault -h
Run inside a a folder with an initialized git repo. 

        Add all your files inside a folder named "open".  
        Use sub-commands "open", "close", "push" and "pull" to update files.

Usage:
  git_vault [command]

Available Commands:
  close       close your vault
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  open        open a *.tgz.enc file and place the content inside a folder named 'open'
  pull        TODO: implement pull
  push        TODO: implement push

Flags:
  -h, --help     help for git_vault
  -t, --toggle   Help message for toggle

Use "git_vault [command] --help" for more information about a command.
```
