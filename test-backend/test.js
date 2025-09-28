const axios = require("axios");

const API_BASE = "http://localhost:8082/api/admin";
const ADMIN_TOKEN = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTkwNzkxMzcsInJvbGUiOiJhZG1pbiIsInR5cGUiOiJhY2Nlc3MiLCJ1c2VyX2lkIjoxfQ.x7Bq-wZKGs6C9YwcaEeHpgZ4-fhD8VgbFWN-I-efdy4"; // replace with your admin token
const SCHEDULE_ID = 3;      // schedule ID to fetch
const SCHEDULE_SNACK_ID = 3; // specific snack ID within the schedule

async function main() {
  try {
    // ---------------- GET SCHEDULE ----------------
    const scheduleRes = await axios.get(
      `${API_BASE}/schedules/${SCHEDULE_ID}`,
      { headers: { Authorization: `Bearer ${ADMIN_TOKEN}` } }
    );
    console.log("🎬 Schedule details:", scheduleRes.data);

    // ---------------- GET ALL SNACKS FOR SCHEDULE ----------------
    const snacksRes = await axios.get(
      `${API_BASE}/schedules/${SCHEDULE_ID}/snacks`,
      { headers: { Authorization: `Bearer ${ADMIN_TOKEN}` } }
    );
    console.log("🍿 All schedule snacks:", snacksRes.data.schedule_snacks);

    // ---------------- GET SPECIFIC SCHEDULE SNACK ----------------
    const singleSnackRes = await axios.get(
      `${API_BASE}/schedules/${SCHEDULE_ID}/snacks/${SCHEDULE_SNACK_ID}`,
      { headers: { Authorization: `Bearer ${ADMIN_TOKEN}` } }
    );
    console.log("🥨 Specific schedule snack:", singleSnackRes.data);

  } catch (err) {
    console.error("❌ Error:", err.response?.data || err.message);
  }
}

main();
