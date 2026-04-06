## 23/2/2026
### Auth refactor+setup
 - config.go: takes secret key and return data structure of config (exprire time, name domain, ...)
 - jwt.go: token pair (access and refresh), generate token pairs using user credential to create a JWT
    - JWT creates by: *jwt.NewWithClaims()* 
    - JWT stored by: *jwt.MapClaims*
    - Validate token: Parse token and return if ok
 - middleware.go: require auth: extract bearer token from cookie ,validate token from cookie, extract claim from token cookie.