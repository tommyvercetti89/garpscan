param (
    [Parameter(Mandatory=$true)]
    [ValidateSet("test", "build", "example", "tidy", "all")]
    $Task
)

switch ($Task) {
    "test" {
        Write-Host "Running tests..." -ForegroundColor Cyan
        go test -v -race ./...
    }
    "build" {
        Write-Host "Building..." -ForegroundColor Cyan
        go build -v ./...
    }
    "example" {
        Write-Host "Running example..." -ForegroundColor Cyan
        go run examples/basic_scan/main.go
    }
    "tidy" {
        Write-Host "Tidying module dependencies..." -ForegroundColor Cyan
        go mod tidy
    }
    "all" {
        Write-Host "Running all tasks..." -ForegroundColor Cyan
        go mod tidy
        go test -v -race ./...
        go build -v ./...
    }
}
