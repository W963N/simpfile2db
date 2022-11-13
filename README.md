# W963N Memory ~~ DB REG T00l

Setup:

 - `https://github.com/W963N/simpfile2db.git`
 - `cd simpfile2db`
 - `go build .`
 - `vi /your/path/env.toml`

Usage:

 - Store data: `simpfile2db -t [DB NAME] -r -f [KEY(FILE NAME)].[Ext]` 
 - Search data: `simpfile2db -t [DB NAME] -s -k [KEY]`
 - Delete data: `simpfile2db -t [DB NAME] -d -k [KEY]`
 - Output data: `simpfile2db -t [DB NAME] -o -k [KEY]`

Options:

```
  -d	Delete file.  
  -e string  
    	path of env.toml. (default "./env.toml")  
  -f string  
    	file path  
  -k string  
    	key  
  -o	Output file  
  -r	Register file  
  -s	Search file  
  -t string  
    	Use db name.  
  -v string  
    	Select types(info, warn, error). (default "info")  
```

## Use case

1. Register hello.txt.
  - `bind -t txt -r -f hello.txt`
  - `bind -t txt -s -k hello`
2. Restore hello.txt.
  - `bind -t txt -o -k hello`

Hobbyright 2022 walnut 🐿🐿🐿 .
