const axios = require("axios");
const FormData = require("form-data");
const fs = require("fs");
const path = require("path");

const API_BASE = "http://localhost:8082/api/admin";
const TOKEN = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTgzMjMzNzUsInJvbGUiOiJhZG1pbiIsInR5cGUiOiJhY2Nlc3MiLCJ1c2VyX2lkIjoyfQ.CWva3TpF84x6A1LojvzQrG1copslBigZFe4ZJW8fERA"; // Replace with your JWT

async function main() {
  try {
    // ---------------- MOVIES ----------------
    
    const movieForm = new FormData();
    movieForm.append("title", "american pie");
    movieForm.append("description", "A mind-bending thriller about dreams within dreams.");
    movieForm.append("duration", 148);
    movieForm.append("genre", "Sci-Fi,Thriller");
    movieForm.append("release_year", 2010);
    movieForm.append("rating", 8.8);
    movieForm.append("image_poster_url", fs.createReadStream(path.join("C:/Users/nesre/Downloads", "inception.jpg")));

    const movieRes = await axios.post(`${API_BASE}/movies`, movieForm, {
      headers: { Authorization: `Bearer ${TOKEN}`, ...movieForm.getHeaders() },
    });
    console.log("Movie created:", movieRes.data);
    const movieId = movieRes.data.movie.id;

    // // Get all movies
    // const movies = await axios.get(`${API_BASE}/movies`, { headers: { Authorization: `Bearer ${TOKEN}` } });
    // console.log("All movies:", movies.data);

    // Update movie
    // await axios.put(`${API_BASE}/movies/${movieId}`, { title: "Inception Updated" }, { headers: { Authorization: `Bearer ${TOKEN}` } });

    // // Delete movie
    // await axios.delete(`${API_BASE}/movies/${movieId}`, { headers: { Authorization: `Bearer ${TOKEN}` } });
    

    // ---------------- SNACKS ----------------
    
    // const snackForm = new FormData();
    // snackForm.append("name", "hotdog");
    // snackForm.append("price", 5.5);
    // snackForm.append("description", "A delicious hotdog with mustard and ketchup");
    // snackForm.append("category", "Food");
    // snackForm.append("snack_image_url", fs.createReadStream(path.join("C:/Users/nesre/Downloads", "popcorn.jpg")));

    // const snackRes = await axios.post(`${API_BASE}/snacks`, snackForm, {
    //   headers: { Authorization: `Bearer ${TOKEN}`, ...snackForm.getHeaders() },
    // });
    // console.log("Snack created:", snackRes.data);
    // const snackId = snackRes.data.snack.id;

    // Get all snacks
    // const snacks = await axios.get(`${API_BASE}/snacks`, { headers: { Authorization: `Bearer ${TOKEN}` } });
    // console.log("All snacks:", snacks.data);

    // Update snack
    // await axios.put(`${API_BASE}/snacks/${snackId}`, { name:"coffee",price: 2.0,category:"drink" }, { headers: { Authorization: `Bearer ${TOKEN}` } });

    // Delete snack
    // await axios.delete(`${API_BASE}/snacks/${snackId}`, { headers: { Authorization: `Bearer ${TOKEN}` } });
    

    // ---------------- HALLS ----------------
    
    const hallData = { name: "Hall 1", capacity: 120 , location: "First Floor"};
    const hallRes = await axios.post(`${API_BASE}/halls`, hallData, { headers: { Authorization: `Bearer ${TOKEN}` } });
    console.log("Hall created:", hallRes.data);
    const hallId = hallRes.data.hall.id;

    // Get all halls
    // const halls = await axios.get(`${API_BASE}/halls`, { headers: { Authorization: `Bearer ${TOKEN}` } });
    // console.log("All halls:", halls.data);

    // Update hall
    // await axios.put(`${API_BASE}/halls/${hallId}`, { capacity: 150 }, { headers: { Authorization: `Bearer ${TOKEN}` } });

    // Delete hall
    // await axios.delete(`${API_BASE}/halls/${hallId}`, { headers: { Authorization: `Bearer ${TOKEN}` } });
    

    // ---------------- SCHEDULES ----------------
    
    const scheduleData = {
      movie_id: movieId,
      hall_id: hallId,
      show_time: new Date().toISOString(),
      available_seats: 120
    };

    const scheduleRes = await axios.post(`${API_BASE}/schedules`, scheduleData, { headers: { Authorization: `Bearer ${TOKEN}` } });
    console.log("Schedule created:", scheduleRes.data);
    const scheduleId = scheduleRes.data.schedule.id;

    // Get all schedules
    const schedules = await axios.get(`${API_BASE}/schedules`, { headers: { Authorization: `Bearer ${TOKEN}` } });
    console.log("All schedules:", schedules.data);

    // Update schedule
    // await axios.put(`${API_BASE}/schedules/${scheduleId}`, { available_seats: 100 }, { headers: { Authorization: `Bearer ${TOKEN}` } });

    // Delete schedule
    // await axios.delete(`${API_BASE}/schedules/${scheduleId}`, { headers: { Authorization: `Bearer ${TOKEN}` } });
    

    // ---------------- SCHEDULE-SPECIFIC SNACKS ----------------
    /*
    const scheduleSnackData = { snack_id: snackId, price: 6.0 };
    const scheduleSnackRes = await axios.post(`${API_BASE}/schedules/${scheduleId}/snacks`, scheduleSnackData, { headers: { Authorization: `Bearer ${TOKEN}` } });
    console.log("Schedule Snack added:", scheduleSnackRes.data);
    const scheduleSnackId = scheduleSnackRes.data.schedule_snack.id;

    // Update schedule snack
    await axios.put(`${API_BASE}/schedules/snacks/${scheduleSnackId}`, { price: 6.5 }, { headers: { Authorization: `Bearer ${TOKEN}` } });

    // Delete schedule snack
    await axios.delete(`${API_BASE}/schedules/snacks/${scheduleSnackId}`, { headers: { Authorization: `Bearer ${TOKEN}` } });
    */

  } catch (error) {
    console.error("Error:", error.response ? error.response.data : error.message);
  }
}

main();
