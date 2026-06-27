const { Client } = require('pg');

const client = new Client({
  connectionString: 'postgres://postgres.mqolpawlannysqjokzoq:Chia.Florist@21@aws-1-ap-northeast-2.pooler.supabase.com:6543/postgres?sslmode=disable'
});

async function check() {
  await client.connect();
  
  const shopsRes = await client.query('SELECT id, name FROM shops');
  console.log('Shops:', shopsRes.rows);
  
  const shopCouriersRes = await client.query('SELECT * FROM shop_couriers');
  console.log('Shop Couriers:', shopCouriersRes.rows);
  
  await client.end();
}

check().catch(console.error);
