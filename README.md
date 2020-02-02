# network-probe

----

* [ ] ip
* [ ] http/https
* [ ] domain



### 构建编译可执行文件

```bash
docker-compose -f ./develop-docker-compose.yml run --rm develop bash -c "cd /work/src;go build -o /work/bin/service"

```

### 启动服务

```bash

docker-compose -f ./develop-docker-compose.yml run --rm develop bash -c "/work/bin/service"
```