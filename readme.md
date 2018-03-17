# gobalance

## TODO
- [x] Init from config
- [x] Pool manager and workers
- [x] CLI
- [x] Development pool
- [x] Round robin
- [x] Header adjustments
- [ ] Documentation
- [ ] CI pipeline (build, test, coverage)
- [ ] Header modification
- [ ] Statistics
- [ ] Logging
- [ ] SSL termination
- [ ] Least connection
- [ ] Ip hash
- [ ] REST API
- [ ] Management UI
- [ ] IP blacklisting/whitelisting

## Development
* To spin up downstream test servers, open a separate terminal and run ```go run test/testpool.go```
* This will create two local endpoints with catch all route handlers with the following information
    * Methods: GET, POST, PUT, DELETE
    * Ports: 8080, 8081
* Ex:
    * ```curl localhost:8080/products -X GET```
    * ```curl localhost:8081/animal/1 -X GET```
    * ```curl localhost:8080/books -X POST```
    * ```curl localhost:8081/games -X POST```