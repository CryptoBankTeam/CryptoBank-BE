require("dotenv").config();
const { ethers } = require("ethers");
const { Client } = require("pg");

// 🔌 Ethereum провайдер и контракт
const provider = new ethers.JsonRpcProvider(process.env.ALCHEMY_URL);
const contractAddress = process.env.CONTRACT_ADDRESS;
const abi = [
  "event LoanCreated(uint256 indexed loanId, address indexed lender, uint256 amount, uint256 interest, uint256 collateral, uint256 duration, uint8 status)",
  "event LoanAccepted(uint256 indexed loanId, address indexed borrower, uint8 status)",
  "event LoanRepaid(uint256 indexed loanId, address indexed borrower, uint8 status)",
  "event CollateralClaimed(uint256 indexed loanId, address indexed lender, uint8 status)",
  "event LoanOverdue(uint256 indexed loanId, address indexed lender, uint8 status)"
];
const contract = new ethers.Contract(contractAddress, abi, provider);

//PostgreSQL подключение
const db = new Client({
  user: process.env.DB_USER,
  host: process.env.DB_HOST,
  database: process.env.DB_NAME,
  password: process.env.DB_PASSWORD,
  port: process.env.DB_PORT
});

db.connect().then(() => console.log("DB connected")).catch(console.error);

// Пересчёт рейтинга
async function recalculateRating(userId) {
  const res = await db.query(
    `SELECT clean_loans, overdue_loans, offers_accepted FROM users WHERE id = $1`,
    [userId]
  );

  if (res.rows.length === 0) return;

  const { clean_loans, overdue_loans, offers_accepted } = res.rows[0];
  const total = clean_loans + overdue_loans + offers_accepted;

  if (total === 0) return;

  const rating = parseFloat(((clean_loans * 2 + offers_accepted) / total).toFixed(2));

  await db.query(`UPDATE users SET rating = $1 WHERE id = $2`, [rating, userId]);
  console.log(`⭐ Rating updated for user #${userId}:`, rating);
}

// LoanCreated
contract.on("LoanCreated", async (loanId, lender, amount, interest, collateral, duration, status) => {
  try {
    const res = await db.query("SELECT id FROM users WHERE adress_wallet = $1 LIMIT 1", [lender]);
    if (res.rows.length === 0) return;

    const lenderId = res.rows[0].id;
    await db.query(`
      INSERT INTO loans (id, amount, interest, collateral, duration, status, lender_id)
      VALUES ($1, $2, $3, $4, $5, $6, $7)
      ON CONFLICT (id) DO NOTHING
    `, [loanId, amount, interest, collateral, duration, status, lenderId]);

    console.log(`LoanCreated #${loanId} by lender ID ${lenderId}`);
  } catch (err) {
    console.error("LoanCreated error:", err);
  }
});

// LoanAccepted
contract.on("LoanAccepted", async (loanId, borrower, status) => {
  try {
    const res = await db.query("SELECT id FROM users WHERE adress_wallet = $1 LIMIT 1", [borrower]);
    if (res.rows.length === 0) return;

    const borrowerId = res.rows[0].id;
    const loan = await db.query("SELECT lender_id FROM loans WHERE id = $1", [loanId]);
    const lenderId = loan.rows[0]?.lender_id;

    await db.query(`
      UPDATE loans
      SET borrower_id = $1, status = $2
      WHERE id = $3
    `, [borrowerId, status, loanId]);

    if (borrowerId !== lenderId) {
      await db.query(`UPDATE users SET offers_accepted = offers_accepted + 1 WHERE id = $1`, [lenderId]);
      await recalculateRating(lenderId);
    }

    console.log(`LoanAccepted #${loanId} by borrower ID ${borrowerId}`);
  } catch (err) {
    console.error("LoanAccepted error:", err);
  }
});

// LoanRepaid
contract.on("LoanRepaid", async (loanId, borrower, status) => {
  try {
    await db.query("UPDATE loans SET status = $1 WHERE id = $2", [status, loanId]);

    const res = await db.query("SELECT lender_id FROM loans WHERE id = $1", [loanId]);
    const lenderId = res.rows[0]?.lender_id;
    if (lenderId) {
      await db.query(`UPDATE users SET clean_loans = clean_loans + 1 WHERE id = $1`, [lenderId]);
      await recalculateRating(lenderId);
    }

    console.log(`LoanRepaid #${loanId} by borrower`);
  } catch (err) {
    console.error("LoanRepaid error:", err);
  }
});

// CollateralClaimed
contract.on("CollateralClaimed", async (loanId, lender, status) => {
  try {
    await db.query("UPDATE loans SET status = $1 WHERE id = $2", [status, loanId]);

    const res = await db.query("SELECT lender_id FROM loans WHERE id = $1", [loanId]);
    const lenderId = res.rows[0]?.lender_id;
    if (lenderId) {
      await db.query(`UPDATE users SET overdue_loans = overdue_loans + 1 WHERE id = $1`, [lenderId]);
      await recalculateRating(lenderId);
    }

    console.log(`CollateralClaimed #${loanId} → lender #${lenderId}`);
  } catch (err) {
    console.error("CollateralClaimed error:", err);
  }
});

// LoanOverdue
contract.on("LoanOverdue", async (loanId, lender, status) => {
  try {
    await db.query("UPDATE loans SET status = $1 WHERE id = $2", [status, loanId]);

    const res = await db.query("SELECT id FROM users WHERE adress_wallet = $1 LIMIT 1", [lender]);
    if (res.rows.length > 0) {
      const userId = res.rows[0].id;
      await db.query(`UPDATE users SET overdue_loans = overdue_loans + 1 WHERE id = $1`, [userId]);
      await recalculateRating(userId);
    }

    console.log(`LoanOverdue #${loanId} → lender ${lender}`);
  } catch (err) {
    console.error("LoanOverdue error:", err);
  }
});
