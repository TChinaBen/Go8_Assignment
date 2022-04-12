### 部署
REPO={{REPO}} \
USERNAME={{USERNAME}} \
PASSWORD={{PASSWORD}} \
MODE=fs \
HARBOR_URL="172.18.8.210:5000/library" \
./deploy.sh

MODE=fs 本地模式
MODE=s3 Ceph模式

### 删除
./delete.sh
### 重试
./retry.sh


