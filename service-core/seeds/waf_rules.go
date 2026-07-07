package seeds

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

const WAFRULES_SEED_NAME = "waf_rule_v1"
const WAFRULES_SEED_VERSION = "1.0"

func SeedWafRules(ctx context.Context, pool *pgxpool.Pool) error {
	seeded, err := wafRulesAlreadySeeded(ctx, pool)
	if err != nil {
		return err
	}

	if seeded {
		log.Println("database: waf rules already seeded, skipping")
		return nil
	}

	log.Println("database: seeding waf rules")

	query := `
		INSERT INTO waf_rules
			(id, description, pattern, tags, impact, enabled, created_at, updated_at)
		VALUES
			('1', 'HTML breaking injections including whitespace attacks', '(?:"[^"]*[^-]?>)|(?:[^\w\s]\s*\/>)|(?:>")', ARRAY['xss', 'csrf'], '4', TRUE, NOW(), NOW()),
			('2', 'Attribute breaking injections including whitespace attacks', '(?:"+\.*[<=]\s*"[^"]+")|(?:"\s*\w+\s*=)|(?:>\w=\/)|(?:#.+\)[\"\s]*>)|(?:"\s*(?:src|style|on\w+)\s*=\s*")|(?:[^\"]?\"[,;\s]+\w*[\[\(])', ARRAY['xss', 'csrf'], '4', TRUE, NOW(), NOW()),
			('3', 'Unquoted attribute breaking injections', '(?:^>[\w\s]*<\/?\w{2,}>)', ARRAY['xss', 'csrf'], '2', TRUE, NOW(), NOW()),
			('10', 'Basic directory traversal', '(?:(?:\/|\\)?\.+(\/|\\)(?:\.+)?)|(?:\w+\.exe\??\s)|(?:;\s*\w+\s*\/[\w*-]+\/)|(?:\d\.\dx\|)|(?:%(?:c0\.|af\.|5c\.))|(?:\/(?:%2e){2})', ARRAY['dt', 'id', 'lfi'], '5', TRUE, NOW(), NOW()),
			('12', 'etc/passwd inclusion attempts', '(?:etc\/[^\w\s]*passwd)', ARRAY['dt', 'id', 'lfi'], '5', TRUE, NOW(), NOW()),
			('40', 'MySQL comments, conditions and ch(a)r injections', '(?:\)\s*when\s*\d+\s*then)|(?:"\s*(?:#|--|{))|(?:\/\*!\s?\d+)|(?:ch(?:a)?r\s*\(\s*\d)|(?:(?:(n?and|x?or|not)\s+|\|\||&&)\s*\w+\()', ARRAY['sqli', 'id', 'lfi'], '6', TRUE, NOW(), NOW()),
			('42', 'Classic SQL injection probings 1/2', '(?:"\s*or\s*"?[0-9])|(?:\\x(?:23|27|3d))|(?:^.?"$)|(?:(?:^["\\]*(?:[\d"]+|[^"]+"))+\s*(?:n?and|x?or|not|\|\||&&)\s*[\w"[+&!@(),.-])|(?:[^\w\s]\w+\s*[|-]\s*"\s*\w)|(?:@w+\s+(and|or)\s*["\d]+)|(?:@[\w-]+\s(and|or)\s*[^\w\s])|(?:[^\w\s:]\s*\d\W+[^\w\s]\s*".)|(?:\Winformation_schema|table_name\W)', ARRAY['sqli', 'id', 'lfi'], '6', TRUE, NOW(), NOW()),
			('1001', 'SQL Injection Simple (Union/Select)', '(?i)(union\s+select|select\s+.*\s+from|delete\s+from|drop\s+table|update\s+set)', ARRAY['sqli'], '5', TRUE, NOW(), NOW()),
			('1002', 'XSS Detection (Script Tag)', '(?i)(<script>|javascript:|on\w+=)', ARRAY['xss'], '5', TRUE, NOW(), NOW()),
			('1003', 'SQLi Detection (OR Logic / Tautology)', '(?i)(\s+OR\s+[''\"\w]+=[''\"\w]+)', ARRAY['sqli'], '5', TRUE, NOW(), NOW()),
			('1004', 'Remote Code Execution (RCE) / Command Injection', '(?i)(;\s*ls\s+-\w+|base64_decode|exec\s+.*|system\s*\(|bin\/sh|bin\/bash)', ARRAY['rce', 'ci'], '8', TRUE, NOW(), NOW()),
			('1005', 'Malicious User-Agent / Scanner Detection', '(?i)(sqlmap|nikto|w3af|nmap|acunetix|dirbuster|gobuster)', ARRAY['scanner', 'malicious-ua'], '5', TRUE, NOW(), NOW()),
			('1006', 'Log4Shell Vulnerability Exploit (CVE-2021-44228)', '(?i)(\$\{jndi:(ldap|rmi|ldaps|dns):)', ARRAY['log4j', 'rce'], '9', TRUE, NOW(), NOW()),
			('1007', 'Shellshock Vulnerability Exploit (CVE-2014-6271)', '(?i)(\(\s*\)\s*\{\s*:\s*;\s*\}[;\s]*echo)', ARRAY['shellshock', 'rce'], '9', TRUE, NOW(), NOW()),
			('1008', 'Admin Directory / Local Probe Scan', '(?i)(\/admin|\/wp-admin|\/config|\/setup|\/install)', ARRAY['probe', 'recon'], '3', TRUE, NOW(), NOW()),
			('1009', 'Local File Inclusion (LFI) & Windows Paths', '(?i)(?:etc\/(?:passwd|hosts|shadow|issue|hostname))|(?:boot\.ini|win\.ini|system\.ini)|(?:proc\/self\/environ)|(?:Windows[\\/]System32)|(?:[a-zA-Z]:[\\/](?:windows|winnt|program files|users|etc))', ARRAY['lfi', 'dt'], '7', TRUE, NOW(), NOW()),
			('1010', 'SQL Authentication Bypass & Comments', '(?i)(?:''\s*(?:--|#|/\*))|(?:''\s*;\s*\w+)', ARRAY['sqli'], '6', TRUE, NOW(), NOW()),
			('1011', 'HTML Entity Exploit Obfuscation', '(?i)(&#(?:x[0-9a-fA-F]+|[0-9]+);)', ARRAY['sqli', 'xss', 'lfi', 'obfuscation'], '5', TRUE, NOW(), NOW())
		ON CONFLICT (id) DO NOTHING;
	`

	if _, err := pool.Exec(ctx, query); err != nil {
		return fmt.Errorf("failed to seed waf rules: %w", err)
	}

	if err := markWafRulesSeeded(ctx, pool); err != nil {
		return fmt.Errorf("failed to mark waf rules seed: %w", err)
	}

	return nil
}

func wafRulesAlreadySeeded(ctx context.Context, pool *pgxpool.Pool) (bool, error) {
	var exists bool

	err := pool.QueryRow(
		ctx, `
			SELECT EXISTS (
				SELECT 1 FROM seed_versions WHERE name = $1
			)
		`,
		WAFRULES_SEED_NAME,
	).Scan(&exists)

	return exists, err
}

func markWafRulesSeeded(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(
		ctx,
		`
			INSERT INTO seed_versions (name, version)
			VALUES ($1,$2)
			ON CONFLICT (name) DO NOTHING
		`,
		WAFRULES_SEED_NAME,
		WAFRULES_SEED_VERSION,
	)

	return err
}
