# file-org

# Overview

Copy files from one directory to another. 

- Walk a source directory for file types defined in a task json config file.


```json
{
    "tasks": [
        {
            "enabled": true,
            "recursive": true,
            "filetype": [".jpeg"],
            "fileprefix":"fo_",
            "scriptprefix":"fo_jpeg_",
            "sourcepath": "/Users/username",
            "destinationpath":"./test",
            "scriptpath":"./scripts"
        }
    ]
}
```