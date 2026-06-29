const http = require('http');

const targetUrl = 'http://localhost:8080';

const testPayloads = [
    {
        name: "Standard SQL Injection (Control)",
        path: "/search?q=%27%20OR%201%3D1%20--",
        desc: "Plain SQL injection query parameter"
    },
    {
        name: "Double URL Encoded SQL Injection",
        path: "/search?q=%27%2520OR%25201%253D1%2520--", // ' %20OR%201%253D1%20-- (double encoded)
        desc: "Tests WAF's triple URL decoding resilience"
    },
    {
        name: "SQL Comment Obfuscation",
        path: "/search?q=1%20UNION/*/*/SELECT%20username,password%20FROM%20users",
        desc: "Injects inline SQL comments to break basic pattern match"
    },
    {
        name: "Mixed-Case Script Injection (XSS)",
        path: "/search?q=%3CjAvAsCrIpT%3Ealert(1)%3C/jAvAsCrIpT%3E",
        desc: "Alters script tag capitalization to bypass case-sensitive checks"
    },
    {
        name: "Hexadecimal-Encoded Script Payload",
        path: "/search?q=%3cscript%3ealert(String.fromCharCode(0x41,0x42,0x43))%3c/script%3e",
        desc: "Uses JavaScript String.fromCharCode in Hex format to bypass keyword filters"
    },
    {
        name: "Null Byte Injection Bypass Attempt",
        path: "/search?q=%00%27%20OR%201%3D1%20--",
        desc: "Prepends null byte %00 to truncate WAF search string"
    }
];

const sendRequest = (payload) => {
    return new Promise((resolve) => {
        const options = {
            hostname: 'localhost',
            port: 8080,
            path: payload.path,
            method: 'GET',
            headers: {
                'User-Agent': 'WAFNinja-Sim/1.0',
                'x-simulated-ip': '88.99.100.54' // spoof test IP
            }
        };

        const req = http.request(options, (res) => {
            let data = '';
            res.on('data', chunk => { data += chunk; });
            res.on('end', () => {
                resolve({
                    name: payload.name,
                    path: payload.path,
                    statusCode: res.statusCode,
                    passed: res.statusCode === 200
                });
            });
        });

        req.on('error', (e) => {
            resolve({
                name: payload.name,
                path: payload.path,
                statusCode: 500,
                passed: false,
                error: e.message
            });
        });
        req.end();
    });
};

async function runTests() {
    console.log(`
===================================================
   🛡️  WAF EVASION & OBFUSCATION AUDIT SIMULATOR
===================================================
Target: http://localhost:8080
Testing for common WAF bypass techniques...
`);

    let blocked = 0;
    let bypassed = 0;

    for (const payload of testPayloads) {
        console.log(`Testing: ${payload.name}`);
        console.log(`- Path  : ${payload.path}`);
        console.log(`- Intent: ${payload.desc}`);
        
        const result = await sendRequest(payload);
        
        if (result.passed) {
            console.log(`🔴 RESULT: BYPASSED (Status: 200 OK) - WAF failed to intercept!`);
            bypassed++;
        } else {
            console.log(`` + (result.statusCode === 403 ? '🟢' : '🟡') + ` RESULT: BLOCKED (Status: ${result.statusCode}) - WAF successfully intercepted.`);
            blocked++;
        }
        console.log(`---------------------------------------------------`);
    }

    console.log(`
===================================================
             📊 AUDIT SUMMARY REPORT
===================================================
Total Evasion Attempts : ${testPayloads.length}
Successfully Blocked   : ${blocked}
Successful Bypasses    : ${bypassed}
WAF Shield Integrity   : ${Math.round((blocked / testPayloads.length) * 100)}%
===================================================
`);
}

runTests();
