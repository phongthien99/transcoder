#!/bin/bash

cd "$1" || exit 1

# Lấy đường dẫn thư mục từ đối số
target_directory="$1"

# Kiểm tra xem thư mục types có tồn tại không
if [ ! -d "$target_directory/types" ]; then
    echo "Thư mục 'types' không tồn tại. Khởi tạo..."
    mkdir -p "types"
fi

# Chạy lệnh yaml-to-go trong thư mục này
/usr/bin/yaml-to-go/yaml-to-go -i default.yaml -o types/type.go

# ./scripts/gen_type_config.sh apps/example/src/config