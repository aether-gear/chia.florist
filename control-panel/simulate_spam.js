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
        port: 8080, // Updated to Go Backend Port
        path: '/search?q=spam_request_' + count,
        method: 'GET',
        headers: {
            'X-Forwarded-For': ip, // Updated header for Go WAF
            'User-Agent': 'SpamBot/1.0'
        }
    };

    const req = http.request(options, (res) => {
        // console.log(`[${ip}] Spam #${count} sent. Status: ${res.statusCode}`);
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

        if (sent >= 50) {
            clearInterval(interval);
            console.log(`[${ip}] Finished spamming.`);
        }
    }, 20);
});
