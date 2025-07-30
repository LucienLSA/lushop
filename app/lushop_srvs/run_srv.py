import os
import asyncio

services = [
    r"E:\Study\Proj\lushop\app\lushop_srvs\user_srv\main.exe",
    r"E:\Study\Proj\lushop\app\lushop_srvs\userop_srv\main.exe",
    r"E:\Study\Proj\lushop\app\lushop_srvs\goods_srv\main.exe",
    r"E:\Study\Proj\lushop\app\lushop_srvs\order_srv\main.exe",
    r"E:\Study\Proj\lushop\app\lushop_srvs\inventory_srv\main.exe",
]

# user_web = "E:\\zhuomian\\project\\mxshop\\web_golang\\order-web\\main.go"
# userop_web = "E:\\zhuomian\\project\\mxshop\\web_golang\\userop-web\\main.go"
# goods_web = "E:\\zhuomian\\project\\mxshop\\web_golang\\goods-web\\main.go"
# order_web = "E:\\zhuomian\\project\\mxshop\\web_golang\\order-web\\main.go"
# oss_web = "E:\\zhuomian\\project\\mxshop\\web_golang\\oss-web\\main.go"

for exe in services:
    exe_dir = os.path.dirname(exe)
    exe_file = os.path.basename(exe)
    os.system(f'start cmd /k "cd /d {exe_dir} && {exe_file}"')

print("所有服务已在新窗口启动！")