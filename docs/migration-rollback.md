# 数据库迁移回滚说明

本文档用于说明 `cmd/migrate` 的回滚策略、常见失败场景与处理步骤。

## 1. 迁移命令

```bash
# 应用所有未执行迁移
go run cmd/migrate/main.go up

# 回滚最近 1 个版本（默认）
go run cmd/migrate/main.go down

# 回滚最近 N 个版本
go run cmd/migrate/main.go down 2

# 回滚全部已执行版本
go run cmd/migrate/main.go down all
```

## 2. 执行机制

1. 迁移文件目录：`internal/migrations`
2. 文件命名规则：`<version>_<name>.up.sql` / `<version>_<name>.down.sql`
3. 执行记录表：`schema_migrations`
4. 每个版本迁移在独立事务中执行：
   - `up` 成功后写入 `schema_migrations`
   - `down` 成功后删除 `schema_migrations` 记录
5. 任一 SQL 失败将回滚当前版本事务，不影响已成功版本。

## 3. 失败回滚流程（推荐）

1. 立即停止发布流程，冻结应用写流量。
2. 执行：

```bash
go run cmd/migrate/main.go down
```

3. 确认 `schema_migrations` 最新版本已回退。
4. 验证关键表结构与索引状态。
5. 修复迁移文件后重新执行 `up`。

## 4. 常见问题

### 4.1 提示 “missing down migration file”

原因：存在 `up` 文件但缺少配套 `down` 文件。  
处理：补齐同版本 `down` 文件后再执行。

### 4.2 提示 “migration name mismatch”

原因：同版本的 `up/down` 文件 `<name>` 不一致。  
处理：统一文件名中的 `<name>` 段后重试。

### 4.3 迁移执行中断后状态不一致

原因：手工执行 SQL 或中断导致表结构与 `schema_migrations` 记录不一致。  
处理：

1. 对比目标版本 SQL 与数据库实际结构。
2. 必要时手动修正数据库对象。
3. 调整 `schema_migrations` 到真实状态。
4. 再执行 `up/down` 使版本恢复一致。

## 5. 生产建议

1. 在生产执行迁移前必须完成备份。
2. 高风险迁移（删除列、重建索引）需先在预发环境演练。
3. 每次发布仅包含少量迁移版本，降低回滚复杂度。
4. 迁移窗口内监控数据库连接数、锁等待、慢查询。
