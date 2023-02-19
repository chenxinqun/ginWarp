#### ginWarp-example.users 

| 序号 | 名称 | 描述 | 类型 | 键 | 为空 | 额外 | 默认值 |
| :--: | :--: | :--: | :--: | :--: | :--: | :--: | :--: |
| 1 | id | 用户ID | bigint(20) | PRI | NO |  | 0 |
| 2 | tenant_id | 企业ID | bigint(20) |  | NO |  | 0 |
| 3 | login_time | 登陆时间 | datetime |  | YES |  |  |
| 4 | created_at | 创建时间 | datetime |  | YES |  |  |
| 5 | updated_at | 更新时间 | datetime |  | YES |  |  |
| 6 | deleted_at | 删除时间 | datetime |  | YES |  |  |
| 7 | account | 账号 | varchar(255) | UNI | NO |  |  |
| 8 | password | 密码 | varchar(255) |  | NO |  |  |
| 9 | login_ip | 登陆IP | varchar(255) |  | NO |  |  |
