$ErrorActionPreference = 'Stop'

Write-Host 'Turning off Go module mode for build...'
$env:GO111MODULE = 'off'

Write-Host 'Installing dependencies...'
go get github.com/lib/pq@v1.10.9
Go get github.com/golang-jwt/jwt/v5@v5.0.0
Go get golang.org/x/crypto@v0.14.0
Go get github.com/joho/godotenv@v1.5.1
Go get github.com/skip2/go-qrcode@v0.0.0-20200617195104-da3256ba3586
Go get github.com/signintech/gopdf@v0.9.13
Go get github.com/go-redis/redis/v8@v8.11.5
Go get github.com/stripe/stripe-go/v75@v75.2.0

Write-Host 'Building application...'
go build -o booking-system
Write-Host 'Build complete: .\booking-system.exe'