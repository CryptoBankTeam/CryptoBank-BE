require("dotenv").config();
const { ethers } = require("ethers");
const { Client } = require("pg");

const provider = new ethers.providers.JsonRpcProvider(process.env.ALCHEMY_URL);
const signer = new ethers.Wallet(process.env.PRIVATE_KEY, provider);

const abi = [ "function batchCheckOverdue(uint256[] loanIds)" ];
const contract = new ethers.Contract(process.env.CONTRACT_ADDRESS, abi, signer);

const db = new Client({
  user: process.env.DB_USER,
  host: process.env.DB_HOST,
  database: process.env.DB_NAME,
  password: process.env.DB_PASSWORD,
  port: process.env.DB_PORT
});

db.connect();

async function checkOverdueLoans() {
  try {
    const now = Math.floor(Date.now() / 1000);

    const res = await db.query(`
      SELECT id FROM loans
      WHERE status = 1 AND due_date < $1
    `, [now]);

    const overdueLoanIds = res.rows.map(row => Number(row.id));

    if (overdueLoanIds.length === 0) {
      console.log("Нет просроченных займов");
      return;
    }

    const chunk = overdueLoanIds.slice(0, 20);
    const tx = await contract.batchCheckOverdue(chunk);
    await tx.wait();
    console.log("Смарт-контракт обновил статусы:", chunk);

    for (const id of chunk) {
      await db.query(`UPDATE loans SET status = 3 WHERE id = $1`, [id]);
    }
  } catch (err) {
    console.error("Ошибка при проверке просрочек:", err);
  }
}

setInterval(checkOverdueLoans, 60 * 1000); 
