# Rental Property API

Go + Beego REST API serving rental properties from an in-memory JSON file.

## Run

```bash
go mod download
bee run
```

The server listens on the `httpport` in `conf/app.conf` (8080), and the data file path comes from `datafile`.

If `bee` is not installed:

```bash
go install github.com/beego/bee/v2@latest
```

then `bee run` will work.

Alternatively, since the mod file has `bee` declared as a tool:

```bash
go tool bee run
```

## Endpoints

```bash
# List (with filters and limit)
curl "http://localhost:8080/v1/properties?feed=11&published=false&min_price=50&max_price=150&property_type=Hotel&amenities=Internet,Parking&limit=5"

# Get by ID
curl http://localhost:8080/v1/properties/BC-1000001
```

## Test

```bash
go test ./... -v
go vet ./...
```
