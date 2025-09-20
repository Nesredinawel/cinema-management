const axios = require("axios");
const FormData = require("form-data");
const fs = require("fs");
const path = require("path");

const API_BASE = "http://localhost:8082/api/admin";
const TOKEN ="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTgzNjU2MzMsInJvbGUiOiJhZG1pbiIsInR5cGUiOiJhY2Nlc3MiLCJ1c2VyX2lkIjoxfQ.uJlxad3hYK69xnxE-b3fNBTGh_9AIP1GDNOO4rB696s"; // Replace with your JWT

async function main() {
  try {
    // ---------------- GENRES ----------------
    const genre1 = await axios.post(
      `${API_BASE}/genres`,
      { name: "Sci-Fi" },
      { headers: { Authorization: `Bearer ${TOKEN}` } }
    );
    const genre2 = await axios.post(
      `${API_BASE}/genres`,
      { name: "Thriller" },
      { headers: { Authorization: `Bearer ${TOKEN}` } }
    );

    console.log("Genres created:", genre1.data, genre2.data);
    const genreIds = [genre1.data.genre.id, genre2.data.genre.id];

    // ---------------- MOVIES ----------------
    const movieForm = new FormData();
    movieForm.append("title", "Inception");
    movieForm.append(
      "description",
      "A mind-bending thriller about dreams within dreams."
    );
    movieForm.append("duration", 148);
    movieForm.append("release_year", 2010);
    movieForm.append("rating", 8.8);
    movieForm.append(
      "image_poster_url",
      fs.createReadStream(path.join("C:/Users/nesre/Downloads", "inception.jpg"))
    );

    // 👇 Send genre IDs as JSON
    movieForm.append("genre_ids", JSON.stringify(genreIds));

    const movieRes = await axios.post(`${API_BASE}/movies`, movieForm, {
      headers: { Authorization: `Bearer ${TOKEN}`, ...movieForm.getHeaders() },
    });
    console.log("Movie created:", movieRes.data);
    const movieId = movieRes.data.movie.id;

    // ---------------- HALLS ----------------
    const hallData = { name: "Hall 1", capacity: 120, location: "First Floor" };
    const hallRes = await axios.post(`${API_BASE}/halls`, hallData, {
      headers: { Authorization: `Bearer ${TOKEN}` },
    });
    console.log("Hall created:", hallRes.data);
    const hallId = hallRes.data.hall.id;

    // ---------------- SCHEDULES ----------------
    const scheduleData = {
      movie_id: movieId,
      hall_id: hallId,
      show_time: new Date().toISOString(),
      available_seats: 120,
    };

    const scheduleRes = await axios.post(`${API_BASE}/schedules`, scheduleData, {
      headers: { Authorization: `Bearer ${TOKEN}` },
    });
    console.log("Schedule created:", scheduleRes.data);
    const scheduleId = scheduleRes.data.schedule.id;

    // Get all schedules
    const schedules = await axios.get(`${API_BASE}/schedules`, {
      headers: { Authorization: `Bearer ${TOKEN}` },
    });
    console.log("All schedules:", schedules.data);
  } catch (error) {
    console.error(
      "Error:",
      error.response ? error.response.data : error.message
    );
  }
}

main();
