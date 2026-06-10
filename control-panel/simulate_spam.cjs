const http = require('http');

const spammers = [
    '71.6.146.185',
    '80.82.77.139',
    '61.146.163.4',
    '70.34.245.215'
];

const sendSpam = (ip, count) => {
    const options = {
        hostname: 'localhost',
        port: 8080, // Go Backend Port
        path: '/search?q=spam_request_' + count,
        method: 'GET',
        headers: {
            'X-Forwarded-For': ip, // Client IP spoofing
            'User-Agent': 'SpamBot/1.0'
        }
    };

    const req = http.request(options, (res) => {
        if (res.statusCode === 429) {
            console.log(`\n🚨 [${ip}] Rate Limit Triggered! Server returned: ${res.statusCode} Too Many Requests`);
        }
        res.on('data', () => { });
    });

    req.on('error', (e) => { });
    req.end();
};

console.log("🚀 Launching Spam/DDoS Simulation (Target: Golang WAF :8080)...");

spammers.forEach((ip) => {
    let sent = 0;
    const interval = setInterval(() => {
        sendSpam(ip, sent++);
        if (sent % 10 === 0) console.log(`[${ip}] Sent ${sent} spam requests...`);

        if (sent >= 70) {
            clearInterval(interval);
            console.log(`[${ip}] Finished spamming.`);
        }
    }, 15);
});
