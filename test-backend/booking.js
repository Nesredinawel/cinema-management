const axios = require("axios");

const API_BASE = "http://localhost:8083/api/v1";
const TOKEN = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTkwNzg2MDIsInJvbGUiOiJhZG1pbiIsInR5cGUiOiJhY2Nlc3MiLCJ1c2VyX2lkIjoxfQ.-NJbc4NV7y40NIgWITdh6TcXWILfS7j1dAr26qdKTr0";
const SCHEDULE_TOKEN = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbnRpdHlfdHlwZSI6InNjaGVkdWxlIiwiZXhwIjoxNzU5MTY0MTIyLCJpZCI6MywiaXNzdWVkX2F0IjoxNzU5MDc3NzIyLCJpc3N1ZXIiOiJjaW5lbWEtc2NoZWR1bGluZyJ9._Kx8GE87aiexkvkZKWHG__yEoQ86YqzDoISGhvWnOI4";
const SNACK_TOKEN_1 = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbnRpdHlfdHlwZSI6InNjaGVkdWxlX3NuYWNrIiwiZXhwIjoxNzU5MTY0MTIyLCJpZCI6MywiaXNzdWVkX2F0IjoxNzU5MDc3NzIyLCJpc3N1ZXIiOiJjaW5lbWEtc2NoZWR1bGluZyJ9.8VAqvELFmly9lA0kUYunsgDWPk2BoVaUeAGOd2DTsWI";


async function main() {
  console.log("🚀 Starting booking test with multiple snacks...");

  try {
    // ---------------- CREATE BOOKING ----------------
    const bookingReq = {
      schedule_id: 3, // existing schedule ID
      schedule_token: SCHEDULE_TOKEN,
      seats: ["A1", "A2"], // valid seat numbers
      snacks: [
        {
          schedule_snack_id: 3, // first snack for schedule
          snack_token: SNACK_TOKEN_1,
          quantity: 2,
          price: 3.5
        },
        
      ],
      total_amount: 30.5
    };

    console.log("📤 Booking payload:", bookingReq);

    const bookingRes = await axios.post(`${API_BASE}/bookings`, bookingReq, {
      headers: { Authorization: `Bearer ${TOKEN}` }
    });

   
  console.log("✅ Booking created:", bookingRes.data.snacks);
  

  } catch (error) {
    console.error(
      "🚨 Unexpected error:",
      error.response ? error.response.data : error.message
    );
  }
}

main();
