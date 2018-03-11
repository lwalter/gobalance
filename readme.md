## TODO
- [x] Init from config
- [x] Pool manager and workers
- [x] CLI
- [x] Development pool
- [ ] Round robin
- [ ] Documentation
- [ ] Header modification
- [ ] Statistics
- [ ] Logging
- [ ] SSL termination
- [ ] Least connection
- [ ] Ip hash
- [ ] REST API
- [ ] Management UI

## Development
* To spin up downstream test servers, open a separate terminal and run ```go run test/testpool.go```
* This will create two local endpoints with the following information
    * Route: "/"
    * Methods: GET, POST, PUT, DELETE
    * Ports: 4000, 4001
* Ex:
    * ```curl localhost:4000/ -X GET```
    * ```curl localhost:4001/ -X GET```
    * ```curl localhost:4000/ -X POST```
    * ```curl localhost:4001/ -X POST```