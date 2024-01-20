# update-proton-ge

An updater and installer for proton-ge on linux systems

## Usage

```
Usage of ./update-proton-ge:
  -force
     force download even when up-to-date
  -l string
     log level (trace, debug, info, warn, error, fatal, panic) (default "info")
  -y skip confirmations
```

## Configuration

By default the install locations are:

- archived versions: `~/.local/proton-ge`
- compatiblity tools: `~/.steam/root/compatibilitytools.d`

But you can the following environment variables to override these locations: 

- `PROTON_GE_ARCHIVE`
- `COMPATIBILITY_TOOLS_DIR`

## Roadmap

- Get automatic versioning and a pipeline up and ready
- Unit test the repo
- More clean up

