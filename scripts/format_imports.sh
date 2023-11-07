# Step1
# goimports-reviser -rm-unused -project-name "github.com/huaweicloud/terraform-provider-hcso" -company-prefixes "github.com/chnsz/golangsdk,github.com/huaweicloud/huaweicloud-sdk-go-v3,github.com/huaweicloud/terraform-provider-huaweicloud"  -imports-order "std,general,company,project,blanked,dotted"  -format ./...

file_list=$(find ./ -type f -name "*.go")

for file in $file_list
do
  echo "$file"
  python3 scripts/group_import.py $file "github.com/chnsz/golangsdk"
  
  python3 scripts/group_import.py $file "github.com/huaweicloud/huaweicloud-sdk-go-v3"
  # 在这里执行其他操作
done