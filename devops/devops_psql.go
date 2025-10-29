package devops

const (
	psqlDropTable           = `DROP TABLE IF EXISTS "public"."%s";`
	psqlCreateDatabase      = `CREATE DATABASE %s WITH LC_COLLATE = 'C' LC_CTYPE = 'en_US.utf8' TABLESPACE = pg_default;`
	psqlBeforeInstallTables = `
DEALLOCATE ALL;       -- 清除当前会话的所有预处理语句
DISCARD PLANS;        -- PostgreSQL ≥13 替代方案`
)
