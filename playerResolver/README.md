### Player Resolver

Tool for manual updating player names and avatars in mongodb

For configuration use same `config.yaml` as PickupStats

1. Build
```bash
go build -o bin/playerResolver ./playerResolver
```

2. Run script
```bash
./bin/playerResolver --config <config path>
```