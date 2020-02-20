# network-probe

----

* [ ] ip
* [ ] http/https
* [ ] domain

### 启动开发环境 

```bash

docker-compose -f ./develop-docker-compose.yml build && docker-compose -f ./develop-docker-compose.yml up -d

```

### 构建编译可执行文件

```bash
docker-compose -f ./develop-docker-compose.yml run --rm develop bash -c "cd /work/src;go build -o /work/bin/service"

```

### 打包服务的docker镜像

```bash
docker build ./produce -t tomhjx/network-probe:latest
docker push tomhjx/network-probe:latest
```


### 启动服务

```bash

docker-compose -f ./produce-docker-compose.yml build && docker-compose -f ./produce-docker-compose.yml up -d
```

### 配置

配置方式  | 参数名                    |参数类型    | 备注
---------|--------------------------|----------|-----
环境变量  | PROBE_INTERVAL_SECOND    |int       | 探测间隔，单位秒，默认为10秒
环境变量  | TARGET_SOURCE_URL        |string    | 目标地址列表（ip、域名、url）接口地址，以换行符（\n）分隔，每1分钟从该地址获取目标列表更新到内存